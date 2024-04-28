package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID             string           `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt      time.Time        `json:"created_at,omitempty"`
	UpdatedAt      time.Time        `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt   `json:"deleted_at,omitempty"`
	Title          string           `json:"title,omitempty"`
	Tag            string           `json:"tag,omitempty"`
	UserID         string           `json:"user_id,omitempty"`
	Done           bool             `json:"done,omitempty"`
	RegisteredTime []RegisteredTime `gorm:"foreignKey:TaskID" json:"registered_times,omitempty"`
}
