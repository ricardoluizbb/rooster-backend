package account

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	User struct {
		ID           string    `json:"id,omitempty" gorm:"primaryKey"`
		Name         string    `json:"name,omitempty"`
		Email        string    `json:"email,omitempty"`
		Password     string    `json:"password,omitempty"`
		Token        string    `json:"token,omitempty"`
		RefreshToken string    `json:"refreshToken,omitempty"`
		CreatedAt    time.Time `json:"createdAt,omitempty"`
		UpdatedAt    time.Time `json:"updatedAt,omitempty"`
		DeletedAt    time.Time `json:"deletedAt,omitempty"`
	}
)

func NewUser(name, email, id, password string) (*User, error) {
	if id == "" {
		id = uuid.NewV4().String()
	}

	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}, nil
}
