package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type RegisteredTime struct {
	ID        string         `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	TaskID    string         `json:"task_id,omitempty"`
	StartTime time.Time      `json:"start_time,omitempty"`
	EndTime   *time.Time     `json:"end_time,omitempty"`
	TotalTime float64        `gorm:"-" json:"total_time,omitempty"`
	Paused    bool           `json:"paused,omitempty"`
	UserID    string         `json:"user_id,omitempty"`
}

type TaskTotalTime struct {
	TaskID          string            `json:"task_id,omitempty"`
	RegisteredTimes []*RegisteredTime `json:"registeredTimes,omitempty"`
	TotalTime       float64           `json:"total_time,omitempty"`
}

func (r *RegisteredTime) SetTotalTime() {
	if r.EndTime == nil {
		r.TotalTime = 0
		return
	}
	r.TotalTime = math.Abs(float64(r.EndTime.UnixMilli() - r.StartTime.UnixMilli()))
}
