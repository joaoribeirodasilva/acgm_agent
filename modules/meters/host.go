package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	evnt "github.com/jonhoo/go-events"
	"github.com/shirou/gopsutil/host"
)

var METRIC_HOST_NAME = "host"

type HostMetric struct {
	MetricTimes
	Host         host.InfoStat          `json:"host" yaml:"host"`
	LSB          host.LSB               `json:"lsb" yaml:"lsb"`
	Temperatures []host.TemperatureStat `json:"temperatures" yaml:"temperatures"`
	Users        []host.UserStat        `json:"users" yaml:"users"`
}

type HostMeter struct {
	MeterControl
	HostID  string       `json:"host_id" yaml:"host_id"`
	Metrics []HostMetric `json:"metrics" yaml:"metrics"`
}

func NewHostMeter(conf *config.Config, db *database.Db) *HostMeter {
	nm := &HostMeter{}
	nm.name = METRIC_HOST_NAME
	nm.conf = conf
	nm.status = STATUS_STOPPED
	nm.db = db
	go nm.EventHandler()
	return nm
}

func (nm *HostMeter) GetHostID() (string, error) {

	host_id, err := host.HostID()
	nm.HostID = host_id

	return host_id, err
}

func (nm *HostMeter) Start() {

	nm.FireEvent(STATUS_STARTING)

}

func (nm *HostMeter) Stop() {

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

func (nm *HostMeter) GetStatus() MeterStatus {
	return nm.status
}

func (nm *HostMeter) FireEvent(status MeterStatus) {

	nm.status = status
	evnt.Signal("meter.changed.host")
}

func (nm *HostMeter) EventHandler() {

	evnt.Verbosity = 0
	if nm.conf.Settings.Debug {
		evnt.Verbosity = 3
	}

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

func (nm *HostMeter) GetName() string {
	return METRIC_HOST_NAME
}

func (nm *HostMeter) Collect() {

	if nm.status != STATUS_WAITING {
		return
	}
	nm.FireEvent(STATUS_COLLECTING)
	metric := HostMetric{}

	metric.DateTime = time.Now()

	// // CPU load average
	// average, err := load.Avg()
	// if err != nil || average == nil {
	// 	logger.Log.Error().Err(err).Msg("failed to get CPU load average data")
	// } else {
	// 	metric.Total.Average = *average
	// }

	// // CPU miscellaneous data
	// misc, err := load.Misc()
	// if err != nil || misc == nil {
	// 	logger.Log.Error().Err(err).Msg("failed to get CPU load miscellaneous data")
	// } else {
	// 	metric.Total.Misc = *misc
	// }

	// // Physical CPU (Cores)
	// physical, err := cpu.Counts(false)
	// if err != nil {
	// 	logger.Log.Error().Err(err).Msg("failed to get host physical CPU count")
	// }
	// metric.Total.Physical = physical

	// // Logical CPU (Threads)
	// logical, err := cpu.Counts(true)
	// if err != nil {
	// 	logger.Log.Error().Err(err).Msg("failed to get host logical CPU count")
	// }
	// metric.Total.Logical = logical

	// // Total CPU times
	// total_times, err := cpu.Times(false)
	// if err != nil || len(total_times) == 0 {
	// 	logger.Log.Error().Err(err).Msg("failed to get host CPU total times")
	// } else {
	// 	metric.Total.Times = total_times[0]
	// }

	// // Total CPU percentage usage
	// percent, err := cpu.Percent(0, false)
	// if err != nil || len(percent) == 0 {
	// 	logger.Log.Error().Err(err).Msg("Failed to get host CPU percent total")
	// } else {
	// 	metric.Total.Percent = percent[0]
	// }

	// // CPU threads information
	// info, err := cpu.Info()
	// if err != nil || len(info) == 0 {
	// 	logger.Log.Error().Err(err).Msg("failed to get host CPU and CPU cores information")
	// 	nm.Metrics = append(nm.Metrics, metric)
	// 	return
	// } else {
	// 	metric.Total.Info = info[0]
	// 	for _, cpu := range info {
	// 		core := CPUCore{
	// 			Info: cpu,
	// 		}
	// 		metric.Cores = append(metric.Cores, core)
	// 	}
	// }

	// // CPU threads times
	// core_times, err := cpu.Times(true)
	// if err != nil {
	// 	logger.Log.Error().Err(err).Msg("failed to get host CPU cores times")
	// } else {
	// 	for idx, times := range core_times {
	// 		metric.Cores[idx].Times = times
	// 	}
	// }

	// // CPU threads percentage usage
	// percents, err := cpu.Percent(0, true)
	// if err != nil || len(info) == 0 {
	// 	logger.Log.Error().Err(err).Msg("failed to get host CPU cores percent")
	// } else {
	// 	for idx, percent := range percents {
	// 		metric.Cores[idx].Percent = percent
	// 	}
	// }

	nm.AppendMetric(&metric)
}

func (nm *HostMeter) AppendMetric(metric *HostMetric) {
	metric.Duration = (time.Now().UnixNano() / int64(time.Millisecond)) - (metric.DateTime.UnixNano() / int64(time.Millisecond))
	nm.Metrics = append(nm.Metrics, *metric)
	nm.Aggregate()
}

func (nm *HostMeter) Aggregate() {

	if len(nm.Metrics) == 0 || (len(nm.Metrics) < nm.conf.Metrics.Cpu.Aggregate && nm.status != STATUS_STOPPING) {
		nm.FireEvent(STATUS_WAITING)
		return
	}

	nm.FireEvent(STATUS_AGGREGATING)

	nm.FireEvent(STATUS_WAITING)
}
