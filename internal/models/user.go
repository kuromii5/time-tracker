package models

import "time"

type User struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Passport  Passport  `json:"passport"`
	People    People    `json:"people"`
}

type Passport struct {
	Serie  string `json:"serie"`
	Number string `json:"number"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

// FilterBy represents filtering criteria for users.
// swagger:model
type FilterBy struct {
	Name           string    `json:"name"`
	Surname        string    `json:"surname"`
	Patronymic     string    `json:"patronymic"`
	Address        string    `json:"address"`
	CreatedAfter   time.Time `json:"created_after"`
	CreatedBefore  time.Time `json:"created_before"`
	PassportSerie  string    `json:"passport_serie"`
	PassportNumber string    `json:"passport_number"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
