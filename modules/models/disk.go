package models

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/disk"
	"golang.org/x/exp/slices"
)

// Place in meter module
// const (
// 	MODULE_DISK                        = "disk"
// 	MODULE_METRIC_DISK_PARTITIONS_NAME = "disk.partitions"
// 	MODULE_METRIC_DISK_USAGES_NAME     = "disk.usages"
// 	MODULE_METRIC_DISK_IOS_NAME        = "disk.ios"
// )

type DiskModel struct {
	DiskPartitions DiskPartitions `json:"disk_partitions"`
	DiskUsages     DiskUsages     `json:"disk_usages"`
	DiskIOs        DiskIOs        `json:"disk_ios"`
}

// Disk partition data model
// For each partition only on row is generated on the
// database.
type DiskPartition struct {
	ModelCommon
	Device     string `json:"device" gorm:"column:device;type:string;size:100"`
	MountPoint string `json:"mount_point" gorm:"column:mount_point;type:string;size:255"`
	FSType     string `json:"fs_type" gorm:"column:fs_type;type:string;size:100"`
	Options    string `json:"options" gorm:"column:options;type:string"`
}

// DiskPartition database table name
func (d *DiskPartition) TableName() string {
	return "disk_partitions"
}

// Disk partitions structure
// Holds the disk partitions information and it's
// model records
type DiskPartitions struct {
	ModelListCommon
	PartitionID int64
	Stats       []*disk.PartitionStat
	Aggregated  map[string]*DiskPartition
}

// NewDiskPartitions returns a new DiskPartitions object and
// sets it into the DiskModel object receiver member with the
// same name
func (o *DiskModel) NewDiskPartitions(host_id int64) *DiskPartitions {
	o.DiskPartitions = DiskPartitions{}
	o.DiskPartitions.HostID = host_id
	o.DiskPartitions.Single = true
	return &o.DiskPartitions
}

// AddStat adds a "github.com/shirou/gopsutil/disk" PartitionStat
// object into the DiskPartitions receiver Stats array
func (o *DiskPartitions) AddStat(i []*disk.PartitionStat, collect_duration int64) {

	if len(i) == 0 {
		return
	}
	if len(o.Stats) == 0 {
		o.CollectedStart = time.Now()
	}
	o.CollectedEnd = time.Now()
	o.CollectedTotal += collect_duration
	o.CollectedAvg = float64(o.CollectedTotal) / float64(len(o.Stats))
	if o.CollectedMin > collect_duration {
		o.CollectedMin = collect_duration
	}
	if o.CollectedMax > collect_duration {
		o.CollectedMax = collect_duration
	}

	o.Stats = append(o.Stats, i...)

}

// Aggregate aggregates the several "github.com/shirou/gopsutil/disk"
// PartitionStat present in the Stats array performing any necessary
// calculations filling the DiskPartitions models inserting them into
// it's Aggregated member to be stored in the
// database.
// If the member Stats is an empty array an error is returned
func (o *DiskPartitions) Aggregate() error {

	if len(o.Stats) == 0 {
		return fmt.Errorf("no partition stats to aggregate")
	}

	mounts := make([]string, 0)
	for _, stat := range o.Stats {
		if !slices.Contains(mounts, stat.Mountpoint) {
			disk_partition := DiskPartition{
				Device:     stat.Device,
				MountPoint: stat.Mountpoint,
				FSType:     stat.Fstype,
				Options:    stat.Opts,
			}
			disk_partition.HostID = o.HostID
			disk_partition.CollectedAt = o.CollectedEnd
			mounts = append(mounts, stat.Mountpoint)
			o.Aggregated = append(o.Aggregated, &disk_partition)
		}
	}
	return nil
}

// Empty resets the DiskPartitions members
func (o *DiskPartitions) Empty() {

	o.Stats = make([]*disk.PartitionStat, 0)
	o.Aggregated = make(map[string]*DiskPartition, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
}

// GetModels returns the DiskPartitions Aggregated member that contains a
// DiskPartition array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskPartitions) GetModels() (*[]*DiskPartition, error) {
	if len(o.Aggregated) == 0 {
		return nil, fmt.Errorf("no partitions found")
	}
	return &o.Aggregated, nil
}

// GetStats returns the DiskPartitions Stats member that contains a
// "github.com/shirou/gopsutil/disk" PartitionStat array.
// If the Stats member is an empty array an error is returned
func (o *DiskPartitions) GetStats() (*[]*disk.PartitionStat, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no partitions stats found")
	}
	return &o.Stats, nil
}

