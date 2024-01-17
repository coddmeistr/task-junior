package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IRepository interface {
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
