package model

import "github.com/google/uuid"

// User is the db schema for the user table
type User struct {
	Surname    string `json:"surname" db:"surname"`
	Name       string `json:"name" db:"name"`
	Address    string `json:"address" db:"address"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Base
}

type AddUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}

type UpdateUserRequest struct {
	ID         uuid.UUID `json:"id"`
	Surname    string    `json:"surname"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Patronymic string    `json:"patronymic,omitempty"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type UserFilter struct {
	SurnameFilter    string
	NameFilter       string
	AddressFilter    string
	PatronymicFilter string
}
