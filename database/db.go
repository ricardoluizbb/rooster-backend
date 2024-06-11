package database

import (
	"time-tracker-backend/config"
	"time-tracker-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	firebaseAuth "firebase.google.com/go/v4/auth"
)

var (
	DB           *gorm.DB
	FirebaseAuth *firebaseAuth.Client
)

func ConnectDataBase() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{})
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
	err = db.AutoMigrate(&models.Task{}, &models.Report{}, &models.RegisteredTime{}, &models.User{})
	if err != nil {
		return err
	}

	return nil
}
