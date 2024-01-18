package service

import "github.com/maxik12233/task-junior/internal/domain"

type CombinedInfo struct {
	Name        string
	Age         int
	Gender      string
	Nationality string
}

func (c *CombinedInfo) ToDomainCharactaristic() domain.Characteristic {
	return domain.Characteristic{
		Age:         c.Age,
		Gender:      c.Gender,
		Nationality: c.Nationality,
	}
}
