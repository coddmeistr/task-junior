package models

import "gorm.io/gorm"

type NameStatistic struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}
