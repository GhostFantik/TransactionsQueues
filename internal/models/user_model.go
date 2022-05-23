package models

import "github.com/google/uuid"

type User struct {
	Id      uuid.UUID `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Balance int       `json:"balance" db:"balance"`
}
