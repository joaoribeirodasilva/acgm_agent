package disks

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/plugins/base"
	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

const (
	NAME                               = "disk"
	PLUGIN_METRIC_DISK_PARTITIONS_NAME = "disk.partitions"
	PLUGIN_METRIC_DISK_USAGES_NAME     = "disk.usages"
	PLUGIN_METRIC_DISK_IOS_NAME        = "disk.ios"
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

type Plugin struct {
	name    string
	db      *base.Db
	metrics []string
	data    DiskModel
	status  base.PluginStatus
}

var P Plugin

func (p *Plugin) Factory(conf *map[string]string, db *base.Db) error {
	p.db = db
	p.name = NAME
	p.metrics = []string{
		PLUGIN_METRIC_DISK_PARTITIONS_NAME,
		PLUGIN_METRIC_DISK_USAGES_NAME,
		PLUGIN_METRIC_DISK_IOS_NAME,
	}
	if err := p.data.Factory(0); err != nil {
		p.status = base.STATUS_ERROR
		return err
	}
	p.status = base.STATUS_STOPPED
	return nil
}

func (p *Plugin) Start() error {

	p.status = base.STATUS_STARTING

	return nil
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) GetName() string {
	return p.name
}

func (p *Plugin) GetMetrics() []string {
	return p.metrics
}

func (p *Plugin) GetData(metric string) (interface{}, error) {

	switch metric {
	case PLUGIN_METRIC_DISK_PARTITIONS_NAME:
		return p.data.DiskPartitions.Stats, nil
	case PLUGIN_METRIC_DISK_USAGES_NAME:
		return p.data.DiskUsages.Stats, nil
	case PLUGIN_METRIC_DISK_IOS_NAME:
		return p.data.DiskIOs.Stats, nil
	}

	return 0, fmt.Errorf("metric not found")
}
func (p *Plugin) GetDataCount(metric string) (int64, error) {

	switch metric {
	case PLUGIN_METRIC_DISK_PARTITIONS_NAME:
		return p.data.DiskPartitions.Count, nil
	case PLUGIN_METRIC_DISK_USAGES_NAME:
		return p.data.DiskUsages.Count, nil
	case PLUGIN_METRIC_DISK_IOS_NAME:
		return p.data.DiskIOs.Count, nil
	}

	return 0, fmt.Errorf("metric not found")
}
func (p *Plugin) GetStatus() base.PluginStatus {
	return p.status
}

func (p *Plugin) SetStatus(status base.PluginStatus) {
	p.status = status
}

func (p *Plugin) Polling() error {

	p.status = base.STATUS_COLLECTING

	p.GetPartitions()
	p.GetUsage()
	p.GetIO()

	p.status = base.STATUS_RUNNING

	return nil
}

// GetPartitions reads the partition data from the OS and
// stores it into a partition data array of database models.
// It returns and error if it fails to read the partitions data.
func (p *Plugin) GetPartitions() error {

	// get the metric collection start time
	time_start := time.Now()

	// get the disk partitions data
	temp_partitions, err := disk.Partitions(true)
	if err != nil {
		return err
	}

	// filter the partitions by fstype ignoring those that we don't want
	partitions := make([]disk.PartitionStat, 0)
	for _, partition := range temp_partitions {
		if slices.Contains(IgnoreFileSystems, partition.Fstype) {
			continue
		}
		partitions = append(partitions, partition)
	}

	// get the final metric collection time
	collect_duration := time.Now().UnixMilli() - time_start.UnixMilli()

	// add the partition data array to the disk partitions list
	if len(partitions) > 0 {
		p.data.DiskPartitions.AddStat(&partitions, collect_duration)
	}

	return nil
}

// GetUsage reads the partition usage data from the OS and
// stores it into a partition usage data array of database models.
// It returns and error if it fails to read the partitions usage data
// or not partition data is present.
func (p *Plugin) GetUsage() error {

	// get the metric collection start time
	time_start := time.Now()

	usages := make([]disk.UsageStat, 0)

	// we can only collect usage data if we already have partitions data
	// stored
	if len(p.data.DiskPartitions.Stats) > 0 {
		for _, partition := range p.data.DiskPartitions.Stats {
			usage, err := disk.Usage(partition.MountPoint)
			if err == nil && usage != nil {
				usages = append(usages, *usage)
			}
		}
	} else {
		return fmt.Errorf("no partitions to collect usage data from")
	}

	// get the final metric collection time
	collect_duration := time.Now().UnixMilli() - time_start.UnixMilli()

	// add the partition usage data array to the disk usage list
	if len(usages) > 0 {
		p.data.DiskUsages.AddStat(&usages, collect_duration)
	}

	return nil
}

// GetIO reads the partition io counters data from the OS and
// stores it into a partition io counters data array of database models.
// It returns and error if it fails to read the partitions io counters data
// or not partition data is present.
func (p *Plugin) GetIO() error {

	// get the metric collection start time
	time_start := time.Now()

	ios := make(map[string]disk.IOCountersStat, 0)

	// we can only collect io data if we already have partitions data
	// stored
	if len(p.data.DiskPartitions.Stats) > 0 {
		for _, partition := range p.data.DiskPartitions.Stats {
			io, err := disk.IOCounters(partition.Device)
			if err == nil && io != nil {
				for k, v := range io {
					ios[k] = v
				}
			}
		}
	} else {
		return fmt.Errorf("no partitions to collect io data from")
	}

	// get the final metric collection time
	collect_duration := time.Now().UnixMilli() - time_start.UnixMilli()

	if len(ios) > 0 {
		p.data.DiskIOs.AddStat(&ios, collect_duration)
	}

	return nil
}
