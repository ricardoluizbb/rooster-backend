package account

import "gorm.io/gorm"

type (
	Persistence struct {
		DB *gorm.DB
	}
	AccountPersistence interface {
		Create(u *User) error
		Get(email string) (*User, error)
		GetUserByRefreshToken(token string) (*User, error)
		GetUserByToken(token string) (*User, error)
		Save(u *User) error
	}
)

func NewPersistence(db *gorm.DB) AccountPersistence {
	return &Persistence{
		DB: db,
	}
}

func (p *Persistence) Create(u *User) error {
	return p.DB.Model(&User{}).Create(u).Error
}

func (p *Persistence) Get(email string) (*User, error) {
	var user *User
	tx := p.DB.Model(&User{}).Where("email = ?", email).First(&user)

	return user, tx.Error
}

func (p *Persistence) GetUserByRefreshToken(token string) (*User, error) {
	var user *User
	tx := p.DB.Model(&User{}).Where("refresh_token = ?", token).First(&user)

	return user, tx.Error
}

func (p *Persistence) GetUserByToken(token string) (*User, error) {
	var user *User
	tx := p.DB.Model(&User{}).Where("token = ?", token).First(&user)

	return user, tx.Error
}

func (p *Persistence) Save(u *User) error {
	return p.DB.Save(u).Error
}
