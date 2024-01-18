package repository

import (
	"context"

	app "github.com/maxik12233/task-junior"
	"github.com/maxik12233/task-junior/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IRepository interface {
	CreatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error
	DeletePersonById(ctx context.Context, id int) error
	UpdatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error
}

//

type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRepository(db *gorm.DB, logger *zap.Logger) IRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error {
	createPerson := Person{}
	createChar := Characteristic{}
	createPerson.FromDomain(person)
	createChar.FromDomain(char)
	createPerson.Characteristic = createChar

	result := r.db.Create(&createPerson)
	if result.Error != nil {
		return app.ErrInternal
	}

	return nil
}

func (r *Repository) DeletePersonById(ctx context.Context, id int) error {

	var person Person
	result := r.db.Where("id = ?", id).Find(&person)
	if result.RowsAffected == 0 {
		return app.ErrNotFound
	}
	if result.Error != nil {
		return app.ErrInternal
	}

	result = r.db.Unscoped().Delete(&Person{}, id)
	if result.Error != nil {
		return app.ErrInternal
	}

	result = r.db.Unscoped().Delete(&Characteristic{}, person.CharacteristicID)
	if result.Error != nil {
		return app.ErrInternal
	}

	return nil
}

func (r *Repository) UpdatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		updatePerson := Person{}
		updateChar := Characteristic{}
		updatePerson.FromDomain(person)
		updateChar.FromDomain(char)

		var p Person
		result := r.db.Where("id = ?", updatePerson.ID).Find(&p)
		if result.RowsAffected == 0 {
			return app.ErrNotFound
		}
		if result.Error != nil {
			return app.ErrInternal
		}

		updateChar.ID = uint(p.CharacteristicID)

		result = r.db.Save(&updateChar)
		if result.Error != nil {
			return app.ErrInternal
		}

		result = r.db.Omit("Characteristic").Omit("CharacteristicID").Save(&updatePerson)
		if result.Error != nil {
			return app.ErrInternal
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
