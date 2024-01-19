package service

import (
	"context"
	"errors"

	app "github.com/maxik12233/task-junior"
	"github.com/maxik12233/task-junior/internal/domain"
	"github.com/maxik12233/task-junior/internal/repository"
	"github.com/maxik12233/task-junior/pkg/api/paginate"
	"github.com/maxik12233/task-junior/pkg/api/sort"
	"github.com/maxik12233/task-junior/pkg/name_info_sdk"
	"go.uber.org/zap"
)

type IService interface {
	CreatePersonInfo(ctx context.Context, person domain.Person) (CombinedInfo, error)
	DeletePersonInfo(ctx context.Context, id int) error
	UpdatePersonInfo(ctx context.Context, person domain.Person, char domain.Characteristic) error
	GetAllPersonInfo(ctx context.Context, sortOption *sort.Options, paginateOption *paginate.Options) ([]*domain.Person, error)
	GetPersonInfo(ctx context.Context, id uint) (*domain.Person, error)
	GetPersonCount(ctx context.Context) (int, error)
}

//

type Service struct {
	repo          repository.IRepository
	logger        *zap.Logger
	byNameService name_info_sdk.INameInfo
}

func NewService(repo repository.IRepository, logger *zap.Logger, byNameService name_info_sdk.INameInfo) IService {
	return &Service{
		repo:          repo,
		logger:        logger,
		byNameService: byNameService,
	}
}

func (s *Service) GetPersonCount(ctx context.Context) (int, error) {
	count, err := s.repo.GetPersonCount(ctx)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *Service) GetPersonInfo(ctx context.Context, id uint) (*domain.Person, error) {
	person, err := s.repo.GetPersonById(ctx, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (s *Service) GetAllPersonInfo(ctx context.Context, sortOption *sort.Options, paginateOption *paginate.Options) ([]*domain.Person, error) {

	var (
		OptionSort     repository.SortOptions
		OptionPaginate repository.PaginateOptions
	)

	if sortOption != nil {
		OptionSort = repository.NewSortOptions(sortOption.Field, sortOption.Order)
	}

	if paginateOption != nil {
		OptionPaginate = repository.NewPaginateOptions(paginateOption.Page, paginateOption.PerPage)
	}

	persons, err := s.repo.GetPersonAll(ctx, OptionSort, OptionPaginate)
	if err != nil {
		return nil, err
	}

	return persons, nil
}

func (s *Service) CreatePersonInfo(ctx context.Context, person domain.Person) (CombinedInfo, error) {

	info, err := s.fetchAllNameInfo(person.Name)
	if err != nil {
		s.logger.Error("Error fetching data from foreign api", zap.Error(err))
		return CombinedInfo{}, app.ErrInternal
	}

	if err := s.repo.CreatePerson(ctx, person, info.ToDomainCharactaristic()); err != nil {
		return CombinedInfo{}, err
	}

	return info, nil
}

func (s *Service) DeletePersonInfo(ctx context.Context, id int) error {

	if err := s.repo.DeletePerson(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePersonInfo(ctx context.Context, person domain.Person, char domain.Characteristic) error {

	if err := s.repo.UpdatePerson(ctx, person, char); err != nil {
		return err
	}

	return nil
}

func (s *Service) fetchAllNameInfo(name string) (CombinedInfo, error) {
	var (
		ageChan    = make(chan *name_info_sdk.LikelyAge)
		genderChan = make(chan *name_info_sdk.LikelyGender)
		natChan    = make(chan *name_info_sdk.LikelyNationality)
	)

	go func() {
		resp, err := s.byNameService.GetAgeInfoByName(name)
		if err != nil {
			s.logger.Error("Error while fetching age info", zap.Error(err))
			ageChan <- nil
			return
		}

		ageChan <- resp
	}()
	go func() {
		resp, err := s.byNameService.GetGenderInfoByName(name)
		if err != nil {
			s.logger.Error("Error while fetching gender info", zap.Error(err))
			genderChan <- nil
			return
		}

		genderChan <- resp
	}()
	go func() {
		resp, err := s.byNameService.GetLikelyNationalityInfoByName(name)
		if err != nil {
			s.logger.Error("Error while fetching nationality info", zap.Error(err))
			natChan <- nil
			return
		}

		natChan <- resp
	}()

	age := <-ageChan
	gender := <-genderChan
	nationality := <-natChan

	if age == nil || gender == nil || nationality == nil {
		return CombinedInfo{}, errors.New("Some request ended badly")
	}

	return CombinedInfo{
		Name:        name,
		Age:         age.Age,
		Gender:      gender.Gender,
		Nationality: nationality.Nationality,
	}, nil
}
