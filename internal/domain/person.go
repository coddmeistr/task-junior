package domain

type Person struct {
	ID               uint
	Name             string
	Surname          string
	Patronymic       string
	CharacteristicID int
}

type Characteristic struct {
	ID          uint
	Age         int
	Gender      string
	Nationality string
}
