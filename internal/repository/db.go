package repository

import (
	"github.com/maxik12233/task-junior/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishDatabaseConnection(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.NameStatistic{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
