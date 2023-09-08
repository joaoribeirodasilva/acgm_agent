package meters

import (
	"time"

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
	nm.status = STATUS_STARTING
	nm.EventHandler()
	return nm
}

func (m *Meters) Start() {

	m.status = STATUS_STARTING
	logger.Log.Debug().Msg("starting")

	m.Host = NewHostMeter(m.conf, m.db)
	name, err := m.Host.GetHostID()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get hostname")
		return
	}
	m.CPU = NewCPUMeter(name, m.conf, m.db)
	go m.loop()

}

func (m *Meters) Stop() {

	m.status = STATUS_STOPPING
	evnt.Signal("meter.meter.stop")

	// Wait for all other meters to stop
	not_stopped := true
	for not_stopped {
		not_stopped = false
		time.Sleep(50 * time.Millisecond)
		// fmt.Printf("Hosts status: %d\n", m.Host.GetStatus())
		if m.Host.GetStatus() != STATUS_STOPPED {
			not_stopped = true
		}
		// fmt.Printf("CPU status: %d\n", m.Host.GetStatus())
		if m.CPU.GetStatus() != STATUS_STOPPED {
			not_stopped = true
		}
		if !not_stopped {
			logger.Log.Debug().Msg("all meters stopped")
			break
		}
	}

	m.status = STATUS_STOPPED
	evnt.Signal("meter.meter.stopped")

	m.Host = nil
	m.CPU = nil

	logger.Log.Debug().Msg("stopped")
}

func (m *Meters) loop() {

	m.status = STATUS_STARTED
	evnt.Signal("meter.meter.start")

	// Wait for all other meters to start
	for {
		not_started := false
		time.Sleep(50 * time.Millisecond)
		if m.Host.GetStatus() != STATUS_WAITING {
			not_started = true
		}
		if m.CPU.GetStatus() != STATUS_WAITING {
			not_started = true
		}
		if !not_started {
			break
		}
	}

	logger.Log.Debug().Msg("started")

	for m.status == STATUS_WAITING || m.status == STATUS_COLLECTING || m.status == STATUS_STARTED {
		m.status = STATUS_WAITING
		time.Sleep(1000 * time.Millisecond)
		evnt.Signal("meter.meter.collect")
	}

	logger.Log.Debug().Msg("stopping")
}

func (nm *Meters) EventHandler() {

	logger.Log.Debug().Msg("event handler start")
	events := evnt.Listen("meter.")

	go func() {
		for nm.status != STATUS_STOPPED {
			for event := range events {
				logger.Log.Debug().Msgf("received event: %s", event.Tag)
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
				case "meter.meter.stop":
					nm.status = STATUS_STOPPING
				case "meter.meter.stopped":
					close(events)
				}
			}
		}
		logger.Log.Debug().Msg("event handler stopped")
	}()

}
