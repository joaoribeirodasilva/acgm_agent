package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/logger"
	"github.com/shirou/gopsutil/load"
)

var LOAD_THREAD_NAME = "load"

type Load struct {
	DateTime time.Time     `json:"date_time" yaml:"date_time"`
	Average  load.AvgStat  `json:"average" yaml:"average"`
	Misc     load.MiscStat `json:"misc" yaml:"misc"`
}

type MetricsLoad struct {
	Metrics     []Load         `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func NewMetricsLoad(conf *config.Config) *MetricsLoad {
	nm := &MetricsLoad{
		config:      conf,
		init_failed: false,
		collecting:  false,
		cutting:     false,
	}
	return nm
}

func (nm *MetricsLoad) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

func (nm *MetricsLoad) IsInitFailed() bool {
	return nm.init_failed
}

func (nm *MetricsLoad) GetThreadName() string {
	return LOAD_THREAD_NAME
}

func (nm *MetricsLoad) Collect() error {
	var start, end, diff int64
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	if nm.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	nm.collecting = true
	for nm.cutting {
		logger.Log.Debug().Msg("waiting for cut to finish")
		time.Sleep(50 * time.Millisecond)
	}
	logger.Log.Debug().Msg("collecting data")

	metrics := Load{}
	metrics.DateTime = time.Now()

	if err := nm.GetAverage(&metrics); err != nil {
		return err
	}

	if err := nm.GetMisc(&metrics); err != nil {
		return err
	}

	nm.Metrics = append(nm.Metrics, metrics)

	logger.Log.Debug().Msg("finish data collection for network")
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("collect took %d ms", diff)
	}
	nm.collecting = false
	return nil
}

func (nm *MetricsLoad) Cut() (*[]Load, error) {

	var start, end, diff int64
	if nm.init_failed {
		return nil, nil
	}
	if nm.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	nm.cutting = true
	for nm.collecting {
		logger.Log.Debug().Msg("waiting for collect to finish")
		time.Sleep(50 * time.Millisecond)
	}
	logger.Log.Debug().Msg("cutting data")

	metrics := nm.Metrics
	nm.Metrics = []Load{}
	logger.Log.Debug().Msgf("finish data cut with %d metrics", len(metrics))

	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("cut took %d ms", diff)
	}
	nm.cutting = false

	return &metrics, nil
}

func (nm *MetricsLoad) GetAverage(l *Load) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	average, err := load.Avg()
	if err != nil || average == nil {
		logger.Log.Error().Err(err).Msg("failed to get host load average data")
		return err
	}

	l.Average = *average

	// fmt.Printf("Average: %+v\n", l.Average)

	return nil
}

func (nm *MetricsLoad) GetMisc(l *Load) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	misc, err := load.Misc()
	if err != nil || misc == nil {
		logger.Log.Error().Err(err).Msg("failed to get host load miscellaneous data")
		return err
	}

	l.Misc = *misc

	// fmt.Printf("Average: %+v\n", l.Misc)

	return nil
}

func (nm *MetricsLoad) error_init() error {
	str_err := "initialization failed, init must be run again"
	err := fmt.Errorf(str_err)
	return err
}
