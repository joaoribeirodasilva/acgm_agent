package models

import "time"

type HostInfo struct {
	ID                   int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	Hostname             string    `json:"hostname" gorm:"column:hostname;type:string;size:255"`
	Uptime               time.Time `json:"uptime" gorm:"column:uptime;type:timestamp;"`
	BootTime             time.Time `json:"boot_time" gorm:"column:boot_time;type:timestamp;"`
	OS                   string    `json:"os" gorm:"column:os;type:string;size:255"`                             // ex: freebsd, linux
	Platform             string    `json:"platform" gorm:"column:platform;type:string;size:255"`                 // ex: ubuntu, linuxmint
	PlatformFamily       string    `json:"platform_family" gorm:"column:platform_family;type:string;size:255"`   // ex: debian, rhel
	PlatformVersion      string    `json:"platform_version" gorm:"column:platform_version;type:string;size:255"` // version of the complete OS
	KernelVersion        string    `json:"kernel_version" gorm:"column:kernel_version;type:string;size:255"`     // version of the OS kernel (if available)
	KernelArch           string    `json:"kernel_arch" gorm:"column:kernel_arch;type:string;size:255"`           // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string    `json:"virtualization_system" gorm:"column:virtualization_system;type:string;size:255"`
	VirtualizationRole   string    `json:"virtualization_role" gorm:"column:virtualization_role;type:string;size:255"` // guest or host
	HostID               string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`                         // ex: uuid
	CollectedAt          time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis      uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (HostInfo) TableName() string {
	return "host_info"
}

type HostTemperatures struct {
	ID              int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID          string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`
	SensorKey       string    `json:"sensor_key" gorm:"column:sensor_key;type:string;size:255"`
	TemperatureAvg  float64   `json:"temperature_avg" gorm:"column:temperature_avg;type:float"`
	TemperatureMax  float64   `json:"temperature_max" gorm:"column:temperature_max;type:float"`
	TemperatureMin  float64   `json:"temperature_min" gorm:"column:temperature_min;type:float"`
	CollectedAt     time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (HostTemperatures) TableName() string {
	return "host_temperatures"
}

type HostUsers struct {
	ID              int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID          string    `json:"host_id" gorm:"column:host_id;type:string;size:255"`
	User            string    `json:"user" gorm:"column:user;type:string;size:255"`
	Terminal        string    `json:"terminal" gorm:"column:terminal;type:string;size:255"`
	Host            string    `json:"host" gorm:"column:host;type:string;size:255"`
	Started         time.Time `json:"started" gorm:"column:started;type:timestamp"`
	CollectedAt     time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	CollectedMillis uint64    `json:"collected_milliseconds" gorm:"column:collected_milliseconds;type:int;size:20;"`
}

func (HostUsers) TableName() string {
	return "host_users"
}
