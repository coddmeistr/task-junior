package transport

import "github.com/maxik12233/task-junior/internal/domain"

type AddPersonInfoRequest struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
}

func (person *AddPersonInfoRequest) ToDomain() domain.Person {
	return domain.Person{
		Name:       person.Name,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	}
}

type AddPersonInfoResponse struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type DeletePersonInfoRequest struct {
	Id uint `json:"id" validate:"required,gte=1"`
}

func (person *DeletePersonInfoRequest) ToDomain() domain.Person {
	return domain.Person{
		ID: person.Id,
	}
}

type UpdatePersonInfoRequest struct {
	Id          uint   `json:"id" validate:"required,gte=1"`
	Name        string `json:"name" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender" validate:"required"`
	Age         int    `json:"age" validate:"required,gte=0,lte=200"`
	Nationality string `json:"nationality" validate:"required"`
}

func (d *UpdatePersonInfoRequest) ToDomain() (domain.Person, domain.Characteristic) {
	return domain.Person{
			ID:         d.Id,
			Name:       d.Name,
			Surname:    d.Surname,
			Patronymic: d.Patronymic,
		}, domain.Characteristic{
			Age:         d.Age,
			Gender:      d.Gender,
			Nationality: d.Nationality,
		}
}

type PersonResponse struct {
	Id          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"genderv"`
	Age         int    `json:"age,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type GetPersonInfoResponse struct {
	PersonResponse
	Persons []PersonResponse `json:"persons"`
}
