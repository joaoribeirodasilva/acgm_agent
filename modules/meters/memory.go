package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/logger"
	"github.com/shirou/gopsutil/mem"
)

var MEMORY_THREAD_NAME = "memory"

type Memory struct {
	DateTime      time.Time             `json:"date_time" yaml:"date_time"`
	VirtualMemory mem.VirtualMemoryStat `json:"memory" yaml:"memory"`
	SwapDevices   []mem.SwapDevice      `json:"swap_devices" yaml:"swap_devices"`
	SwapMemory    mem.SwapMemoryStat    `json:"swap_memory" yaml:"swap_memory"`
}

type MetricsMemory struct {
	Metrics     []Memory       `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func NewMetricsMemory(conf *config.Config) *MetricsMemory {
	nm := &MetricsMemory{
		config:      conf,
		init_failed: false,
		collecting:  false,
		cutting:     false,
	}
	return nm
}

func (nm *MetricsMemory) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

func (nm *MetricsMemory) IsInitFailed() bool {
	return nm.init_failed
}

func (nm *MetricsMemory) GetThreadName() string {
	return MEMORY_THREAD_NAME
}

func (nm *MetricsMemory) Collect() error {
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

	metrics := Memory{}
	metrics.DateTime = time.Now()

	if err := nm.GetVirtualMemory(&metrics); err != nil {
		return err
	}

	if err := nm.GetSwapDevices(&metrics); err != nil {
		return err
	}

	if err := nm.GetSwapMemory(&metrics); err != nil {
		return err
	}

	nm.Metrics = append(nm.Metrics, metrics)

	logger.Log.Debug().Msg("finish data collection for memory")
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("collect took %d ms", diff)
	}
	nm.collecting = false
	return nil
}

func (nm *MetricsMemory) Cut() (*[]Memory, error) {

	var start, end, diff int64
	if nm.init_failed {
		return nil, nil
	}
	if nm.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	nm.cutting = true
	for nm.collecting {
		logger.Log.Debug().Msg("Waiting for collect to finish")
		time.Sleep(50 * time.Millisecond)
	}
	logger.Log.Debug().Msg("Cutting data")

	metrics := nm.Metrics
	nm.Metrics = []Memory{}
	logger.Log.Debug().Msgf("finish data cut with %d metrics", len(metrics))

	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("Cut took %d ms", diff)
	}
	nm.cutting = false

	return &metrics, nil
}

func (nm *MetricsMemory) GetVirtualMemory(m *Memory) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	virtual_memory, err := mem.VirtualMemory()
	if err != nil || virtual_memory == nil {
		logger.Log.Error().Err(err).Msg("failed to get host virtual memory data")
		return err
	}

	m.VirtualMemory = *virtual_memory

	return nil
}

func (nm *MetricsMemory) GetSwapDevices(m *Memory) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	swap_devices, err := mem.SwapDevices()
	if err != nil || swap_devices == nil {
		logger.Log.Error().Err(err).Msg("failed to get host swap devices")
		return err
	}

	for _, swap_device := range swap_devices {
		m.SwapDevices = append(m.SwapDevices, *swap_device)
	}

	return nil
}

func (nm *MetricsMemory) GetSwapMemory(m *Memory) error {

	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}

	swap_memory, err := mem.SwapMemory()
	if err != nil || swap_memory == nil {
		logger.Log.Error().Err(err).Msg("failed to get host swap memory data")
		return err
	}

	m.SwapMemory = *swap_memory

	return nil
}

func (nm *MetricsMemory) error_init() error {
	str_err := "initialization failed, init must be run again"
	err := fmt.Errorf(str_err)
	return err
}
