package repository

import (
	_ "database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishPostgresConnection(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DoAutoMigration(db *gorm.DB) error {
	err := db.AutoMigrate(Person{}, Characteristic{})
	if err != nil {
		return err
	}

	return nil
}

func DoCommonMigration(dbUrl string) error {
	m, err := migrate.New("file://internal/repository/migrations",
		dbUrl)
	if err != nil {
		return err
	}

	m.Up()
	return nil
}

func Rollback(dbUrl string) error {
	m, err := migrate.New("file://internal/repository/migrations",
		dbUrl)
	if err != nil {
		return err
	}

	m.Down()
	return nil
}
