package repository

import "github.com/maxik12233/task-junior/internal/domain"

type Person struct {
	ID               uint   `gorm:"primary key"`
	Name             string `gorm:"not null"`
	Surname          string `gorm:"not null"`
	Patronymic       string
	CharacteristicID int
	Characteristic   Characteristic
}

type Characteristic struct {
	ID          uint   `gorm:"primary key"`
	Age         int    `gorm:"not null"`
	Gender      string `gorm:"not null"`
	Nationality string `gorm:"not null"`
}

func (p *Person) FromDomain(person domain.Person) {
	if person.ID != 0 {
		p.ID = person.ID
	}
	p.Name = person.Name
	p.Surname = person.Surname
	p.Patronymic = person.Patronymic
}

func (p *Characteristic) FromDomain(person domain.Characteristic) {
	p.Age = person.Age
	p.Gender = person.Gender
	p.Nationality = person.Nationality
}