// TODO: func (o *DiskPartitions) CollectStats()

func (o *DiskPartitions) FindDevicePartition(device string) (int64, error) {

	for _, partition := range o.Aggregated {
		if partition.Device == device {
			return partition.ID, nil
		}
	}
	return 0, fmt.Errorf("partition device %s not found", device)
}

func (o *DiskPartitions) FindPathPartition(path string) (int64, error) {

	for _, partition := range o.Aggregated {
		if partition.MountPoint == path {
			return partition.ID, nil
		}
	}
	return 0, fmt.Errorf("partition path %s not found", path)
}

/**
 * Disk partition usage aggregated statistics
 */
type DiskUsage struct {
	ModelCommon
	PartitionID          int64   `json:"partition_id" gorm:"column:partition_id;type:int;"`
	Path                 string  `json:"path" gorm:"column:path;type:int;"`
	Fstype               string  `json:"fstype" gorm:"column:fstype;type:int;"`
	Total                uint64  `json:"total" gorm:"column:total;type:int;"`
	Free                 uint64  `json:"free" gorm:"column:free;type:int;"`
	FreeAvg              float64 `json:"free_avg" gorm:"column:free_avg;type:float;"`
	FreeMin              uint64  `json:"free_min" gorm:"column:free_min;type:int;"`
	FreeMax              uint64  `json:"free_max" gorm:"column:free_max;type:int;"`
	Used                 uint64  `json:"used" gorm:"column:used;type:int;"`
	UsedAvg              float64 `json:"used_avg" gorm:"column:used_avg;type:float;"`
	UsedMin              uint64  `json:"used_min" gorm:"column:fused_min;type:int;"`
	UsedMax              uint64  `json:"used_max" gorm:"column:used_max;type:int;"`
	UsedPercent          float64 `json:"-" gorm:"-"`
	UsedPercentAvg       float64 `json:"used_percent_avg" gorm:"column:used_percent_avg;type:float;"`
	UsedPercentMin       float64 `json:"used_percent_min" gorm:"column:used_percent_min;type:float;"`
	UsedPercentMax       float64 `json:"used_percent_max" gorm:"column:used_percent_max;type:float;"`
	InodesTotal          uint64  `json:"-" gorm:"-"`
	InodesTotalAvg       float64 `json:"inodes_total_avg" gorm:"column:inodes_total_avg;type:float;"`
	InodesTotalMin       uint64  `json:"inodes_total_min" gorm:"column:inodes_total_min;type:int;"`
	InodesTotalMax       uint64  `json:"inodes_total_max" gorm:"column:inodes_total_max;type:int;"`
	InodesUsed           uint64  `json:"-" gorm:"-"`
	InodesUsedAvg        float64 `json:"inodes_used_avg" gorm:"column:inodes_used_avg;type:float;"`
	InodesUsedMin        uint64  `json:"inodes_used_min" gorm:"column:inodes_used_min;type:int;"`
	InodesUsedMax        uint64  `json:"inodes_used_max" gorm:"column:inodes_used_max;type:int;"`
	InodesFree           uint64  `json:"-" gorm:"-"`
	InodesFreeAvg        float64 `json:"inodes_free_avg" gorm:"column:inodes_free_avg;type:float;"`
	InodesFreeMin        uint64  `json:"inodes_free_min" gorm:"column:inodes_free_min;type:int;"`
	InodesFreeMax        uint64  `json:"inodes_free_max" gorm:"column:inodes_free_max;type:int;"`
	InodesUsedPercent    float64 `json:"-" gorm:"-"`
	InodesUsedPercentAvg float64 `json:"inodes_used_percent_avg" gorm:"column:inodes_used_percent_avg;type:float;"`
	InodesUsedPercentMin float64 `json:"inodes_used_percent_min" gorm:"column:inodes_used_percent_min;type:float;"`
	InodesUsedPercentMax float64 `json:"inodes_used_percent_max" gorm:"column:inodes_used_percent_max;type:float;"`
}

func (o *DiskUsage) TableName() string {
	return "disk_usages"
}

type DiskUsages struct {
	ModelListCommon
	Stats      []*disk.UsageStat
	Aggregated map[string]*DiskUsage
}

func (o *DiskModel) NewDiskUsages(host_id int64) *DiskUsages {

	o.DiskUsages = DiskUsages{}
	o.DiskUsages.HostID = host_id
	o.DiskUsages.Single = false
	return &o.DiskUsages
}

