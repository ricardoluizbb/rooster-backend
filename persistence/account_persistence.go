package persistence

import "gorm.io/gorm"

type AccountPersistence struct {
	db *gorm.DB
}

// Criar e atualizar usuário
// Essa função será usada para atualizar o token do usuário
func (p *AccountPersistence) SaveUser() {

}

func (p *AccountPersistence) GetUser() {

}
