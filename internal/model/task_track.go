package model

import (
	"github.com/google/uuid"
	"time"
)

type TaskTrack struct {
	TaskID uuid.UUID
	Task   Task
	Time   *time.Duration
	Base
}

type TaskTrackRequest struct {
	TaskID uuid.UUID `json:"task_id"`
}

type CalcTimeRequest struct {
	ID uuid.UUID `json:"id"`
}
