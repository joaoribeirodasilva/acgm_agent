package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	evnt "github.com/jonhoo/go-events"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

const METRIC_CPU_NAME = "cpu"

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
	nm.status = STATUS_STOPPED
	nm.db = db
	go nm.EventHandler()
	return nm
}

func (nm *CPUMeter) Start() {

	nm.FireEvent(STATUS_STARTING)

	nm.FireEvent(STATUS_STARTED)

}

func (nm *CPUMeter) Stop() {

	for nm.status == STATUS_COLLECTING || nm.status == STATUS_AGGREGATING {
		time.Sleep(50 * time.Millisecond)
	}

	nm.FireEvent(STATUS_STOPPING)
	nm.Aggregate()

	for nm.status == STATUS_STOPPING || nm.status == STATUS_AGGREGATING {

		time.Sleep(50 * time.Millisecond)
	}

	nm.FireEvent(STATUS_STOPPED)

}

func (nm *CPUMeter) GetStatus() MeterStatus {
	return nm.status
}

func (nm *CPUMeter) GetName() string {
	return METRIC_CPU_NAME
}

func (nm *CPUMeter) Collect() {

	if nm.status != STATUS_WAITING {
		return
	}
	nm.FireEvent(STATUS_COLLECTING)
	metric := CPUMetric{}

	metric.DateTime = time.Now()

	// CPU load average
	average, err := load.Avg()
	if err != nil || average == nil {
		logger.Log.Error().Err(err).Msg("failed to get CPU load average data")
	} else {
		metric.Total.Average = *average
	}

	// CPU miscellaneous data
	misc, err := load.Misc()
	if err != nil || misc == nil {
		logger.Log.Error().Err(err).Msg("failed to get CPU load miscellaneous data")
	} else {
		metric.Total.Misc = *misc
	}

	// Physical CPU (Cores)
	physical, err := cpu.Counts(false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host physical CPU count")
	}
	metric.Total.Physical = physical

	// Logical CPU (Threads)
	logical, err := cpu.Counts(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host logical CPU count")
	}
	metric.Total.Logical = logical

	// Total CPU times
	total_times, err := cpu.Times(false)
	if err != nil || len(total_times) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU total times")
	} else {
		metric.Total.Times = total_times[0]
	}

	// Total CPU percentage usage
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		logger.Log.Error().Err(err).Msg("Failed to get host CPU percent total")
	} else {
		metric.Total.Percent = percent[0]
	}

	// CPU threads information
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU and CPU cores information")
		nm.Metrics = append(nm.Metrics, metric)
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
	} else {
		for idx, times := range core_times {
			metric.Cores[idx].Times = times
		}
	}

	// CPU threads percentage usage
	percents, err := cpu.Percent(0, true)
	if err != nil || len(info) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host CPU cores percent")
	} else {
		for idx, percent := range percents {
			metric.Cores[idx].Percent = percent
		}
	}

	nm.AppendMetric(&metric)
}

func (nm *CPUMeter) AppendMetric(metric *CPUMetric) {
	metric.Duration = (time.Now().UnixNano() / int64(time.Millisecond)) - (metric.DateTime.UnixNano() / int64(time.Millisecond))
	nm.Metrics = append(nm.Metrics, *metric)
	nm.Aggregate()
}

func (nm *CPUMeter) Aggregate() {

	if len(nm.Metrics) == 0 || (len(nm.Metrics) < nm.conf.Metrics.Cpu.Aggregate && nm.status != STATUS_STOPPING) {
		nm.FireEvent(STATUS_WAITING)
		return
	}

	nm.FireEvent(STATUS_AGGREGATING)

	nm.FireEvent(STATUS_WAITING)
}

func (nm *CPUMeter) FireEvent(status MeterStatus) {

	nm.status = status
	logger.Log.Debug().Msgf("status changed to %d", nm.status)
	evnt.Signal("meter.changed.cpu")
}

func (nm *CPUMeter) EventHandler() {

	events := evnt.Listen("meter.meter.")

	for event := range events {
		switch event.Tag {
		case "meter.meter.start":
			nm.Start()
		case "meter.meter.collect":
			nm.Collect()
		case "meter.meter.stop":
			nm.Stop()
		}
		if nm.status == STATUS_STOPPED {
			break
		}
	}
}
