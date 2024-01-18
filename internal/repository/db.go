package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishDatabaseConnection(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(Person{}, Characteristic{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
