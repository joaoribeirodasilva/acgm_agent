package meters

import (
	"errors"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/converters"
)

type Meter struct {
	StartTime      time.Time
	EndTime        time.Time
	Interval       time.Duration
	Processes      []Process
	running        bool
	stop_requested bool
}

func New() *Meter {

	m := new(Meter)
	return m
}

func (m *Meter) Start(conf *config.Config) error {

	interval, err := converters.StringToInterval(conf.Metrics.CollectInterval)
	if err != nil {
		return err
	}

	m.Interval = *interval
	m.stop_requested = false

	go m.loop()

	return nil
}

func (m *Meter) Stop() error {

	m.stop_requested = true

	for {
		attempts := 0

		if !m.running {
			break
		}

		time.Sleep(500 * time.Millisecond)
		attempts++
		if attempts == 4 {
			return errors.New("ERROR: Process meters timeout exiting")
		}
	}

	return nil
}

func (m *Meter) loop() {

	m.running = true

	for {
		if m.stop_requested {
			break
		}
		// List all processes
		// Loop through processes
		// Filter process by path (as in config)
		// Get several process metrics
		// Store process metrics in the database
	}

	m.stop_requested = false
	m.running = false
}
