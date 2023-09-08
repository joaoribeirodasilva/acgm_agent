package meters

import (
	"strings"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	"biqx.com.br/acgm_agent/modules/models"
	evnt "github.com/jonhoo/go-events"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

const METRIC_CPU_NAME = "cpu"

type CPUAggregator struct {
	HostID        string
	CPU           int32
	Percent       MeterMinMaxAvgFloat64
	Load1         MeterMinMaxAvgFloat64
	Load5         MeterMinMaxAvgFloat64
	Load15        MeterMinMaxAvgFloat64
	ProcsTotal    MeterMinMaxAvgInt
	ProcsCreated  MeterMinMaxAvgInt
	ProcsRunning  MeterMinMaxAvgInt
	ProcsBlocked  MeterMinMaxAvgInt
	Ctxt          MeterMinMaxAvgInt
	PhysicalCores int
	LogicalCores  int
}

type CPUAggregatorTimes struct {
	HostID         string
	CPU            string
	CoreIndex      int32
	Percent        MeterMinMaxAvgFloat64
	TimesTotal     MeterMinMaxAvgFloat64
	TimesUser      MeterMinMaxAvgFloat64
	TimesSystem    MeterMinMaxAvgFloat64
	TimesIdle      MeterMinMaxAvgFloat64
	TimesNice      MeterMinMaxAvgFloat64
	TimesIOWait    MeterMinMaxAvgFloat64
	TimesIrq       MeterMinMaxAvgFloat64
	TimesSoftirq   MeterMinMaxAvgFloat64
	TimesSteal     MeterMinMaxAvgFloat64
	TimesGuest     MeterMinMaxAvgFloat64
	TimesGuestNice MeterMinMaxAvgFloat64
}

type CPUTotal struct {
	Info     cpu.InfoStat  `json:"info" yaml:"info"`
	Percent  float64       `json:"percent" yaml:"percent"`
	Times    cpu.TimesStat `json:"times" yaml:"times"`
	Threads  int           `json:"threads" yaml:"threads"`
	Average  load.AvgStat  `json:"average" yaml:"average"`
	Misc     load.MiscStat `json:"misc" yaml:"misc"`
	Physical int           `json:"physical" yaml:"physical"`
	Logical  int           `json:"logical" yaml:"logical"`
}

type CPUCore struct {
	Info    cpu.InfoStat  `json:"info" yaml:"info"`
	Percent float64       `json:"percent" yaml:"percent"`
	Times   cpu.TimesStat `json:"times" yaml:"times"`
}

type CPUMetric struct {
	MetricTimes
	Total CPUTotal  `json:"total" yaml:"total"`
	Cores []CPUCore `json:"cores" yaml:"cores"`
}

type CPUMeter struct {
	MeterControl
	Metrics []CPUMetric `json:"metrics" yaml:"metrics"`
}

func NewCPUMeter(host string, conf *config.Config, db *database.Db) *CPUMeter {
	nm := &CPUMeter{}
	nm.host = host
	nm.name = METRIC_CPU_NAME
	nm.conf = conf
	nm.status = STATUS_STARTING
	nm.db = db
	nm.EventHandler()
	return nm
}

func (nm *CPUMeter) Start() {

	logger.Log.Debug().Msg("starting")
	nm.FireEvent(STATUS_STARTING)

	// TODO: Validate tables

	logger.Log.Debug().Msg("started")
	nm.FireEvent(STATUS_STARTED)

	logger.Log.Debug().Msg("waiting")
	nm.FireEvent(STATUS_WAITING)

}

func (nm *CPUMeter) Stop() {

	logger.Log.Debug().Msg("waiting for running operations")
	for nm.status == STATUS_COLLECTING || nm.status == STATUS_AGGREGATING {
		time.Sleep(50 * time.Millisecond)
	}

	logger.Log.Debug().Msg("stopping")
	nm.FireEvent(STATUS_STOPPING)
	nm.Aggregate()

	nm.status = STATUS_STOPPED
	evnt.Signal("meter.cpu.stopped")

}

func (nm *CPUMeter) GetStatus() MeterStatus {
	return nm.status
}

func (nm *CPUMeter) GetName() string {
	return METRIC_CPU_NAME
}

func (nm *CPUMeter) Collect() {

	logger.Log.Debug().Msg("checking for collect")
	if nm.status != STATUS_WAITING {
		return
	}

	logger.Log.Debug().Msg("collecting")
	nm.FireEvent(STATUS_COLLECTING)
	metric := CPUMetric{}

	metric.DateTime = time.Now()

	// CPU load average
	average, err := load.Avg()
	if err != nil || average == nil {
		logger.Log.Error().Err(err).Msg("failed to get CPU load average data")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Total.Average = *average
	}

	// CPU miscellaneous data
	misc, err := load.Misc()
	if err != nil || misc == nil {
		logger.Log.Error().Err(err).Msg("failed to get CPU load miscellaneous data")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Total.Misc = *misc
	}

	// Physical CPUs (Cores)
	physical, err := cpu.Counts(false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host physical CPU count")
		logger.Log.Debug().Msg("collect error")
	}
	metric.Total.Physical = physical

	// Logical CPUs
	logical, err := cpu.Counts(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host logical CPU count")
		logger.Log.Debug().Msg("collect error")
	}
	metric.Total.Logical = logical

	// Total CPU times
	total_times, err := cpu.Times(false)
	if err != nil || len(total_times) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU total times")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Total.Times = total_times[0]
	}

	// Total CPU percentage usage
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		logger.Log.Error().Err(err).Msg("Failed to get host CPU percent total")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Total.Percent = percent[0]
	}

	// Multiple
	// CPU threads information
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU and CPU cores information")
		logger.Log.Debug().Msg("collect error")
		nm.AppendMetric(&metric)
		return
	} else {
		metric.Total.Info = info[0]
		for _, cpu := range info {
			core := CPUCore{
				Info: cpu,
			}
			metric.Cores = append(metric.Cores, core)
		}
	}

	// CPU threads times
	core_times, err := cpu.Times(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host CPU cores times")
		logger.Log.Debug().Msg("collect error")
	} else {
		for idx, times := range core_times {
			metric.Cores[idx].Times = times
		}
	}

	// CPU threads percentage usage
	percents, err := cpu.Percent(0, true)
	if err != nil || len(info) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU cores percent")
		logger.Log.Debug().Msg("collect error")
	} else {
		for idx, percent := range percents {
			metric.Cores[idx].Percent = percent
		}
	}

	nm.AppendMetric(&metric)
	logger.Log.Debug().Msg("collect success")
}

