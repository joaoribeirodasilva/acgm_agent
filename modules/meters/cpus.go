package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/logger"
	"github.com/shirou/gopsutil/cpu"
)

const CPUS_THREAD_NAME = "cpus"

// Data fo a single host CPU
type CPU struct {
	// CPU information
	Info cpu.InfoStat `json:"info" yaml:"info"`
	// CPU usage percent
	Percent float64 `json:"percent" yaml:"percent"`
	// CPU Times
	Times cpu.TimesStat `json:"times" yaml:"times"`
	// Threads per CPU
	Threads int `json:"threads" yaml:"threads"`
}

type CPUs struct {
	DateTime time.Time `json:"date_time" yaml:"date_time"`
	// Number of physical cores
	Physical int `json:"physical" yaml:"physical"`
	// Number of logical cores
	Logical int `json:"logical" yaml:"logical"`
	// Usage percent for all cores
	Percent float64 `json:"percent" yaml:"percent"`
	// Times for all cores
	Times cpu.TimesStat `json:"times" yaml:"times"`
	// CPU core list
	CPUs []CPU `json:"cpus" yaml:"cpus"`
}

type MetricsCPU struct {
	// Name of this thread - CPUs
	Metrics     []CPUs         `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func NewMetricsCPU(conf *config.Config) *MetricsCPU {
	nm := &MetricsCPU{
		config:      conf,
		init_failed: false,
		collecting:  false,
		cutting:     false,
	}
	return nm

}

// Initializes the CPUs list. This function runs only once at startup because the CPUs don't change during execution
func (nm *MetricsCPU) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

// Initialization function failed
func (nm *MetricsCPU) IsInitFailed() bool {
	return nm.init_failed
}

func (nm *MetricsCPU) GetThreadName() string {
	return CPUS_THREAD_NAME
}

func (nm *MetricsCPU) Collect() error {
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

	metric := CPUs{}
	metric.DateTime = time.Now()

	if err := nm.GetCpusCount(&metric); err != nil {
		return err
	}
	if err := nm.GetCpusInfo(&metric); err != nil {
		return err
	}
	if err := nm.GetCpusTimesTotal(&metric); err != nil {
		return err
	}
	if err := nm.GetCpusPercentTotal(&metric); err != nil {
		return err
	}
	if err := nm.GetCpusTimes(&metric); err != nil {
		return err
	}
	if err := nm.GetCpusPercent(&metric); err != nil {
		return err
	}

	nm.Metrics = append(nm.Metrics, metric)

	logger.Log.Debug().Msgf("finish data collection for %d metrics", len(nm.Metrics))
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("collect took %d ms", diff)
	}
	nm.collecting = false
	return nil
}

func (nm *MetricsCPU) Cut() (*[]CPUs, error) {
	var start, end, diff int64
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return nil, err
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
	nm.Metrics = []CPUs{}
	logger.Log.Debug().Msgf("finish data cut with %d metrics", len(metrics))
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("cut took %d ms", diff)
	}
	nm.cutting = false
	return &metrics, nil
}

func (nm *MetricsCPU) GetCpusCount(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	physical, err := cpu.Counts(false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host physical CPU count")
		return err
	}
	c.Physical = physical

	logical, err := cpu.Counts(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host logical CPU count")
		return err
	}
	c.Logical = logical
	return nil
}

func (nm *MetricsCPU) GetCpusInfo(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	info, err := cpu.Info()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host CPUs information")
		return err
	}

	c.CPUs = []CPU{}
	for _, item := range info {
		ncpu := CPU{
			Info:    item,
			Threads: c.Logical / c.Physical,
		}
		c.CPUs = append(c.CPUs, ncpu)
	}
	return nil
}

func (nm *MetricsCPU) GetCpusTimesTotal(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	total, err := cpu.Times(false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host CPUs total times")
		return err
	}

	c.Times = total[0]

	return nil
}

func (nm *MetricsCPU) GetCpusTimes(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	times, err := cpu.Times(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host CPUs times")
	}

	for idx, item := range times {
		c.CPUs[idx].Times = item
	}

	return nil
}

func (nm *MetricsCPU) GetCpusPercentTotal(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	percent, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get host CPUs percent total")
	}

	c.Percent = percent[0]

	return nil
}

func (nm *MetricsCPU) GetCpusPercent(c *CPUs) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	percent, err := cpu.Percent(0, true)
	if err != nil {
		logger.Log.Error().Str("namespace", "meters::cpus::GetCpusPercent").Err(err).Msg("Failed to get host CPUs percent")
	}

	for idx, item := range percent {
		c.CPUs[idx].Percent = item
	}

	return nil
}

func (nm *MetricsCPU) error_init() error {
	str_err := "initialization failed, init must be run again"
	err := fmt.Errorf(str_err)
	return err
}
