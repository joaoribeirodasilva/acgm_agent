package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog/log"
)

type MeterInterface interface {
	GetThreadName() string
	Collect()
}

type Thread struct {
	meter        MeterInterface
	is_running   bool
	request_stop bool
	config       *config.Config
}

func NewThread(ti MeterInterface, conf *config.Config) *Thread {
	t := new(Thread)
	t.config = conf
	t.meter = ti
	return t
}

func (t *Thread) Start() error {
	t.request_stop = false
	log.Info().Str("namespace", "meters::thread::Start").Msgf("Starting thread %s", t.meter.GetThreadName())
	go t.run()
	return nil
}

func (t *Thread) run() {

	t.is_running = true
	log.Info().Str("namespace", "meters::thread::run").Msgf("Thread %s started", t.meter.GetThreadName())
	for {
		time.Sleep(time.Duration(t.config.Metrics.CollectInterval) * time.Millisecond)
		if t.request_stop {
			log.Info().Str("namespace", "meters::thread::run").Msgf("Stop thread %s requested", t.meter.GetThreadName())
			break
		}
		t.meter.Collect()
	}
	t.is_running = false
}

func (t *Thread) Stop() error {
	t.request_stop = true
	for t.is_running {
		time.Sleep(time.Duration(t.config.Metrics.CollectInterval/2) * time.Millisecond)
	}
	log.Info().Str("namespace", "meters::thread::Stop").Msgf("Thread %s stopped", t.meter.GetThreadName())
	return nil
}

func (t *Thread) IsRunning() bool {
	return t.is_running
}
