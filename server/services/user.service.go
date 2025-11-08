package services

import (
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type IUserService interface {
	GetUserById(id int) (*repositories.UserEntity, error)
	GetUsers(pageSize int, offset int, searchString string, statusFilter string, userTypeFilter string) (*models.PaginatedResponse, error)
	UpdateUser(dto *models.UserUpdateDTO, userId *int) (*repositories.UserEntity, error)
	ToggleUserArchived(userId *int) error
}

type UserService struct {
	userRepository repositories.IUserRepository
}

func NewUserService(userRepository repositories.IUserRepository) IUserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) ToggleUserArchived(userId *int) error {
	return s.userRepository.ToggleUserArchived(userId)
}

func (s *UserService) UpdateUser(dto *models.UserUpdateDTO, userId *int) (*repositories.UserEntity, error) {
	updatedUser, err := s.userRepository.UpdateUser(dto, userId)
	if err != nil {
		return nil, err
	}
	if updatedUser == nil {
		return nil, nil // Return nil if user not found
	}
	return updatedUser, nil
}

func (s *UserService) GetUserById(id int) (*repositories.UserEntity, error) {
	user, err := s.userRepository.GetUserById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Return nil if user not found
	}
	return user, nil
}

func (s *UserService) GetUsers(pageSize int, offset int, searchString string, statusFilter string, userTypeFilter string) (*models.PaginatedResponse, error) {
	users, err := s.userRepository.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)
	if err != nil {
		return nil, err
	}
	count, err := s.userRepository.GetUsersCount(searchString, statusFilter, userTypeFilter)
	if err != nil {
		return nil, err
	}
	results := make([]any, len(users))
	for i, user := range users {
		results[i] = user
	}

	page := offset / pageSize
	if offset%pageSize != 0 {
		page++
	}
	if page <= 0 {
		page = 1
	}

	paginatedResponse := &models.PaginatedResponse{
		PageSize: pageSize,
		Page:     page,
		Total:    *count,
		Results:  results,
	}
	return paginatedResponse, nil
}
