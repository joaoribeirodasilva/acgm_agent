package disks

const (
	MODULE                             = "disk"
	MODULE_METRIC_DISK_PARTITIONS_NAME = "disk.partitions"
	MODULE_METRIC_DISK_USAGES_NAME     = "disk.usages"
	MODULE_METRIC_DISK_IOS_NAME        = "disk.ios"
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

var M Meter

func (m *Meter) Factory(host_id int64, active bool, interval int64, aggregate int64) (*Meter, error) {

	m.model = &DiskModel{}
	m.name = MODULE
	m.metrics = []string{
		MODULE_METRIC_DISK_PARTITIONS_NAME,
		MODULE_METRIC_DISK_USAGES_NAME,
		MODULE_METRIC_DISK_IOS_NAME,
	}
	m.conf.host_id = host_id
	m.conf.active = active
	m.conf.interval = interval
	m.conf.aggregate = aggregate
	m.status = STATUS_STOPPED
	model := m.model.(*DiskModel)
	model.Factory(host_id)

	return m, nil
}

func (m *Meter) Start() error {
	if !m.conf.active {
		return nil
	}
	m.Polling()
	return nil
}

func (m *Meter) Stop() error {
	if !m.conf.active {
		return nil
	}
	return nil
}

func (m *Meter) GetData(meter string) (interface{}, error) {

	return nil, nil
}

func (m *Meter) GetStatus() PluginStatus {
	return m.status
}

func (m *Meter) Polling() error {
	return nil
}

// type Partition struct {
// 	Info     disk.PartitionStat  `json:"info" yaml:"info"`
// 	Usage    disk.UsageStat      `json:"usage" yaml:"usage"`
// 	Counters disk.IOCountersStat `json:"counters" yaml:"counters"`
// }

// type Partitions struct {
// 	DateTime  time.Time   `json:"date_time" yaml:"date_time"`
// 	Partition []Partition `json:"Partition" yaml:"Partition"`
// }

// type MetricsPartition struct {
// 	Metrics     []Partitions   `json:"metrics" yaml:"metrics"`
// 	config      *config.Config `json:"-" yaml:"-"`
// 	init_failed bool           `json:"-" yaml:"-"`
// 	collecting  bool           `json:"-" yaml:"-"`
// 	cutting     bool           `json:"-" yaml:"-"`
// }

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