func (o *DiskUsages) AddStat(i []*disk.UsageStat, collect_duration int64) {

	if len(i) == 0 {
		return
	}
	if len(o.Stats) == 0 {
		o.CollectedStart = time.Now()
	}
	o.CollectedEnd = time.Now()
	o.CollectedTotal += collect_duration
	o.CollectedAvg = float64(o.CollectedTotal) / float64(len(o.Stats))
	if o.CollectedMin > collect_duration {
		o.CollectedMin = collect_duration
	}
	if o.CollectedMax > collect_duration {
		o.CollectedMax = collect_duration
	}

	o.Stats = append(o.Stats, i...)

}

func (o *DiskUsages) Aggregate() error {

	if len(o.Stats) == 0 {
		return fmt.Errorf("no usage stats to aggregate")
	}

	for _, stat := range o.Stats {
		if !slices.Contains(paths, stat.Path) {
			m := DiskUsage{}
		}
	}

	return nil
}

func (o *DiskUsages) Empty() {

	o.Stats = make([]*disk.UsageStat, 0)
	o.Aggregated = make([]*DiskUsage, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
}

// GetModels returns the DiskUsage Aggregated member that contains a
// DiskUsage array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskUsages) GetModels() (*[]*DiskUsage, error) {
	if len(o.Aggregated) == 0 {
		return nil, fmt.Errorf("no usages found")
	}
	return &o.Aggregated, nil
}

// GetStats returns the DiskUsages Stats member that contains a
// "github.com/shirou/gopsutil/disk" DiskUsages array.
// If the Stats member is an empty array an error is returned
func (o *DiskUsages) GetStats() (*[]*disk.UsageStat, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no usages stats found")
	}
	return &o.Stats, nil
}

/**
 * Disk partition IO aggregated statistics
 */
type DiskIO struct {
	ModelCommon
	PartitionID         int64   `json:"partition_id" gorm:"column:partition_id;type:int;"`
	Name                string  `json:"name" gorm:"column:name;type:string;size:255;"`
	SerialNumber        string  `json:"serial_number" gorm:"column:serial_number;type:string;size:100;"`
	Label               string  `json:"label" gorm:"column:label;type:string;size:255;"`
	ReadCount           uint64  `json:"-" gorm:"-"`
	ReadCountAvg        float64 `json:"read_count_avg" gorm:"column:read_count_avg;type:float;"`
	ReadCountMin        uint64  `json:"read_count_min" gorm:"column:read_count_min;type:int;"`
	ReadCountMax        uint64  `json:"read_count_max" gorm:"column:read_count_max;type:int;"`
	MergedReadCount     uint64  `json:"-" gorm:"-"`
	MergedReadCountAvg  float64 `json:"merged_read_read_count_avg" gorm:"column:merged_read_count_min;type:float;"`
	MergedReadCountMin  uint64  `json:"merged_read_count_min" gorm:"column:merged_read_count_min;type:int;"`
	MergedReadCountMax  uint64  `json:"merged_read_count_max" gorm:"column:merged_read_count_max;type:int;"`
	WriteCount          uint64  `json:"-" gorm:"-"`
	WriteCountAvg       float64 `json:"write_count_avg" gorm:"column:write_count_avg;type:float;"`
	WriteCountMin       uint64  `json:"write_count_min" gorm:"column:write_count_min;type:int;"`
	WriteCountMax       uint64  `json:"write_count_max" gorm:"column:write_count_max;type:int;"`
	MergedWriteCount    uint64  `json:"-" gorm:"-"`
	MergedWriteCountAvg float64 `json:"merged_write_count_avg" gorm:"column:merged_write_count_avg;type:float;"`
	MergedWriteCountMin uint64  `json:"merged_write_count_min" gorm:"column:merged_write_count_min;type:int;"`
	MergedWriteCountMax uint64  `json:"merged_write_count_max" gorm:"column:merged_write_count_max;type:int;"`
	ReadBytes           float64 `json:"-" gorm:"-"`
	ReadBytesAvg        float64 `json:"read_bytes_avg" gorm:"column:read_bytes_avg;type:float;"`
	ReadBytesMin        uint64  `json:"read_bytes_min" gorm:"column:read_bytes_min;type:int;"`
	ReadBytesMax        uint64  `json:"read_bytes_max" gorm:"column:read_bytes_max;type:int;"`
	WriteBytes          uint64  `json:"-" gorm:"-"`
	WriteBytesAvg       float64 `json:"write_bytes_avg" gorm:"column:write_bytes_avg;type:float;"`
	WriteBytesMin       uint64  `json:"write_bytes_min" gorm:"column:write_bytes_min;type:int;"`
	WriteBytesMax       uint64  `json:"write_bytes_max" gorm:"column:write_bytes_max;type:int;"`
	ReadTime            uint64  `json:"-" gorm:"-"`
	ReadTimeAvg         float64 `json:"read_time_avg" gorm:"column:read_time_avg;type:float;"`
	ReadTimeMin         uint64  `json:"read_time_min" gorm:"column:read_time_min;type:int;"`
	ReadTimeMax         uint64  `json:"read_time_max" gorm:"column:read_time_max;type:int;"`
	WriteTime           uint64  `json:"-" gorm:"-"`
	WriteTimeAvg        float64 `json:"write_time_avg" gorm:"column:write_time_avg;type:float;"`
	WriteTimeMin        uint64  `json:"write_time_min" gorm:"column:write_time_min;type:int;"`
	WriteTimeMax        uint64  `json:"write_time_max" gorm:"column:write_time_max;type:int;"`
	IopsInProgress      uint64  `json:"-" gorm:"-"`
	IopsInProgressAvg   float64 `json:"iops_in_progress_avg" gorm:"column:iops_in_progress_avg;type:float;"`
	IopsInProgressMin   uint64  `json:"iops_in_progress_min" gorm:"column:iops_in_progress_min;type:int;"`
	IopsInProgressMax   uint64  `json:"iops_in_progress_max" gorm:"column:iops_in_progress_max;type:int;"`
	IoTime              uint64  `json:"-" gorm:"-"`
	IoTimeAvg           float64 `json:"io_time_avg" gorm:"column:io_time_avg;type:float;"`
	IoTimeMin           uint64  `json:"io_time_min" gorm:"column:io_time_min;type:int;"`
	IoTimeMax           uint64  `json:"io_time_max" gorm:"column:io_time_max;type:int;"`
	WeightedIO          uint64  `json:"-" gorm:"-"`
	WeightedIOAvg       float64 `json:"weighted_io_avg" gorm:"column:weighted_io_avg;type:float;"`
	WeightedIOMin       uint64  `json:"weighted_io_min" gorm:"column:weighted_io_min;type:int;"`
	WeightedIOMax       uint64  `json:"weighted_io_max" gorm:"column:weighted_io_max;type:int;"`
}

