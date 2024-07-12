package storage

import (
	"os"
	"timeTracker/internal/model"

	_ "github.com/golang-migrate/migrate/source/file" // import file driver for migrate
	"github.com/sash20m/go-api-template/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*Storage, error) {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " dbname=" + os.Getenv("DB_NAME") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"

	logger.Log.Info(dsn, " eee")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{db: db}, nil
}

// MigratePostgres migrates the postgres db to a new version.
func (s *Storage) MigratePostgres() error {
	err := s.db.AutoMigrate(&model.User{})

	if err != nil {
		return err
	}

	err = s.db.AutoMigrate(&model.Task{})
	if err != nil {
		return err
	}

	err = s.db.AutoMigrate(&model.TaskTrack{})
	if err != nil {
		return err
	}

	logger.Log.Info("Migrations script ran successfully")

	return nil
}
