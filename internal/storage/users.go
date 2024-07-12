package storage

import (
	"context"
	"github.com/google/uuid"
	"timeTracker/internal/model"
	"timeTracker/internal/utils"
)

func (s *Storage) GetUsers(ctx context.Context, filters model.UserFilter, pagination utils.Pagination) ([]model.User, error) {
	var users []model.User

	err := s.db.Scopes(utils.Paginate(users, &pagination, s.db)).Where(&model.User{Name: filters.NameFilter, Surname: filters.SurnameFilter, Address: filters.AddressFilter, Patronymic: filters.PatronymicFilter}).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) AddUser(ctx context.Context, user model.User) (model.User, error) {
	err := s.db.Create(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	err := s.db.Save(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	err := s.db.Where("id = ?", userID).Delete(&model.User{}).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
