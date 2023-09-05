package meters

import (
	"time"

	"biqx.com.br/acgm_agent/modules/config"
	"biqx.com.br/acgm_agent/modules/database"
)

type MeterStatus int

const (
	STATUS_STOPPED MeterStatus = iota
	STATUS_STARTING
	STATUS_STARTED
	STATUS_COLLECTING
	STATUS_AGGREGATING
	STATUS_WAITING
	STATUS_STOPPING
)

type MeterControl struct {
	name   string
	host   string
	conf   *config.Config
	status MeterStatus
	db     *database.Db
}

type MetricTimes struct {
	DateTime time.Time `json:"date_time" yaml:"date_time"`
	Duration int64     `json:"duration" yaml:"duration"`
}

type InterfaceMeter interface {
	Start()
	Stop()
	GetName() string
	GetStatus()
	Collect()
	Aggregate()
	FireEvent(status MeterStatus)
	EventHandler()
}
