package models

import "time"

type CPU struct {
	ID              int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID          string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`
	CPU             int32     `json:"cpu" gorm:"column:cpu;type:int"`
	PercentAvg      float64   `json:"percent_avg" gorm:"column:percent_avg;type:float"`
	PercentMin      float64   `json:"percent_min" gorm:"column:percent_min;type:float"`
	PercentMax      float64   `json:"percent_max" gorm:"column:percent_max;type:float"`
	Load1Avg        float64   `json:"load_1_avg" gorm:"column:load_1_avg;type:float"`
	Load1Min        float64   `json:"load_1_min" gorm:"column:load_1_min;type:float"`
	Load1Max        float64   `json:"load_1_max" gorm:"column:load_1_max;type:float"`
	Load5Avg        float64   `json:"load_5_avg" gorm:"column:load_5_avg;type:float"`
	Load5Min        float64   `json:"load_5_min" gorm:"column:load_5_min;type:float"`
	Load5Max        float64   `json:"load_5_max" gorm:"column:load_5_max;type:float"`
	Load15Avg       float64   `json:"load_15_avg" gorm:"column:load_15_avg;type:float"`
	Load15Min       float64   `json:"load_15_min" gorm:"column:load_15_min;type:float"`
	Load15Max       float64   `json:"load_15_max" gorm:"column:load_15_max;type:float"`
	ProcsTotalAvg   float64   `json:"procs_total_avg" gorm:"column:procs_total_avg;type:float"`
	ProcsTotalMin   int       `json:"procs_total_min" gorm:"column:procs_total_min;type:int"`
	ProcsTotalMax   int       `json:"procs_total_max" gorm:"column:procs_total_max;type:int"`
	ProcsCreatedAvg float64   `json:"procs_created_avg" gorm:"column:procs_created_avg;type:float"`
	ProcsCreatedMin int       `json:"procs_created_min" gorm:"column:procs_created_min;type:int"`
	ProcsCreatedMax int       `json:"procs_created_max" gorm:"column:procs_created_max;type:int"`
	ProcsRunningAvg float64   `json:"procs_running_avg" gorm:"column:procs_running_avg;type:float"`
	ProcsRunningMin int       `json:"procs_running_min" gorm:"column:procs_running_min;type:int"`
	ProcsRunningMax int       `json:"procs_running_max" gorm:"column:procs_running_max;type:int"`
	ProcsBlockedAvg float64   `json:"procs_blocked_avg" gorm:"column:procs_blocked_avg;type:float"`
	ProcsBlockedMin int       `json:"procs_blocked_min" gorm:"column:procs_blocked_min;type:int"`
	ProcsBlockedMax int       `json:"procs_blocked_max" gorm:"column:procs_blocked_max;type:int"`
	CtxtAvg         float64   `json:"ctxt_avg" gorm:"column:ctxt_avg;type:float"`
	CtxtMin         int       `json:"ctxt_min" gorm:"column:ctxt_min;type:int"`
	CtxtMax         int       `json:"ctxt_max" gorm:"column:ctxt_max;type:int"`
	PhysicalCores   int       `json:"physical_cores" gorm:"column:physical_cores;type:int"`
	LogicalCores    int       `json:"logical_cores" gorm:"column:logical_cores;type:int"`
	CollectedAt     time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (CPU) TableName() string {
	return "cpus"
}

type CPUTimes struct {
	ID                int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID            string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`
	CPU               string    `json:"cpu" gorm:"column:cpu;type:string;size:255"`
	CoreIndex         int32     `json:"core_index" gorm:"column:core_index;type:int"`
	PercentAvg        float64   `json:"percent_avg" gorm:"column:percent_avg;type:float"`
	PercentMin        float64   `json:"percent_min" gorm:"column:percent_min;type:float"`
	PercentMax        float64   `json:"percent_max" gorm:"column:percent_max;type:float"`
	TimesTotalAvg     float64   `json:"times_total_avg" gorm:"column:times_total_avg;type:float"`
	TimesTotalMin     float64   `json:"times_total_min" gorm:"column:times_total_min;type:float"`
	TimesTotalMax     float64   `json:"times_total_max" gorm:"column:times_total_max;type:float"`
	TimesUserAvg      float64   `json:"times_user_avg" gorm:"column:times_user_avg;type:float"`
	TimesUserMin      float64   `json:"times_user_min" gorm:"column:times_user_min;type:float"`
	TimesUserMax      float64   `json:"times_user_max" gorm:"column:times_user_max;type:float"`
	TimesSystemAvg    float64   `json:"times_system_avg" gorm:"column:times_system_avg;type:float"`
	TimesSystemMin    float64   `json:"times_system_min" gorm:"column:times_system_min;type:float"`
	TimesSystemMax    float64   `json:"times_system_max" gorm:"column:times_system_max;type:float"`
	TimesIdleAvg      float64   `json:"times_idle_avg" gorm:"column:times_idle_avg;type:float"`
	TimesIdleMin      float64   `json:"times_idle_min" gorm:"column:times_idle_min;type:float"`
	TimesIdleMax      float64   `json:"times_idle_max" gorm:"column:times_idle_max;type:float"`
	TimesNiceAvg      float64   `json:"times_nice_avg" gorm:"column:times_nice_avg;type:float"`
	TimesNiceMin      float64   `json:"times_nice_min" gorm:"column:times_nice_min;type:float"`
	TimesNiceMax      float64   `json:"times_nice_max" gorm:"column:times_nice_max;type:float"`
	TimesIOWaitAvg    float64   `json:"times_io_wait_avg" gorm:"column:times_io_wait_avg;type:float"`
	TimesIOWaitMin    float64   `json:"times_io_wait_min" gorm:"column:times_io_wait_min;type:float"`
	TimesIOWaitMax    float64   `json:"times_io_wait_max" gorm:"column:times_io_wait_max;type:float"`
	TimesIrqAvg       float64   `json:"times_irq_avg" gorm:"column:times_irq_avg;type:float"`
	TimesIrqMin       float64   `json:"times_irq_min" gorm:"column:times_irq_min;type:float"`
	TimesIrqMax       float64   `json:"times_irq_max" gorm:"column:times_irq_max;type:float"`
	TimesSoftirqAvg   float64   `json:"times_softirq_avg" gorm:"column:times_softirq_avg;type:float"`
	TimesSoftirqMin   float64   `json:"times_softirq_min" gorm:"column:times_softirq_min;type:float"`
	TimesSoftirqMax   float64   `json:"times_softirq_max" gorm:"column:times_softirq_max;type:float"`
	TimesStealAvg     float64   `json:"times_steal_avg" gorm:"column:times_steal_avg;type:float"`
	TimesStealMin     float64   `json:"times_steal_min" gorm:"column:times_steal_min;type:float"`
	TimesStealMax     float64   `json:"times_steal_max" gorm:"column:times_steal_max;type:float"`
	TimesGuestAvg     float64   `json:"times_guest_avg" gorm:"column:times_guest_avg;type:float"`
	TimesGuestMin     float64   `json:"times_guest_min" gorm:"column:times_guest_min;type:float"`
	TimesGuestMax     float64   `json:"times_guest_max" gorm:"column:times_guest_max;type:float"`
	TimesGuestNiceAvg float64   `json:"times_guest_nice_avg" gorm:"column:times_guest_nice_avg;type:float"`
	TimesGuestNiceMin float64   `json:"times_guest_nice_min" gorm:"column:times_guest_nice_min;type:float"`
	TimesGuestNiceMax float64   `json:"times_guest_nice_max" gorm:"column:times_guest_nice_max;type:float"`
	CollectedAt       time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis   uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (CPUTimes) TableName() string {
	return "cpu_times"
}

type CPUCoreInfo struct {
	ID              int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID          string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`
	CPU             int32     `json:"cpu" gorm:"column:cpu;type:int"`
	CoreIndex       *int      `json:"core_index" gorm:"column:core_index;type:int"`
	VendorID        string    `json:"vendor_id" gorm:"column:vendor_id;type:string;size:255"`
	Family          string    `json:"family" gorm:"column:family;type:string;size:255"`
	Model           string    `json:"model" gorm:"column:model;type:string;size:255"`
	Stepping        int32     `json:"stepping" gorm:"column:stepping;type:int"`
	PhysicalID      string    `json:"physical_id" gorm:"column:physical_id;type:string;size:255"`
	CoreID          string    `json:"core_id" gorm:"column:core_id;type:string;size:255"`
	Cores           int32     `json:"cores" gorm:"column:cores;type:int"`
	ModelName       string    `json:"model_name" gorm:"column:model_name;type:string;size:255"`
	Mhz             float64   `json:"mhz" gorm:"column:mhz;type:float"`
	CacheSize       int32     `json:"cache_size" gorm:"column:cache_size;type:int"`
	Flags           string    `json:"flags" gorm:"column:flags;type:string;"`
	Microcode       string    `json:"microcode" gorm:"column:microcode;type:string;size:255"`
	CollectedAt     time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (CPUCoreInfo) TableName() string {
	return "cpu_core_infos"
}
