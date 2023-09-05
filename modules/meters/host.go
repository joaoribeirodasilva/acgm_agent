package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	"biqx.com.br/acgm_agent/modules/models"
	evnt "github.com/jonhoo/go-events"
	"github.com/shirou/gopsutil/host"
)

var METRIC_HOST_NAME = "host"

type HostMetric struct {
	MetricTimes
	Host         host.InfoStat          `json:"host" yaml:"host"`
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

	model_info := models.HostInfo{}
	temps := make(map[string]MeterMinMaxAvgFloat64)
	users := make(map[string]models.HostUsers)

	for idx, metric := range nm.Metrics {
		if idx == 0 {
			model_info.HostID = metric.Host.HostID
			model_info.Hostname = metric.Host.Hostname
			uptime := time.Unix(int64(metric.Host.Uptime), 0)
			model_info.Uptime = uptime
			boot_time := time.Unix(int64(metric.Host.BootTime), 0)
			model_info.BootTime = boot_time
			model_info.OS = metric.Host.OS
			model_info.Platform = metric.Host.Platform
			model_info.PlatformFamily = metric.Host.PlatformFamily
			model_info.PlatformVersion = metric.Host.PlatformVersion
			model_info.KernelVersion = metric.Host.KernelVersion
			model_info.KernelArch = metric.Host.KernelArch
			model_info.VirtualizationSystem = metric.Host.VirtualizationSystem
			model_info.VirtualizationRole = metric.Host.VirtualizationRole
			model_info.CollectedAt = metric.DateTime
			model_info.CollectedMillis = uint64(metric.Duration)
		}
		for _, temp := range nm.Metrics[idx].Temperatures {
			key := temp.SensorKey
			val, ok := temps[key]
			if !ok {
				temps[key] = MeterMinMaxAvgFloat64{
					Count: 0,
					Max:   temp.Temperature,
					Min:   temp.Temperature,
					Total: temp.Temperature,
					Avg:   0,
				}

			}
			if val.Max < temp.Temperature {
				val.Max = temp.Temperature
			}
			if val.Min > temp.Temperature {
				val.Min = temp.Temperature
			}
			if val.Total > temp.Temperature {
				val.Total += temp.Temperature
			}
			val.Count++
		}
		for _, usr := range nm.Metrics[idx].Users {
			key := fmt.Sprintf("%s_%s_%s", usr.User, usr.Host, usr.Terminal)
			_, ok := users[key]
			if !ok {
				started := time.Unix(int64(usr.Started), 0)
				users[key] = models.HostUsers{
					User:     usr.User,
					Terminal: usr.Terminal,
					Host:     usr.Host,
					Started:  started,
				}
			}
		}
	}

	tx, err := nm.db.Begin()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to start transaction")
		return
	}

	if err := tx.Create(&model_info).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to create host info record")
		tx.Rollback()
		return
	}

	model_temps := make([]models.HostTemperatures, 0)
	for key, temp := range temps {
		n := models.HostTemperatures{
			HostID:          nm.HostID,
			SensorKey:       key,
			TemperatureAvg:  temp.Total / float64(temp.Count),
			TemperatureMin:  temp.Min,
			TemperatureMax:  temp.Max,
			CollectedAt:     nm.Metrics[0].DateTime,
			CollectedMillis: uint64(nm.Metrics[0].Duration),
		}
		model_temps = append(model_temps, n)
	}

	if len(model_temps) > 0 {
		if err := tx.Create(&model_temps).Error; err != nil {
			logger.Log.Error().Err(err).Msg("failed to create host temperature records")
			tx.Rollback()
			return
		}
	}
	// model_users := make([]models.HostUsers{}, 0)

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
