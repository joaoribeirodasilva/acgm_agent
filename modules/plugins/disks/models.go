package disks

import (
	"fmt"
	"time"

	"biqx.com.br/acgm_agent/modules/plugins/base"
	"github.com/shirou/gopsutil/disk"
)

type DiskModel struct {
	DiskPartitions *DiskPartitions `json:"disk_partitions"`
	DiskUsages     *DiskUsages     `json:"disk_usages"`
	DiskIOs        *DiskIOs        `json:"disk_ios"`
}

func (o *DiskModel) Factory(host_id int64) {
	o.NewDiskPartitions(host_id)
	o.NewDiskUsages(host_id)
	o.NewDiskIOs(host_id)
}

/***************************************************************************************************
 * Disk partitions
 */

// Disk partition data model
// For each partition only on row is generated on the
// database.
type DiskPartition struct {
	base.ModelCommon
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
	base.ModelListCommon
	PartitionID int64
	Stats       map[string]*DiskPartition
}

// NewDiskPartitions returns a new DiskPartitions object and
// sets it into the DiskModel object receiver member with the
// same name
func (o *DiskModel) NewDiskPartitions(host_id int64) *DiskPartitions {
	o.DiskPartitions = &DiskPartitions{}
	o.DiskPartitions.HostID = host_id
	return o.DiskPartitions
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

	for _, stat := range i {
		if _, ok := o.Stats[stat.Device]; !ok {
			o.Stats[stat.Device] = &DiskPartition{}
		}
		o.Stats[stat.Device].Device = stat.Device
		o.Stats[stat.Device].HostID = o.HostID
		o.Stats[stat.Device].FSType = stat.Fstype
		o.Stats[stat.Device].MountPoint = stat.Mountpoint
		o.Stats[stat.Device].Options = stat.Opts
		o.Stats[stat.Device].CollectedAt = o.CollectedEnd
	}
}

// Empty resets the DiskPartitions members
func (o *DiskPartitions) Empty() {

	o.Stats = make(map[string]*DiskPartition, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
	o.Count = 0
}

// GetModels returns the DiskPartitions Aggregated member that contains a
// DiskPartition array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskPartitions) GetModels() (*[]*DiskPartition, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no partitions found")
	}
	rows := make([]*DiskPartition, 0)
	for _, stat := range o.Stats {
		rows = append(rows, stat)
	}
	return &rows, nil
}

// TODO: func (o *DiskPartitions) MigrateUp()
// TODO: func (o *DiskPartitions) MigrateDown()
// TODO: func (o *DiskPartitions) CollectStatsTimes()

func (o *DiskPartitions) FindDevicePartition(device string) (int64, error) {

	stat, ok := o.Stats[device]
	if !ok {
		return 0, fmt.Errorf("partition device %s not found", device)
	}
	return stat.ID, nil
}

func (o *DiskPartitions) FindPathPartition(path string) (int64, error) {

	for _, stat := range o.Stats {
		if stat.MountPoint == path {
			return stat.ID, nil
		}
	}
	return 0, fmt.Errorf("partition path %s not found", path)
}

/***************************************************************************************************
 * Disk usages
 */
type DiskUsage struct {
	base.ModelCommon
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
	base.ModelListCommon
	Stats map[string]*DiskUsage
}

func (o *DiskModel) NewDiskUsages(host_id int64) *DiskUsages {

	o.DiskUsages = &DiskUsages{}
	o.DiskUsages.HostID = host_id
	return o.DiskUsages
}

