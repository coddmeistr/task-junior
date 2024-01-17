package service

import (
	"github.com/maxik12233/task-junior/internal/repository"
	"go.uber.org/zap"
)

type IService interface {
}

//

type Service struct {
	repo   repository.IRepository
	logger *zap.Logger
}

func NewService(repo repository.IRepository, logger *zap.Logger) IService {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}
