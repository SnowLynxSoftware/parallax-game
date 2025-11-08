package services

import (
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsersCount(searchString string, statusFilter string, userTypeFilter string) (*int, error) {
	args := m.Called(searchString, statusFilter, userTypeFilter)
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockUserRepository) GetUsers(pageSize int, offset int, searchString string, statusFilter string, userTypeFilter string) ([]*repositories.UserEntity, error) {
	args := m.Called(pageSize, offset, searchString, statusFilter, userTypeFilter)
	return args.Get(0).([]*repositories.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id int) (*repositories.UserEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserEntity), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*repositories.UserEntity, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserEntity), args.Error(1)
}

func (m *MockUserRepository) CreateNewUser(dto *models.UserCreateDTO) (*repositories.UserEntity, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserEntity), args.Error(1)
}

func (m *MockUserRepository) MarkUserVerified(userId *int) (bool, error) {
	args := m.Called(userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(dto *models.UserUpdateDTO, userId *int) (*repositories.UserEntity, error) {
	args := m.Called(dto, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserEntity), args.Error(1)
}

func (m *MockUserRepository) UpdateUserLastLogin(userId *int) (bool, error) {
	args := m.Called(userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateUserPassword(userId *int, password string) (bool, error) {
	args := m.Called(userId, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) BanUserByIdWithReason(userId *int, reason string) (bool, error) {
	args := m.Called(userId, reason)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UnbanUserById(userId *int) (bool, error) {
	args := m.Called(userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) SetUserTypeKey(userId *int, key string) (bool, error) {
	args := m.Called(userId, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ToggleUserArchived(userId *int) error {
	args := m.Called(userId)
	return args.Error(0)
}

// Metrics methods
func (m *MockUserRepository) GetTotalUsersCount() (*int64, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockUserRepository) GetActiveUsersCount() (*int64, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockUserRepository) GetNewUsersCount() (*int64, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

// Test GetUserById - Success
func TestUserService_GetUserById_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       "test@example.com",
		DisplayName: "Test User",
		IsVerified:  true,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("GetUserById", 123).Return(expectedUser, nil)

	// Act
	result, err := userService.GetUserById(123)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
	mockRepo.AssertExpectations(t)
}

// Test GetUserById - User Not Found
func TestUserService_GetUserById_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserById", 999).Return(nil, errors.New("user not found"))

	// Act
	result, err := userService.GetUserById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Test GetUserById - Repository Returns Nil User
func TestUserService_GetUserById_RepositoryReturnsNil(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("GetUserById", 456).Return(nil, nil) // Repository returns nil user, no error

	// Act
	result, err := userService.GetUserById(456)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Success
func TestUserService_GetUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []*repositories.UserEntity{
		{
			ID:          1,
			Email:       "user1@example.com",
			DisplayName: "User One",
		},
		{
			ID:          2,
			Email:       "user2@example.com",
			DisplayName: "User Two",
		},
	}

	totalCount := 2
	pageSize := 10
	offset := 0
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return(expectedUsers, nil)
	mockRepo.On("GetUsersCount", searchString, statusFilter, userTypeFilter).Return(&totalCount, nil)

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pageSize, result.PageSize)
	assert.Equal(t, 1, result.Page) // Should be calculated as 1 for offset 0
	assert.Equal(t, totalCount, result.Total)
	assert.Len(t, result.Results, 2)
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Repository Error
func TestUserService_GetUsers_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	pageSize := 10
	offset := 0
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return([]*repositories.UserEntity(nil), errors.New("database error"))

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Count Error
func TestUserService_GetUsers_CountError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []*repositories.UserEntity{
		{ID: 1, Email: "user1@example.com"},
	}

	pageSize := 10
	offset := 0
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return(expectedUsers, nil)
	mockRepo.On("GetUsersCount", searchString, statusFilter, userTypeFilter).Return((*int)(nil), errors.New("count error"))

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "count error")
	mockRepo.AssertExpectations(t)
}

// Test UpdateUser - Success
func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userId := 123
	updateDTO := &models.UserUpdateDTO{
		Email:       "updated@example.com",
		DisplayName: "Updated User",
		AvatarURL:   nil,
		ProfileText: nil,
	}

	expectedUpdatedUser := &repositories.UserEntity{
		ID:          int64(userId),
		Email:       updateDTO.Email,
		DisplayName: updateDTO.DisplayName,
		IsVerified:  true,
		CreatedAt:   time.Now(),
	}

	mockRepo.On("UpdateUser", updateDTO, &userId).Return(expectedUpdatedUser, nil)

	// Act
	result, err := userService.UpdateUser(updateDTO, &userId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUpdatedUser.Email, result.Email)
	assert.Equal(t, expectedUpdatedUser.DisplayName, result.DisplayName)
	mockRepo.AssertExpectations(t)
}

// Test UpdateUser - User Not Found
func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userId := 999
	updateDTO := &models.UserUpdateDTO{
		Email:       "updated@example.com",
		DisplayName: "Updated User",
	}

	mockRepo.On("UpdateUser", updateDTO, &userId).Return(nil, nil)

	// Act
	result, err := userService.UpdateUser(updateDTO, &userId)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Test UpdateUser - Repository Error
func TestUserService_UpdateUser_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userId := 123
	updateDTO := &models.UserUpdateDTO{
		Email:       "updated@example.com",
		DisplayName: "Updated User",
	}

	mockRepo.On("UpdateUser", updateDTO, &userId).Return(nil, errors.New("update failed"))

	// Act
	result, err := userService.UpdateUser(updateDTO, &userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "update failed")
	mockRepo.AssertExpectations(t)
}

// Test ToggleUserArchived - Success
func TestUserService_ToggleUserArchived_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userId := 123

	mockRepo.On("ToggleUserArchived", &userId).Return(nil)

	// Act
	err := userService.ToggleUserArchived(&userId)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test ToggleUserArchived - Repository Error
func TestUserService_ToggleUserArchived_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	userId := 123

	mockRepo.On("ToggleUserArchived", &userId).Return(errors.New("archive failed"))

	// Act
	err := userService.ToggleUserArchived(&userId)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "archive failed")
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Pagination Edge Cases
func TestUserService_GetUsers_PaginationEdgeCases(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []*repositories.UserEntity{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}

	totalCount := 25
	pageSize := 10
	offset := 15 // Offset that doesn't divide evenly by pageSize
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return(expectedUsers, nil)
	mockRepo.On("GetUsersCount", searchString, statusFilter, userTypeFilter).Return(&totalCount, nil)

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pageSize, result.PageSize)
	assert.Equal(t, 2, result.Page) // offset 15 / pageSize 10 = 1, but since remainder exists, page should be 2
	assert.Equal(t, totalCount, result.Total)
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Zero Offset
func TestUserService_GetUsers_ZeroOffset(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []*repositories.UserEntity{
		{ID: 1, Email: "user1@example.com"},
	}

	totalCount := 10
	pageSize := 10
	offset := 0
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return(expectedUsers, nil)
	mockRepo.On("GetUsersCount", searchString, statusFilter, userTypeFilter).Return(&totalCount, nil)

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page) // Page should be 1 for offset 0
	mockRepo.AssertExpectations(t)
}

// Test GetUsers - Large Offset
func TestUserService_GetUsers_LargeOffset(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []*repositories.UserEntity{}

	totalCount := 100
	pageSize := 10
	offset := 100
	searchString := ""
	statusFilter := ""
	userTypeFilter := ""

	mockRepo.On("GetUsers", pageSize, offset, searchString, statusFilter, userTypeFilter).Return(expectedUsers, nil)
	mockRepo.On("GetUsersCount", searchString, statusFilter, userTypeFilter).Return(&totalCount, nil)

	// Act
	result, err := userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, result.Page) // offset 100 / pageSize 10 = 10
	assert.Len(t, result.Results, 0)
	mockRepo.AssertExpectations(t)
}