func (o *DiskUsages) AddStat(i []*disk.UsageStat, collect_duration int64) {

	if len(i) == 0 {
		return
	}
	if len(o.Stats) == 0 {
		o.CollectedStart = time.Now()
	}
	o.Count++
	o.CollectedEnd = time.Now()
	o.CollectedTotal += collect_duration
	o.CollectedAvg = float64(o.CollectedTotal) / float64(o.Count)
	if o.CollectedMin > collect_duration {
		o.CollectedMin = collect_duration
	}
	if o.CollectedMax > collect_duration {
		o.CollectedMax = collect_duration
	}

	for _, stat := range i {
		_, ok := o.Stats[stat.Path]
		if !ok {
			o.Stats[stat.Path] = &DiskUsage{}
		}
		o.Stats[stat.Path].Count++
		o.Stats[stat.Path].HostID = o.HostID
		o.Stats[stat.Path].CollectedAt = o.CollectedEnd
		o.Stats[stat.Path].Path = stat.Path
		o.Stats[stat.Path].Fstype = stat.Fstype
		o.Stats[stat.Path].Total = stat.Total
		o.Stats[stat.Path].Free += stat.Free
		o.Stats[stat.Path].FreeAvg = float64(o.Stats[stat.Path].Free) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].FreeMin > stat.Free {
			o.Stats[stat.Path].FreeMin = stat.Free
		}
		if o.Stats[stat.Path].FreeMax < stat.Free {
			o.Stats[stat.Path].FreeMax = stat.Free
		}
		o.Stats[stat.Path].Used += stat.Used
		o.Stats[stat.Path].UsedAvg = float64(o.Stats[stat.Path].Used) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].UsedMin > stat.Used {
			o.Stats[stat.Path].UsedMin = stat.Used
		}
		if o.Stats[stat.Path].UsedMax < stat.Used {
			o.Stats[stat.Path].UsedMax = stat.Used
		}
		o.Stats[stat.Path].UsedPercent += stat.UsedPercent
		o.Stats[stat.Path].UsedPercentAvg = float64(o.Stats[stat.Path].UsedPercent) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].UsedPercentMin > stat.UsedPercent {
			o.Stats[stat.Path].UsedPercentMin = stat.UsedPercent
		}
		if o.Stats[stat.Path].UsedPercentMax < stat.UsedPercent {
			o.Stats[stat.Path].UsedPercentMax = stat.UsedPercent
		}
		o.Stats[stat.Path].InodesTotal = stat.InodesTotal
		o.Stats[stat.Path].InodesUsed += stat.InodesUsed
		o.Stats[stat.Path].InodesUsedAvg = float64(o.Stats[stat.Path].InodesUsed) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].InodesUsedMin > stat.InodesUsed {
			o.Stats[stat.Path].InodesUsedMin = stat.InodesUsed
		}
		if o.Stats[stat.Path].InodesUsedMax < stat.InodesUsed {
			o.Stats[stat.Path].InodesUsedMax = stat.InodesUsed
		}
		o.Stats[stat.Path].InodesFree += stat.InodesFree
		o.Stats[stat.Path].InodesFreeAvg = float64(o.Stats[stat.Path].InodesFree) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].InodesFreeMin > stat.InodesFree {
			o.Stats[stat.Path].InodesFreeMin = stat.InodesFree
		}
		if o.Stats[stat.Path].InodesFreeMin < stat.InodesFree {
			o.Stats[stat.Path].InodesFreeMin = stat.InodesFree
		}
		o.Stats[stat.Path].InodesUsedPercent += stat.InodesUsedPercent
		o.Stats[stat.Path].InodesUsedPercentAvg = float64(o.Stats[stat.Path].InodesUsedPercent) / float64(o.Stats[stat.Path].Count)
		if o.Stats[stat.Path].InodesUsedPercentMin > stat.InodesUsedPercent {
			o.Stats[stat.Path].InodesUsedPercentMin = stat.InodesUsedPercent
		}
		if o.Stats[stat.Path].InodesUsedPercentMax < stat.InodesUsedPercent {
			o.Stats[stat.Path].InodesUsedPercentMax = stat.InodesUsedPercent
		}
	}
}

func (o *DiskUsages) Empty() {

	o.Stats = make(map[string]*DiskUsage, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
	o.Count = 0
}

// GetModels returns the DiskUsage Aggregated member that contains a
// DiskUsage array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskUsages) GetModels() (*[]*DiskUsage, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no usages found")
	}
	rows := make([]*DiskUsage, 0)
	for _, stat := range o.Stats {
		rows = append(rows, stat)
	}
	return &rows, nil
}

/***************************************************************************************************
 * Disk input and output
 */

/**
 * Disk partition IO aggregated statistics
 */
type DiskIO struct {
	base.ModelCommon
	PartitionID         int64   `json:"partition_id" gorm:"column:partition_id;type:int;"`
	Device              string  `json:"device" gorm:"column:device;type:string;size:255;"`
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
	ReadBytes           uint64  `json:"-" gorm:"-"`
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
	base.ModelListCommon
	Stats map[string]*DiskIO
}

