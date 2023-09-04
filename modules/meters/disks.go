package meters

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/logger"
	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

var IgnoreFileSystems = []string{
	"",
	"squashfs",
	"binfmt_misc",
	"tmpfs",
	"sysfs",
	"proc",
	"mqueue",
	"hugetlbfs",
	"debugfs",
	"rpc_pipefs",
	"devtmpfs",
	"devpts",
	"securityfs",
	"pstore",
	"efivarfs",
	"bpf",
	"autofs",
	"tracefs",
	"ramfs",
	"fuse.portal",
	"nsfs",
	"cgroup2",
	"configfs",
	"fusectl",
}

var PARTITIONS_THREAD_NAME = "partitions"

type Partition struct {
	Info     disk.PartitionStat  `json:"info" yaml:"info"`
	Usage    disk.UsageStat      `json:"usage" yaml:"usage"`
	Counters disk.IOCountersStat `json:"counters" yaml:"counters"`
}

type Partitions struct {
	DateTime  time.Time   `json:"date_time" yaml:"date_time"`
	Partition []Partition `json:"Partition" yaml:"Partition"`
}

type MetricsPartition struct {
	Metrics     []Partitions   `json:"metrics" yaml:"metrics"`
	config      *config.Config `json:"-" yaml:"-"`
	init_failed bool           `json:"-" yaml:"-"`
	collecting  bool           `json:"-" yaml:"-"`
	cutting     bool           `json:"-" yaml:"-"`
}

func NewMetricsPartition(conf *config.Config) *MetricsPartition {
	nm := &MetricsPartition{
		config:      conf,
		init_failed: false,
		collecting:  false,
		cutting:     false,
	}
	return nm
}

func (nm *MetricsPartition) Init() error {

	nm.init_failed = false

	if err := nm.Collect(); err != nil {
		nm.init_failed = true
		return err
	}

	return nil
}

func (nm *MetricsPartition) IsInitFailed() bool {
	return nm.init_failed
}

func (nm *MetricsPartition) GetThreadName() string {
	return PARTITIONS_THREAD_NAME
}

func (nm *MetricsPartition) Collect() error {
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

	metric := Partitions{}
	metric.DateTime = time.Now()

	if err := nm.GetPartitions(&metric); err != nil {
		return err
	}
	if err := nm.GetUsage(&metric); err != nil {
		return err
	}
	if err := nm.GetCounters(&metric); err != nil {
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

func (nm *MetricsPartition) Cut() (*[]Partitions, error) {
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
	logger.Log.Debug().Msg("Cutting data")
	metrics := nm.Metrics
	nm.Metrics = []Partitions{}

	logger.Log.Debug().Msgf("finish data cut with %d metrics", len(metrics))
	if nm.config.Settings.Debug {
		end = time.Now().UnixNano() / int64(time.Millisecond)
		diff = end - start
		logger.Log.Debug().Msgf("Cut took %d ms", diff)
	}
	nm.cutting = false
	return &metrics, nil
}

func (nm *MetricsPartition) GetPartitions(p *Partitions) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	partitions, err := disk.Partitions(true)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get host partitions")
		return err
	}

	p.Partition = []Partition{}

	for _, item := range partitions {

		if slices.Contains(IgnoreFileSystems, item.Fstype) {
			continue
		}

		npartition := Partition{
			Info: item,
		}
		p.Partition = append(p.Partition, npartition)
	}

	return nil
}

func (nm *MetricsPartition) GetUsage(p *Partitions) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	for _, partition := range p.Partition {
		usage, err := disk.Usage(partition.Info.Mountpoint)
		if err != nil || usage == nil {
			logger.Log.Error().Err(err).Msgf("failed to get host partition %s usage", partition.Info.Mountpoint)
			return err
		} else {
			partition.Usage = *usage
		}
	}

	return nil
}

func (nm *MetricsPartition) GetCounters(p *Partitions) error {
	if nm.init_failed {
		err := nm.error_init()
		logger.Log.Error().Err(err)
		return err
	}
	for _, partition := range p.Partition {
		counters, err := disk.IOCounters(partition.Info.Device)
		if err != nil || len(counters) == 0 {
			logger.Log.Error().Err(err).Msgf("Failed to get host device %s counters", partition.Info.Device)
			return err
		}

		for key := range counters {
			partition.Counters = counters[key]
		}
	}
	return nil
}

func (nm *MetricsPartition) error_init() error {
	str_err := "initialization failed, init must be run again"
	err := fmt.Errorf(str_err)
	return err
}
