package storage

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Initializes the postgres driver
	"gorm.io/gorm"
	"timeTracker/internal/model"
	"timeTracker/internal/utils"
)

type StorageInterface interface {
	GetUsers(ctx context.Context, filters model.UserFilter, pagination utils.Pagination) ([]model.User, error)
	AddUser(ctx context.Context, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
	TrackTime(ctx context.Context, trackModel model.TaskTrack) (model.TaskTrack, error)
	StopTrackTime(ctx context.Context, taskId uuid.UUID) (model.TaskTrack, error)
	CalcTime(ctx context.Context, userID uuid.UUID) ([]model.TaskTrack, error)
	//GetBook(ctx context.Context, id int) (model.Book, error)
	//GetBooks(ctx context.Context) ([]model.Book, error)
	//UpdateBook(ctx context.Context, book model.UpdateBookRequest) (int, error)
	//DeleteBook(ctx context.Context, id int) error
	//VerifyBookExists(ctx context.Context, id int) (bool, error)
}

// Storage contains an SQL db. Storage implements the StorageInterface.
type Storage struct {
	db *gorm.DB
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()

	if err != nil {
		return err
	}

	if err = sqlDB.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetDB() (*sql.DB, error) {
	return s.db.DB()
}