func (o *DiskModel) NewDiskIOs(host_id int64) *DiskIOs {

	o.DiskIOs = &DiskIOs{}
	o.DiskIOs.HostID = host_id
	return o.DiskIOs
}

func (o *DiskIOs) AddStat(i map[string]*disk.IOCountersStat, collect_duration int64) {

	if len(i) == 0 {
		return
	}
	if len(o.Stats) == 0 {
		o.CollectedStart = time.Now()
	}
	o.Count++
	o.CollectedEnd = time.Now()
	o.CollectedTotal += collect_duration
	o.CollectedAvg = float64(o.CollectedTotal) / float64(o.Count)
	if o.CollectedMin > collect_duration {
		o.CollectedMin = collect_duration
	}
	if o.CollectedMax > collect_duration {
		o.CollectedMax = collect_duration
	}

	for device, stat := range i {
		_, ok := o.Stats[device]
		if !ok {
			o.Stats[device] = &DiskIO{}
		}
		o.Stats[device].Count++
		o.Stats[device].HostID = o.HostID
		o.Stats[device].CollectedAt = o.CollectedEnd
		o.Stats[device].Device = device
		o.Stats[device].Name = stat.Name
		o.Stats[device].SerialNumber = stat.SerialNumber
		o.Stats[device].Label = stat.Label
		o.Stats[device].ReadCount += stat.ReadCount
		o.Stats[device].ReadCountAvg = float64(o.Stats[device].ReadCount) / float64(o.Stats[device].Count)
		if o.Stats[device].ReadCountMin > o.Stats[device].ReadCount {
			o.Stats[device].ReadCountMin = o.Stats[device].ReadCount
		}
		if o.Stats[device].ReadCountMax < o.Stats[device].ReadCount {
			o.Stats[device].ReadCountMax = o.Stats[device].ReadCount
		}
		o.Stats[device].MergedReadCount += stat.MergedReadCount
		o.Stats[device].MergedReadCountAvg = float64(o.Stats[device].MergedReadCount) / float64(o.Stats[device].Count)
		if o.Stats[device].MergedReadCountMin > o.Stats[device].MergedReadCount {
			o.Stats[device].MergedReadCountMin = o.Stats[device].MergedReadCount
		}
		if o.Stats[device].MergedReadCountMax < o.Stats[device].MergedReadCount {
			o.Stats[device].MergedReadCountMax = o.Stats[device].MergedReadCount
		}
		o.Stats[device].WriteCount += stat.WriteCount
		o.Stats[device].WriteCountAvg = float64(o.Stats[device].WriteCount) / float64(o.Stats[device].Count)
		if o.Stats[device].WriteCountMin > o.Stats[device].WriteCount {
			o.Stats[device].WriteCountMin = o.Stats[device].WriteCount
		}
		if o.Stats[device].WriteCountMax < o.Stats[device].WriteCount {
			o.Stats[device].WriteCountMax = o.Stats[device].WriteCount
		}
		o.Stats[device].MergedWriteCount += stat.MergedWriteCount
		o.Stats[device].MergedWriteCountAvg = float64(o.Stats[device].MergedWriteCount) / float64(o.Stats[device].Count)
		if o.Stats[device].MergedWriteCountMin > o.Stats[device].MergedWriteCount {
			o.Stats[device].MergedWriteCountMin = o.Stats[device].MergedWriteCount
		}
		if o.Stats[device].MergedWriteCountMax < o.Stats[device].MergedWriteCount {
			o.Stats[device].MergedWriteCountMax = o.Stats[device].MergedWriteCount
		}
		o.Stats[device].ReadBytes += stat.ReadBytes
		o.Stats[device].ReadBytesAvg = float64(o.Stats[device].ReadBytes) / float64(o.Stats[device].Count)
		if o.Stats[device].ReadBytesMin > o.Stats[device].ReadBytes {
			o.Stats[device].ReadBytesMin = o.Stats[device].ReadBytes
		}
		if o.Stats[device].ReadBytesMax < o.Stats[device].ReadBytes {
			o.Stats[device].ReadBytesMax = o.Stats[device].ReadBytes
		}
		o.Stats[device].WriteBytes += stat.WriteBytes
		o.Stats[device].WriteBytesAvg = float64(o.Stats[device].WriteBytes) / float64(o.Stats[device].Count)
		if o.Stats[device].WriteBytesMin > o.Stats[device].WriteBytes {
			o.Stats[device].WriteBytesMin = o.Stats[device].WriteBytes
		}
		if o.Stats[device].WriteBytesMax < o.Stats[device].WriteBytes {
			o.Stats[device].WriteBytesMax = o.Stats[device].WriteBytes
		}
		o.Stats[device].ReadTime += stat.ReadTime
		o.Stats[device].ReadTimeAvg = float64(o.Stats[device].ReadTime) / float64(o.Stats[device].Count)
		if o.Stats[device].ReadTimeMin > o.Stats[device].ReadTime {
			o.Stats[device].ReadTimeMin = o.Stats[device].ReadTime
		}
		if o.Stats[device].ReadTimeMax < o.Stats[device].ReadTime {
			o.Stats[device].ReadTimeMax = o.Stats[device].ReadTime
		}
		o.Stats[device].WriteTime += stat.WriteTime
		o.Stats[device].WriteTimeAvg = float64(o.Stats[device].WriteTime) / float64(o.Stats[device].Count)
		if o.Stats[device].WriteTimeMin > o.Stats[device].WriteTime {
			o.Stats[device].WriteTimeMin = o.Stats[device].WriteTime
		}
		if o.Stats[device].WriteTimeMax < o.Stats[device].WriteTime {
			o.Stats[device].WriteTimeMax = o.Stats[device].WriteTime
		}
		o.Stats[device].IopsInProgress += stat.IopsInProgress
		o.Stats[device].IopsInProgressAvg = float64(o.Stats[device].IopsInProgress) / float64(o.Stats[device].Count)
		if o.Stats[device].IopsInProgressMin > o.Stats[device].IopsInProgress {
			o.Stats[device].IopsInProgressMin = o.Stats[device].IopsInProgress
		}
		if o.Stats[device].IopsInProgressMax < o.Stats[device].IopsInProgress {
			o.Stats[device].IopsInProgressMax = o.Stats[device].IopsInProgress
		}
		o.Stats[device].IoTime += stat.IoTime
		o.Stats[device].IoTimeAvg = float64(o.Stats[device].IoTime) / float64(o.Stats[device].Count)
		if o.Stats[device].IoTimeMin > o.Stats[device].IoTime {
			o.Stats[device].IoTimeMin = o.Stats[device].IoTime
		}
		if o.Stats[device].IoTimeMax < o.Stats[device].IoTime {
			o.Stats[device].IoTimeMax = o.Stats[device].IoTime
		}
		o.Stats[device].WeightedIO += stat.WeightedIO
		o.Stats[device].WeightedIOAvg = float64(o.Stats[device].WeightedIO) / float64(o.Stats[device].Count)
		if o.Stats[device].WeightedIOMin > o.Stats[device].WeightedIO {
			o.Stats[device].WeightedIOMin = o.Stats[device].WeightedIO
		}
		if o.Stats[device].WeightedIOMax < o.Stats[device].WeightedIO {
			o.Stats[device].WeightedIOMax = o.Stats[device].WeightedIO
		}
	}
}

func (o *DiskIOs) Empty() {

	o.Stats = make(map[string]*DiskIO, 0)
	o.CollectedStart = time.Unix(0, 0)
	o.CollectedEnd = time.Unix(0, 0)
	o.CollectedAvg = 0
	o.CollectedMin = 0
	o.CollectedMax = 0
	o.Count = 0
}

// GetModels returns the DiskIOs Aggregated member that contains a
// DiskIOs array.
// If the Aggregated member is an empty array an error is
// returned
func (o *DiskIOs) GetModels() (*[]*DiskIO, error) {
	if len(o.Stats) == 0 {
		return nil, fmt.Errorf("no disk io data found")
	}
	rows := make([]*DiskIO, 0)
	for _, stat := range o.Stats {
		rows = append(rows, stat)
	}
	return &rows, nil
}
