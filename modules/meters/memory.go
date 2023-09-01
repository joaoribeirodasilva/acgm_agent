package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/mem"
)

var MEMORY_THREAD_NAME = "memory"

type Memory struct {
	VirtualMemory   mem.VirtualMemoryStat   `json:"memory" yaml:"memory"`
	VirtualMemoryEx mem.VirtualMemoryExStat `json:"memory_ex" yaml:"memory_ex"`
	SwapDevices     []mem.SwapDevice        `json:"swap_devices" yaml:"swap_devices"`
	SwapMemory      mem.SwapMemoryStat      `json:"swap_memory" yaml:"swap_memory"`
}

type MetricsMemory struct {
	Metrics     []Memory       `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func (nm *MetricsMemory) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

func (m *Memory) IsInitFailed() bool {
	return m.init_failed
}

func (m *Memory) GetThreadName() string {
	return MEMORY_THREAD_NAME
}

func (m *Memory) Collect() {
	var start, end, diff int64
	if m.init_failed {
		return
	}
	if m.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	m.collecting = true
	for m.cutting {
		log.Debug().Str("namespace", "meters::partitions::Collect").Msg("Waiting for cut to finish")
		time.Sleep(50 * time.Millisecond)
	}
	log.Debug().Str("namespace", "meters::partitions::Collect").Msg("Collecting data")

	// TODO: Code

	log.Debug().Str("namespace", "meters::partitions::Collect").Msg("Finish data collection for memory")
	if m.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		log.Debug().Str("namespace", "meters::partitions::Collect").Msgf("Collect took %d ms", diff)
	}
	m.collecting = false
}

func (m *Memory) Cut() *[]Memory {
	var start, end, diff int64
	if m.init_failed {
		return nil
	}
	if m.config.Settings.Debug {
		start = time.Now().UnixNano() / int64(time.Millisecond)
	}
	m.cutting = true
	for m.collecting {
		log.Debug().Str("namespace", "meters::partitions::Cut").Msg("Waiting for collect to finish")
		time.Sleep(50 * time.Millisecond)
	}
	log.Debug().Str("namespace", "meters::partitions::Cut").Msg("Cutting data")

	// TODO: Code here
	nmemory := []Memory{}
	log.Debug().Str("namespace", "meters::partitions::Cut").Msgf("Finish data cut with %d partitions", len(nmemory))
	if m.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		log.Debug().Str("namespace", "meters::partitions::Cut").Msgf("Cut took %d ms", diff)
	}
	m.cutting = false
	return &nmemory
}
