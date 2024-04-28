package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID         string         `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
	TaskID     string         `json:"task_id,omitempty"`
	Type       string         `json:"type,omitempty"`
	TotalHours time.Duration  `json:"total_hours,omitempty"`
}
