package database

import (
	"time-tracker-backend/account"
	"time-tracker-backend/config"
	"time-tracker-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDataBase() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migra as tabelas do banco de dados
	err = db.AutoMigrate(&models.Task{}, &models.User{}, &models.RegisteredTime{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations() error {
	db, err := ConnectDataBase()
	if err != nil {
		return err
	}

	// Migra as tabelas do banco de dados
	err = db.AutoMigrate(&models.Task{}, &models.User{}, &models.RegisteredTime{}, &account.User{})
	if err != nil {
		return err
	}

	return nil
}