func (nm *CPUMeter) AppendMetric(metric *CPUMetric) {
	metric.Duration = (time.Now().UnixNano() / int64(time.Millisecond)) - (metric.DateTime.UnixNano() / int64(time.Millisecond))
	nm.Metrics = append(nm.Metrics, *metric)
	nm.Aggregate()
}

func (nm *CPUMeter) Aggregate() {

	logger.Log.Debug().Msg("checking for aggregation")
	if len(nm.Metrics) == 0 || (len(nm.Metrics) < nm.conf.Metrics.Cpu.Aggregate && nm.status != STATUS_STOPPING) {
		logger.Log.Debug().Msg("waiting")
		nm.FireEvent(STATUS_WAITING)
		return
	}

	nm.FireEvent(STATUS_AGGREGATING)
	logger.Log.Debug().Msg("aggregating")

	cpu_core_info := make(map[int32]*models.CPUCoreInfo, 0)
	cpu_times := make(map[int32]*CPUAggregatorTimes, 0)
	cpu_agg := &CPUAggregator{
		HostID: nm.host,
		CPU:    nm.Metrics[0].Total.Info.CPU,
	}

	for i, metric := range nm.Metrics {

		//CPU Loads
		cpu_agg.Percent.Aggregate(metric.Total.Percent)
		cpu_agg.Load1.Aggregate(metric.Total.Average.Load1)
		cpu_agg.Load5.Aggregate(metric.Total.Average.Load5)
		cpu_agg.Load15.Aggregate(metric.Total.Average.Load15)
		cpu_agg.ProcsTotal.Aggregate(metric.Total.Misc.ProcsTotal)
		cpu_agg.ProcsCreated.Aggregate(metric.Total.Misc.ProcsCreated)
		cpu_agg.ProcsRunning.Aggregate(metric.Total.Misc.ProcsRunning)
		cpu_agg.ProcsBlocked.Aggregate(metric.Total.Misc.ProcsBlocked)
		cpu_agg.Ctxt.Aggregate(metric.Total.Misc.Ctxt)
		cpu_agg.PhysicalCores = metric.Total.Physical
		cpu_agg.LogicalCores = metric.Total.Logical

		// CPU Info (totals)
		_, ok := cpu_core_info[metric.Total.Info.CPU]
		if !ok {
			cpu_core_info[metric.Total.Info.CPU] = &models.CPUCoreInfo{
				CPU:             metric.Total.Info.CPU,
				HostID:          nm.host,
				CoreIndex:       &i,
				VendorID:        metric.Total.Info.VendorID,
				Family:          metric.Total.Info.Family,
				Model:           metric.Total.Info.Model,
				Stepping:        metric.Total.Info.Stepping,
				PhysicalID:      metric.Total.Info.PhysicalID,
				CoreID:          metric.Total.Info.CoreID,
				Cores:           metric.Total.Info.Cores,
				ModelName:       metric.Total.Info.ModelName,
				Mhz:             metric.Total.Info.Mhz,
				CacheSize:       metric.Total.Info.CacheSize,
				Flags:           strings.Join(metric.Total.Info.Flags, ", "),
				Microcode:       metric.Total.Info.Microcode,
				CollectedAt:     nm.Metrics[0].DateTime,
				CollectedMillis: uint64(nm.Metrics[0].Duration),
			}
		}

		// CPU Times (totals)
		_, ok = cpu_times[metric.Total.Info.CPU]
		if !ok {
			cpu_times[metric.Total.Info.CPU] = &CPUAggregatorTimes{
				HostID:    nm.host,
				CPU:       metric.Total.Times.CPU,
				CoreIndex: metric.Total.Info.CPU,
			}
		}
		cpu_times[metric.Total.Info.CPU].Percent.Aggregate(metric.Total.Percent)
		cpu_times[metric.Total.Info.CPU].TimesTotal.Aggregate(metric.Total.Times.Total())
		cpu_times[metric.Total.Info.CPU].TimesUser.Aggregate(metric.Total.Times.User)
		cpu_times[metric.Total.Info.CPU].TimesSystem.Aggregate(metric.Total.Times.System)
		cpu_times[metric.Total.Info.CPU].TimesIdle.Aggregate(metric.Total.Times.Idle)
		cpu_times[metric.Total.Info.CPU].TimesNice.Aggregate(metric.Total.Times.Nice)
		cpu_times[metric.Total.Info.CPU].TimesIOWait.Aggregate(metric.Total.Times.Iowait)
		cpu_times[metric.Total.Info.CPU].TimesIrq.Aggregate(metric.Total.Times.Irq)
		cpu_times[metric.Total.Info.CPU].TimesSoftirq.Aggregate(metric.Total.Times.Softirq)
		cpu_times[metric.Total.Info.CPU].TimesSteal.Aggregate(metric.Total.Times.Steal)
		cpu_times[metric.Total.Info.CPU].TimesGuest.Aggregate(metric.Total.Times.Guest)
		cpu_times[metric.Total.Info.CPU].TimesGuestNice.Aggregate(metric.Total.Times.GuestNice)

		for idx, core := range metric.Cores {
			// Cores info
			_, ok := cpu_core_info[core.Info.CPU]
			if !ok {
				cpu_core_info[core.Info.CPU] = &models.CPUCoreInfo{
					CPU:             core.Info.CPU,
					HostID:          nm.host,
					CoreIndex:       &idx,
					VendorID:        core.Info.VendorID,
					Family:          core.Info.Family,
					Model:           core.Info.Model,
					Stepping:        core.Info.Stepping,
					PhysicalID:      core.Info.PhysicalID,
					CoreID:          core.Info.CoreID,
					Cores:           core.Info.Cores,
					ModelName:       core.Info.ModelName,
					Mhz:             core.Info.Mhz,
					CacheSize:       core.Info.CacheSize,
					Flags:           strings.Join(core.Info.Flags, ", "),
					Microcode:       core.Info.Microcode,
					CollectedAt:     nm.Metrics[0].DateTime,
					CollectedMillis: uint64(nm.Metrics[0].Duration),
				}
				// TODO: avg duration
			}
			// Cores times
			_, ok = cpu_times[core.Info.CPU]
			if !ok {
				cpu_times[core.Info.CPU] = &CPUAggregatorTimes{
					CPU:       core.Times.CPU,
					CoreIndex: core.Info.CPU,
				}
			}
			cpu_times[core.Info.CPU].Percent.Aggregate(core.Percent)
			cpu_times[core.Info.CPU].TimesTotal.Aggregate(core.Times.Total())
			cpu_times[core.Info.CPU].TimesUser.Aggregate(core.Times.User)
			cpu_times[core.Info.CPU].TimesSystem.Aggregate(core.Times.System)
			cpu_times[core.Info.CPU].TimesIdle.Aggregate(core.Times.Idle)
			cpu_times[core.Info.CPU].TimesNice.Aggregate(core.Times.Nice)
			cpu_times[core.Info.CPU].TimesIOWait.Aggregate(core.Times.Iowait)
			cpu_times[core.Info.CPU].TimesIrq.Aggregate(core.Times.Irq)
			cpu_times[core.Info.CPU].TimesSoftirq.Aggregate(core.Times.Softirq)
			cpu_times[core.Info.CPU].TimesSteal.Aggregate(core.Times.Steal)
			cpu_times[core.Info.CPU].TimesGuest.Aggregate(core.Times.Guest)
			cpu_times[core.Info.CPU].TimesGuestNice.Aggregate(core.Times.GuestNice)
		}
	}

	model_cpu_core_info := make([]*models.CPUCoreInfo, 0)
	for _, info := range cpu_core_info {
		model_cpu_core_info = append(model_cpu_core_info, info)
	}

	model_cpu_times := make([]*models.CPUTimes, 0)
	for _, times := range cpu_times {
		t := &models.CPUTimes{
			HostID:            times.HostID,
			CPU:               times.CPU,
			CoreIndex:         times.CoreIndex,
			PercentAvg:        times.Percent.Avg,
			PercentMin:        times.Percent.Min,
			PercentMax:        times.Percent.Max,
			TimesTotalAvg:     times.TimesTotal.Avg,
			TimesTotalMin:     times.TimesTotal.Min,
			TimesTotalMax:     times.TimesTotal.Max,
			TimesUserAvg:      times.TimesUser.Avg,
			TimesUserMin:      times.TimesUser.Min,
			TimesUserMax:      times.TimesUser.Max,
			TimesSystemAvg:    times.TimesSystem.Avg,
			TimesSystemMin:    times.TimesSystem.Min,
			TimesSystemMax:    times.TimesSystem.Max,
			TimesIdleAvg:      times.TimesIdle.Avg,
			TimesIdleMin:      times.TimesIdle.Min,
			TimesIdleMax:      times.TimesIdle.Max,
			TimesNiceAvg:      times.TimesNice.Avg,
			TimesNiceMin:      times.TimesNice.Min,
			TimesNiceMax:      times.TimesNice.Max,
			TimesIOWaitAvg:    times.TimesIOWait.Avg,
			TimesIOWaitMin:    times.TimesIOWait.Min,
			TimesIOWaitMax:    times.TimesIOWait.Max,
			TimesIrqAvg:       times.TimesIrq.Avg,
			TimesIrqMin:       times.TimesIrq.Min,
			TimesIrqMax:       times.TimesIrq.Max,
			TimesSoftirqAvg:   times.TimesSoftirq.Avg,
			TimesSoftirqMin:   times.TimesSoftirq.Min,
			TimesSoftirqMax:   times.TimesSoftirq.Max,
			TimesStealAvg:     times.TimesSteal.Avg,
			TimesStealMin:     times.TimesSteal.Min,
			TimesStealMax:     times.TimesSteal.Max,
			TimesGuestAvg:     times.TimesGuest.Avg,
			TimesGuestMin:     times.TimesGuest.Min,
			TimesGuestMax:     times.TimesGuest.Max,
			TimesGuestNiceAvg: times.TimesGuestNice.Avg,
			TimesGuestNiceMin: times.TimesGuestNice.Min,
			TimesGuestNiceMax: times.TimesGuestNice.Max,
			CollectedAt:       nm.Metrics[0].DateTime,
			CollectedMillis:   uint64(nm.Metrics[0].Duration),
		}
		model_cpu_times = append(model_cpu_times, t)
	}

	model_cpus := &models.CPU{
		HostID:          cpu_agg.HostID,
		CPU:             cpu_agg.CPU,
		PercentAvg:      cpu_agg.Percent.Avg,
		PercentMin:      cpu_agg.Percent.Min,
		PercentMax:      cpu_agg.Percent.Max,
		Load1Avg:        cpu_agg.Load1.Avg,
		Load1Min:        cpu_agg.Load1.Min,
		Load1Max:        cpu_agg.Load1.Max,
		Load5Avg:        cpu_agg.Load5.Avg,
		Load5Min:        cpu_agg.Load5.Min,
		Load5Max:        cpu_agg.Load5.Max,
		Load15Avg:       cpu_agg.Load15.Avg,
		Load15Min:       cpu_agg.Load15.Min,
		Load15Max:       cpu_agg.Load15.Max,
		ProcsTotalAvg:   cpu_agg.ProcsTotal.Avg,
		ProcsTotalMin:   cpu_agg.ProcsTotal.Min,
		ProcsTotalMax:   cpu_agg.ProcsTotal.Max,
		ProcsCreatedAvg: cpu_agg.ProcsCreated.Avg,
		ProcsCreatedMin: cpu_agg.ProcsCreated.Min,
		ProcsCreatedMax: cpu_agg.ProcsCreated.Max,
		ProcsRunningAvg: cpu_agg.ProcsRunning.Avg,
		ProcsRunningMin: cpu_agg.ProcsRunning.Min,
		ProcsRunningMax: cpu_agg.ProcsRunning.Max,
		ProcsBlockedAvg: cpu_agg.ProcsBlocked.Avg,
		ProcsBlockedMin: cpu_agg.ProcsBlocked.Min,
		ProcsBlockedMax: cpu_agg.ProcsBlocked.Max,
		CtxtAvg:         cpu_agg.Ctxt.Avg,
		CtxtMin:         cpu_agg.Ctxt.Min,
		CtxtMax:         cpu_agg.Ctxt.Max,
		PhysicalCores:   cpu_agg.PhysicalCores,
		LogicalCores:    cpu_agg.LogicalCores,
		CollectedAt:     nm.Metrics[0].DateTime,
		CollectedMillis: uint64(nm.Metrics[0].Duration),
	}

	tx, err := nm.db.Begin()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to start transaction")
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}

	if err := tx.Create(&model_cpu_core_info).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to create cpu info records")
		nm.db.Rollback(tx)
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}
	logger.Log.Debug().Msgf("inserted %d cpu records", tx.RowsAffected)

	if err := tx.Create(&model_cpu_times).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to create cpu times records")
		nm.db.Rollback(tx)
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}
	logger.Log.Debug().Msgf("inserted %d cpu time records", tx.RowsAffected)

	if err := tx.Create(&model_cpus).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to create cpu load records")
		nm.db.Rollback(tx)
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}
	logger.Log.Debug().Msgf("inserted %d cpu load records", tx.RowsAffected)

	nm.db.Commit(tx)

	nm.Metrics = make([]CPUMetric, 0)

	logger.Log.Debug().Msg("aggregated success")
	logger.Log.Debug().Msg("waiting")
	nm.FireEvent(STATUS_WAITING)
}

func (nm *CPUMeter) FireEvent(status MeterStatus) {

	nm.status = status
	evnt.Signal("meter.changed.cpu")
}

func (nm *CPUMeter) EventHandler() {

	logger.Log.Debug().Msg("event handler start")
	events := evnt.Listen("meter.")

	go func() {
		for nm.status != STATUS_STOPPED {
			for event := range events {
				logger.Log.Debug().Msgf("received event: %s", event.Tag)
				switch event.Tag {
				case "meter.meter.start":
					nm.Start()
				case "meter.meter.collect":
					nm.Collect()
				case "meter.meter.stop":
					nm.Stop()
				case "meter.cpu.stopped":
					close(events)
				}
			}
		}
	}()
	logger.Log.Debug().Msg("event handler stopped")
}
