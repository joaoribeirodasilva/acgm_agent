package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
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

type HostTemperaturesAggregate struct {
	SensorKey          string
	AverageTemperature float64
	MaxTemperature     float64
	MinTemperature     float64
}

type HostUsersAggregate struct {
	User     string
	Terminal string
	Host     string
	Started  string
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

	nm.FireEvent(STATUS_STARTED)

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

	info, err := host.Info()
	if err != nil || info == nil {
		logger.Log.Error().Err(err).Msg("failed to get host information")
	} else {
		metric.Host = *info
	}

	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host temperatures")
	} else {
		metric.Temperatures = temps
	}

	users, err := host.Users()
	if err != nil || len(users) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host users")
	} else {
		metric.Users = users
	}

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

	for _, metric := range nm.Metrics {

	}

	nm.FireEvent(STATUS_WAITING)
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
