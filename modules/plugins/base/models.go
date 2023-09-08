package base

import "time"

type ModelListCommon struct {
	HostID         int64
	CollectedStart time.Time `json:"collected_start" gorm:"column:collected_start;type:timestamp;"`
	CollectedEnd   time.Time `json:"collected_end" gorm:"column:collected_end;type:timestamp;"`
	CollectedTotal int64     `json:"-" gorm:"-"`
	CollectedAvg   float64   `json:"collected_avg" gorm:"column:collected_avg;type:float;"`
	CollectedMin   int64     `json:"collected_min" gorm:"column:collected_min;type:int;"`
	CollectedMax   int64     `json:"collected_max" gorm:"column:collected_max;type:int;"`
	Count          int64     `json:"-" gorm:"-"`
}

type ModelCommon struct {
	ID          int64     `json:"id" gorm:"column:id;type:int;autoIncrement;primaryKey;"`
	HostID      int64     `json:"host_id" gorm:"column:host_id;type:int"`
	CollectedAt time.Time `json:"collected_at" gorm:"column:collected_at;type:timestamp;"`
	Count       uint      `json:"-" gorm:"-"`
}

func (m *ModelCommon) SetID(id int64) {
	m.ID = id
}
