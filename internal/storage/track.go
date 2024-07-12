package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
	"timeTracker/internal/model"
)

func (s *Storage) TrackTime(ctx context.Context, trackModel model.TaskTrack) (model.TaskTrack, error) {
	var savedModel model.TaskTrack
	err := s.db.Where("task_id = ?", trackModel.TaskID).First(&savedModel).Error

	if err != nil && err.Error() == "record not found" {
		err = s.db.Create(&trackModel).Error
		s.db.Save(&trackModel)

		if err != nil {
			return model.TaskTrack{}, err
		}
	} else {
		return model.TaskTrack{}, errors.New("task already been started")
	}

	return trackModel, nil
}

func (s *Storage) StopTrackTime(ctx context.Context, taskId uuid.UUID) (model.TaskTrack, error) {
	var savedModel model.TaskTrack
	err := s.db.Where("task_id = ?", taskId).First(&savedModel).Error

	if err != nil && err.Error() == "record not found" {
		return model.TaskTrack{}, errors.New("no task with that id")
	} else {
		if savedModel.Time != nil {
			return model.TaskTrack{}, errors.New("task already been ended")
		}
		now := time.Now().Sub(savedModel.CreatedAt)
		savedModel.Time = &now
		err = s.db.Save(&savedModel).Error

		if err != nil {
			return model.TaskTrack{}, err
		}
	}

	return savedModel, nil
}

func (s *Storage) CalcTime(ctx context.Context, userID uuid.UUID) ([]model.TaskTrack, error) {
	var taskTrackModels []model.TaskTrack

	s.db.Raw("SELECT tt.* FROM task_tracks as tt JOIN tasks ON tt.task_id = tasks.id WHERE tasks.user_id = ? ORDER BY time DESC", userID).Scan(&taskTrackModels)

	return taskTrackModels, nil
}
