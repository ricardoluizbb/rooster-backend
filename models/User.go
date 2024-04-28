package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
