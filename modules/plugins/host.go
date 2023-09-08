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
	nm.status = STATUS_STARTING
	nm.db = db
	nm.EventHandler()
	return nm
}

func (nm *HostMeter) GetHostID() (string, error) {

	host_id, err := host.HostID()
	nm.HostID = host_id

	return host_id, err
}

func (nm *HostMeter) Start() {

	logger.Log.Debug().Msg("starting")
	nm.FireEvent(STATUS_STARTING)

	// TODO: Validate tables

	logger.Log.Debug().Msg("started")
	nm.FireEvent(STATUS_STARTED)

	logger.Log.Debug().Msg("waiting")
	nm.FireEvent(STATUS_WAITING)

}

func (nm *HostMeter) Stop() {

	logger.Log.Debug().Msg("waiting for running operations")
	for nm.status == STATUS_COLLECTING || nm.status == STATUS_AGGREGATING {
		time.Sleep(50 * time.Millisecond)
	}

	logger.Log.Debug().Msg("stopping")
	nm.FireEvent(STATUS_STOPPING)
	nm.Aggregate()

	logger.Log.Debug().Msg("stopped")
	nm.status = STATUS_STOPPED
	evnt.Signal("meter.host.stopped")

}

func (nm *HostMeter) GetStatus() MeterStatus {
	return nm.status
}

func (nm *HostMeter) GetName() string {
	return METRIC_HOST_NAME
}

func (nm *HostMeter) Collect() {

	logger.Log.Debug().Msg("checking for collect")
	if nm.status != STATUS_WAITING {
		return
	}

	logger.Log.Debug().Msg("collecting")
	nm.FireEvent(STATUS_COLLECTING)
	metric := HostMetric{}

	metric.DateTime = time.Now()

	info, err := host.Info()
	if err != nil || info == nil {
		logger.Log.Error().Err(err).Msg("failed to get host information")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Host = *info
	}

	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host temperatures")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Temperatures = temps
	}

	users, err := host.Users()
	if err != nil || len(users) == 0 {
		logger.Log.Error().Err(err).Msg("failed to get host users")
		logger.Log.Debug().Msg("collect error")
	} else {
		metric.Users = users
	}

	nm.AppendMetric(&metric)
	logger.Log.Debug().Msg("collect success")
}

func (nm *HostMeter) AppendMetric(metric *HostMetric) {
	metric.Duration = (time.Now().UnixNano() / int64(time.Millisecond)) - (metric.DateTime.UnixNano() / int64(time.Millisecond))
	nm.Metrics = append(nm.Metrics, *metric)
	nm.Aggregate()
}

func (nm *HostMeter) Aggregate() {

	logger.Log.Debug().Msg("checking for aggregation")
	if len(nm.Metrics) == 0 || (len(nm.Metrics) < nm.conf.Metrics.Cpu.Aggregate && nm.status != STATUS_STOPPING) {
		logger.Log.Debug().Msg("waiting")
		nm.FireEvent(STATUS_WAITING)
		return
	}

	nm.FireEvent(STATUS_AGGREGATING)
	logger.Log.Debug().Msg("aggregating")

	model_info := models.HostInfo{}
	temps := make(map[string]*MeterMinMaxAvgFloat64)
	users := make(map[string]*models.HostUsers)

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
			_, ok := temps[key]
			if !ok {
				temps[key] = &MeterMinMaxAvgFloat64{
					Count: 0,
					Max:   0,
					Min:   0,
					Total: 0,
					Avg:   0,
				}
			}

			if temps[key].Max < temp.Temperature {
				temps[key].Max = temp.Temperature
			}
			if temps[key].Min > temp.Temperature {
				temps[key].Min = temp.Temperature
			}
			temps[key].Total += temp.Temperature

			temps[key].Count++
			temps[key].Avg = float64(temps[key].Total) / float64(temps[key].Count)
			// fmt.Printf("Average %s: %.2f\n", key, temps[key].Avg)
		}

		for _, usr := range nm.Metrics[idx].Users {
			key := fmt.Sprintf("%s_%s_%s", usr.User, usr.Host, usr.Terminal)
			_, ok := users[key]
			if !ok {
				started := time.Unix(int64(usr.Started), 0)
				users[key] = &models.HostUsers{
					User:     usr.User,
					Terminal: usr.Terminal,
					Host:     usr.Host,
					Started:  started,
				}
			}
		}

		// TODO: avg duration
	}

	tx, err := nm.db.Begin()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to start transaction")
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}

	if err := tx.Create(&model_info).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to create host info record")
		nm.db.Rollback(tx)
		logger.Log.Debug().Msg("aggregated error")
		nm.FireEvent(STATUS_WAITING)
		return
	}
	logger.Log.Debug().Msgf("inserted %d host information", tx.RowsAffected)

	model_temps := make([]models.HostTemperatures, 0)
	for key, temp := range temps {
		n := models.HostTemperatures{
			HostID:          nm.HostID,
			SensorKey:       key,
			TemperatureAvg:  temp.Avg,
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
			nm.db.Rollback(tx)
			logger.Log.Debug().Msg("aggregated error")
			nm.FireEvent(STATUS_WAITING)
			return
		}
		logger.Log.Debug().Msgf("inserted %d host temperatures", tx.RowsAffected)
	}

	model_users := make([]models.HostUsers, 0)
	for _, usr := range users {
		n := models.HostUsers{
			HostID:          nm.HostID,
			User:            usr.User,
			Terminal:        usr.Terminal,
			Host:            usr.Host,
			Started:         usr.Started,
			CollectedAt:     nm.Metrics[0].DateTime,
			CollectedMillis: uint64(nm.Metrics[0].Duration),
		}
		model_users = append(model_users, n)
	}

	if len(model_users) > 0 {
		if err := tx.Create(&model_users).Error; err != nil {
			logger.Log.Error().Err(err).Msg("failed to create host users records")
			nm.db.Rollback(tx)
			logger.Log.Debug().Msg("aggregated error")
			nm.FireEvent(STATUS_WAITING)
			return
		}
		logger.Log.Debug().Msgf("inserted %d host users", tx.RowsAffected)
	}

	nm.db.Commit(tx)

	nm.Metrics = make([]HostMetric, 0)

	logger.Log.Debug().Msg("aggregated success")
	logger.Log.Debug().Msg("waiting")
	nm.FireEvent(STATUS_WAITING)
}

func (nm *HostMeter) FireEvent(status MeterStatus) {

	nm.status = status
	evnt.Signal("meter.changed.host")
}

func (nm *HostMeter) EventHandler() {

	logger.Log.Debug().Msg("event handler start")
	events := evnt.Listen("meter.")

	go func() {
		for nm.status != STATUS_STOPPED {
			logger.Log.Debug().Msg("listening")
			for event := range events {
				logger.Log.Debug().Msgf("received event: %s", event.Tag)
				switch event.Tag {
				case "meter.meter.start":
					nm.Start()
				case "meter.meter.collect":
					nm.Collect()
				case "meter.meter.stop":
					nm.Stop()
				case "meter.host.stopped":
					close(events)
				}
			}
		}
		logger.Log.Debug().Msg("event handler stopped")
	}()

}
