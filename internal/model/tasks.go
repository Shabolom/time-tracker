package model

import "github.com/google/uuid"

// Task is the db schema for the task table
type Task struct {
	Name   string `json:"name" db:"name"`
	UserID uuid.UUID
	User   User
	Base
}
