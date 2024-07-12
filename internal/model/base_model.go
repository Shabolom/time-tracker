package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Base creates the default model that every other model is based on.
type Base struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
