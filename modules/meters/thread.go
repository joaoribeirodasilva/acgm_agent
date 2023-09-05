package meters

// import (
// 	"time"

// 	"biqx.com.br/acgm_agent/modules/config"
// 	"biqx.com.br/acgm_agent/modules/logger"
// )

// type MeterInterface interface {
// 	Init() error
// 	IsInitFailed() bool
// 	GetThreadName() string
// 	Collect() error
// }

// type Thread struct {
// 	meter        MeterInterface
// 	is_running   bool
// 	request_stop bool
// 	config       *config.Config
// }

// func NewThread(ti MeterInterface, conf *config.Config) *Thread {
// 	t := new(Thread)
// 	t.config = conf
// 	t.meter = ti
// 	return t
// }

// func (t *Thread) Start() error {
// 	t.request_stop = false
// 	logger.Log.Info().Msgf("Starting thread %s", t.meter.GetThreadName())
// 	go t.run()
// 	return nil
// }

// func (t *Thread) run() {

// 	t.is_running = true
// 	t.meter.Init()
// 	if t.meter.IsInitFailed() {
// 		logger.Log.Error().Msgf("Thread %s initialization failed... aborting thread", t.meter.GetThreadName())
// 		return
// 	}
// 	logger.Log.Info().Msgf("Thread %s started", t.meter.GetThreadName())
// 	for {
// 		time.Sleep(time.Duration(t.config.Metrics.CollectInterval) * time.Millisecond)
// 		if t.request_stop {
// 			logger.Log.Info().Msgf("Stop thread %s requested", t.meter.GetThreadName())
// 			break
// 		}
// 		t.meter.Collect()
// 	}
// 	t.is_running = false
// }

// func (t *Thread) Stop() error {
// 	t.request_stop = true
// 	for t.is_running {
// 		time.Sleep(time.Duration(t.config.Metrics.CollectInterval/2) * time.Millisecond)
// 	}
// 	logger.Log.Info().Msgf("Thread %s stopped", t.meter.GetThreadName())
// 	return nil
// }

// func (t *Thread) IsRunning() bool {
// 	return t.is_running
// }
