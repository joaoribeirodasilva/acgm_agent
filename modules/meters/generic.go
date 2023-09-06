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

type MeterMinMaxAvgInt64 struct {
	Count int
	Min   int64
	Max   int64
	Avg   float64
	Total int64
}

func (m *MeterMinMaxAvgInt64) Aggregate(value int64) {
	m.Total += value
	if m.Min > value {
		m.Min = value
	}
	if m.Max < value {
		m.Max = value
	}
	m.Count++
	m.Avg = float64(m.Total) / float64(m.Count)
}

type MeterMinMaxAvgUInt64 struct {
	Count int
	Min   uint64
	Max   uint64
	Avg   float64
	Total uint64
}

func (m *MeterMinMaxAvgUInt64) Aggregate(value uint64) {
	m.Total += value
	if m.Min > value {
		m.Min = value
	}
	if m.Max < value {
		m.Max = value
	}
	m.Count++
	m.Avg = float64(m.Total) / float64(m.Count)
}

type MeterMinMaxAvgInt struct {
	Count int
	Min   int
	Max   int
	Avg   float64
	Total int64
}

func (m *MeterMinMaxAvgInt) Aggregate(value int) {
	m.Total += int64(value)
	if m.Min > value {
		m.Min = value
	}
	if m.Max < value {
		m.Max = value
	}
	m.Count++
	m.Avg = float64(m.Total) / float64(m.Count)
}

type MeterMinMaxAvgFloat64 struct {
	Count int
	Min   float64
	Max   float64
	Avg   float64
	Total float64
}

func (m *MeterMinMaxAvgFloat64) Aggregate(value float64) {
	m.Total += value
	if m.Min > value {
		m.Min = value
	}
	if m.Max < value {
		m.Max = value
	}
	m.Count++
	m.Avg = m.Total / float64(m.Count)
}

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
