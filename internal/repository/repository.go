package repository

import (
	"context"

	app "github.com/maxik12233/task-junior"
	"github.com/maxik12233/task-junior/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository interface {
	GetPersonCount(ctx context.Context) (int64, error)
	GetPersonAll(ctx context.Context, sortOptions SortOptions, paginateOptions PaginateOptions) ([]*domain.Person, error)
	GetPersonById(ctx context.Context, id uint) (*domain.Person, error)
	CreatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error
	DeletePerson(ctx context.Context, id int) error
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

func (r *Repository) GetPersonAll(ctx context.Context, sortOptions SortOptions, paginateOptions PaginateOptions) ([]*domain.Person, error) {
	var persons []*Person
	result := r.db.Preload(clause.Associations).Offset(int(paginateOptions.GetPage()) * int(paginateOptions.GetPerPage())).
		Limit(int(paginateOptions.GetPerPage())).
		Order(sortOptions.GetOrderBy()).
		Find(&persons)
	if result.Error != nil {
		r.logger.Error("Error getting all person infos", zap.Error(result.Error))
		return nil, app.ErrInternal
	}

	var domainPersons = make([]*domain.Person, len(persons))
	for i, v := range persons {
		domainPerson := v.ToDomain()
		domainPersons[i] = &domainPerson
	}

	return domainPersons, nil
}

func (r *Repository) GetPersonCount(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.Model(&Person{}).Count(&count)
	if result.Error != nil {
		r.logger.Error("Error counting all person infos", zap.Error(result.Error))
		return 0, app.ErrInternal
	}

	return count, nil
}

func (r *Repository) GetPersonById(ctx context.Context, id uint) (*domain.Person, error) {
	var p *Person
	result := r.db.Preload(clause.Associations).Find(&p, id)
	if result.RowsAffected == 0 {
		r.logger.Error("Error not found while getting person by id")
		return nil, app.ErrNotFound
	}
	if result.Error != nil {
		r.logger.Error("Error getting person by id", zap.Error(result.Error))
		return nil, app.ErrInternal
	}

	resultValue := p.ToDomain()
	return &resultValue, nil
}

func (r *Repository) CreatePerson(ctx context.Context, person domain.Person, char domain.Characteristic) error {
	createPerson := Person{}
	createChar := Characteristic{}
	createPerson.FromDomain(person)
	createChar.FromDomain(char)
	createPerson.Characteristic = createChar

	result := r.db.Create(&createPerson)
	if result.Error != nil {
		r.logger.Error("Error creating new person info", zap.Error(result.Error))
		return app.ErrInternal
	}

	return nil
}

func (r *Repository) DeletePerson(ctx context.Context, id int) error {

	var person Person
	result := r.db.Where("id = ?", id).Find(&person)
	if result.RowsAffected == 0 {
		r.logger.Error("Error not found while deleting person info")
		return app.ErrNotFound
	}
	if result.Error != nil {
		r.logger.Error("Error while getting person by id", zap.Error(result.Error))
		return app.ErrInternal
	}

	result = r.db.Unscoped().Delete(&Person{}, id)
	if result.Error != nil {
		r.logger.Error("Error while deleting person info", zap.Error(result.Error))
		return app.ErrInternal
	}

	result = r.db.Unscoped().Delete(&Characteristic{}, person.CharacteristicID)
	if result.Error != nil {
		r.logger.Error("Error while deleting person charactaristic", zap.Error(result.Error))
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
			r.logger.Error("Error not found while getting person by id")
			return app.ErrNotFound
		}
		if result.Error != nil {
			r.logger.Error("Error getting person by id", zap.Error(result.Error))
			return app.ErrInternal
		}

		updateChar.ID = uint(p.CharacteristicID)

		result = r.db.Save(&updateChar)
		if result.Error != nil {
			r.logger.Error("Error updating person's charactatistic", zap.Error(result.Error))
			return app.ErrInternal
		}

		result = r.db.Omit("Characteristic").Omit("CharacteristicID").Save(&updatePerson)
		if result.Error != nil {
			r.logger.Error("Error updating person", zap.Error(result.Error))
			return app.ErrInternal
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