func (o *DiskIO) TableName() string {
	return "disk_io"
}

type DiskIOs struct {
	ModelListCommon
	Stats      []*disk.IOCountersStat
	Aggregated map[string]*DiskIO
}

func (o *DiskModel) NewDiskIOs(host_id int64) *DiskIOs {

	o.DiskIOs = DiskIOs{}
	o.DiskIOs.HostID = host_id
	o.DiskIOs.Single = false
	return &o.DiskIOs
}

func (o *DiskIOs) AddStat(i []*disk.IOCountersStat, collect_duration int64) {

	if len(i) == 0 {
		return
	}
	if len(o.Stats) == 0 {
		o.CollectedStart = time.Now()
	}
	o.CollectedEnd = time.Now()
	o.CollectedTotal += collect_duration
	o.CollectedAvg = float64(o.CollectedTotal) / float64(len(o.Stats))
	if o.CollectedMin > collect_duration {
		o.CollectedMin = collect_duration
	}
	if o.CollectedMax > collect_duration {
		o.CollectedMax = collect_duration
	}

	o.Stats = append(o.Stats, i...)

}

func (o *DiskIOs) Aggregate() error {

	// r := &DiskUsage{}
	return nil
}

func (o *DiskIOs) Empty() {

	o.Stats = make([]*disk.IOCountersStat, 0)
	o.Aggregated = make([]*DiskIO, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
}

// GetModels returns the DiskIOs Aggregated member that contains a
// DiskIOs array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskIOs) GetModels() (*[]*DiskIO, error) {
	if len(o.Aggregated) == 0 {
		return nil, fmt.Errorf("no ios found")
	}
	return &o.Aggregated, nil
}

// GetStats returns the DiskUsages Stats member that contains a
// "github.com/shirou/gopsutil/disk" DiskUsages array.
// If the Stats member is an empty array an error is returned
func (o *DiskIOs) GetStats() (*[]*disk.IOCountersStat, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no is stats found")
	}
	return &o.Stats, nil
}
