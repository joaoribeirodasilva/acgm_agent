package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/shirou/gopsutil/disk"
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

var PARTITIONS_THREAD_NAME = "disk"

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

// func (nm *MetricsPartition) GetPartitions(p *Partitions) error {
// 	if nm.init_failed {
// 		err := nm.error_init()
// 		logger.Log.Error().Err(err)
// 		return err
// 	}
// 	partitions, err := disk.Partitions(true)
// 	if err != nil {
// 		logger.Log.Error().Err(err).Msg("failed to get host partitions")
// 		return err
// 	}

// 	p.Partition = []Partition{}

// 	for _, item := range partitions {

// 		if slices.Contains(IgnoreFileSystems, item.Fstype) {
// 			continue
// 		}

// 		npartition := Partition{
// 			Info: item,
// 		}
// 		p.Partition = append(p.Partition, npartition)
// 	}

// 	return nil
// }

// func (nm *MetricsPartition) GetUsage(p *Partitions) error {
// 	if nm.init_failed {
// 		err := nm.error_init()
// 		logger.Log.Error().Err(err)
// 		return err
// 	}
// 	for _, partition := range p.Partition {
// 		usage, err := disk.Usage(partition.Info.Mountpoint)
// 		if err != nil || usage == nil {
// 			logger.Log.Error().Err(err).Msgf("failed to get host partition %s usage", partition.Info.Mountpoint)
// 			return err
// 		} else {
// 			partition.Usage = *usage
// 		}
// 	}

// 	return nil
// }

// func (nm *MetricsPartition) GetCounters(p *Partitions) error {
// 	if nm.init_failed {
// 		err := nm.error_init()
// 		logger.Log.Error().Err(err)
// 		return err
// 	}
// 	for _, partition := range p.Partition {
// 		counters, err := disk.IOCounters(partition.Info.Device)
// 		if err != nil || len(counters) == 0 {
// 			logger.Log.Error().Err(err).Msgf("Failed to get host device %s counters", partition.Info.Device)
// 			return err
// 		}

// 		for key := range counters {
// 			partition.Counters = counters[key]
// 		}
// 	}
// 	return nil
// }
