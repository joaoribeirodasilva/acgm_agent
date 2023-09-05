package meters

import (
	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
	"biqx.com.br/acgm_agent/modules/logger"
	evnt "github.com/jonhoo/go-events"
)

type Meters struct {
	MeterControl
	Host *HostMeter `json:"host" yaml:"host"`
	CPU  *CPUMeter  `json:"cpu" yaml:"cpu"`
	// Memory     *Memory     `json:"memory" yaml:"memory"`
	// Net        *Net        `json:"network" yaml:"network"`
	// Partitions *Partitions `json:"disk" yaml:"disk"`
	// Processes  *Processes  `json:"processes" yaml:"processes"`
}

func NewMeters(conf *config.Config, db *database.Db) *Meters {
	nm := &Meters{}
	nm.conf = conf
	nm.db = db
	go nm.EventHandler()
	return nm
}

func (m *Meters) Start() {
	m.Host = NewHostMeter(m.conf, m.db)
	name, err := m.Host.GetHostID()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get hostname")
		return
	}
	m.CPU = NewCPUMeter(name, m.conf, m.db)

	evnt.Signal("meter.meter.start")

}

func (m *Meters) Stop() {
	evnt.Signal("meter.meter.stop")
}

func (nm *Meters) EventHandler() {
	events := evnt.Listen("meter.changed.")
	for event := range events {
		switch event.Tag {
		case "meter.changed.host":
			//
		case "meter.changed.cpu":
			//
		case "meter.changed.memory":
			//
		case "meter.changed.network":
			//
		case "meter.changed.disk":
			//
		case "meter.changed.process":
			//
		}
		if nm.status == STATUS_STOPPED {
			break
		}
	}
}
