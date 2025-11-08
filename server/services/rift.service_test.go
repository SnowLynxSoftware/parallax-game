package services

import (
	"errors"
	"testing"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// MockRiftRepository for RiftService tests
type MockRiftRepository struct {
	mock.Mock
}

func (m *MockRiftRepository) GetAllRifts() ([]*repositories.RiftEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.RiftEntity), args.Error(1)
}

func (m *MockRiftRepository) GetRiftById(id int64) (*repositories.RiftEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.RiftEntity), args.Error(1)
}

func (m *MockRiftRepository) GetRiftsByDifficulty(difficulty string) ([]*repositories.RiftEntity, error) {
	args := m.Called(difficulty)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.RiftEntity), args.Error(1)
}

// MockExpeditionRepository for RiftService tests
type MockExpeditionRepository struct {
	mock.Mock
}

func (m *MockExpeditionRepository) CreateExpedition(userId, teamId, riftId int64, durationMinutes int) (*repositories.ExpeditionEntity, error) {
	args := m.Called(userId, teamId, riftId, durationMinutes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepository) GetExpeditionById(expeditionId int64) (*repositories.ExpeditionEntity, error) {
	args := m.Called(expeditionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepository) GetActiveExpeditionsByUserId(userId int64) ([]*repositories.ExpeditionEntity, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepository) GetCompletedExpeditionsByUserId(userId int64, limit int) ([]*repositories.ExpeditionEntity, error) {
	args := m.Called(userId, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepository) GetCompletedExpeditionsCount(userId int64) (int, error) {
	args := m.Called(userId)
	return args.Int(0), args.Error(1)
}

func (m *MockExpeditionRepository) MarkCompleted(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

func (m *MockExpeditionRepository) MarkProcessed(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

func (m *MockExpeditionRepository) MarkClaimed(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

// Test GetAllRifts - Success With Zero Completed Expeditions
func TestRiftService_GetAllRifts_ZeroCompletedExpeditions(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rifts := []*repositories.RiftEntity{
		{ID: 1, Name: "Tutorial Rift", Difficulty: "tutorial", Description: "Test", WorldType: "neon", DurationMinutes: 10, WeakToElement: "fire", UnlockRequirementText: stringPtr("Always available"), Icon: "icon.png"},
		{ID: 2, Name: "Easy Rift", Difficulty: "easy", Description: "Test", WorldType: "crimson", DurationMinutes: 20, WeakToElement: "ice", UnlockRequirementText: stringPtr("Complete tutorial"), Icon: "icon.png"},
		{ID: 3, Name: "Medium Rift", Difficulty: "medium", Description: "Test", WorldType: "frozen", DurationMinutes: 30, WeakToElement: "lightning", UnlockRequirementText: stringPtr("5 completions"), Icon: "icon.png"},
	}

	mockRiftRepo.On("GetAllRifts").Return(rifts, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(0, nil)

	// Act
	result, err := service.GetAllRifts(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.True(t, result[0].IsUnlocked)  // Tutorial always unlocked
	assert.True(t, result[1].IsUnlocked)  // Easy always unlocked
	assert.False(t, result[2].IsUnlocked) // Medium requires 5 completions
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test GetAllRifts - Success With 5+ Completed Expeditions
func TestRiftService_GetAllRifts_FiveCompletedExpeditions(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rifts := []*repositories.RiftEntity{
		{ID: 1, Name: "Tutorial Rift", Difficulty: "tutorial", Description: "Test", WorldType: "neon", DurationMinutes: 10, WeakToElement: "fire", UnlockRequirementText: stringPtr("Always available"), Icon: "icon.png"},
		{ID: 2, Name: "Medium Rift", Difficulty: "medium", Description: "Test", WorldType: "frozen", DurationMinutes: 30, WeakToElement: "lightning", UnlockRequirementText: stringPtr("5 completions"), Icon: "icon.png"},
		{ID: 3, Name: "Hard Rift", Difficulty: "hard", Description: "Test", WorldType: "verdant", DurationMinutes: 45, WeakToElement: "shadow", UnlockRequirementText: stringPtr("15 completions"), Icon: "icon.png"},
	}

	mockRiftRepo.On("GetAllRifts").Return(rifts, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(5, nil)

	// Act
	result, err := service.GetAllRifts(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.True(t, result[0].IsUnlocked)  // Tutorial always unlocked
	assert.True(t, result[1].IsUnlocked)  // Medium unlocked at 5
	assert.False(t, result[2].IsUnlocked) // Hard requires 15
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test GetAllRifts - Rift Repository Error
func TestRiftService_GetAllRifts_RiftRepositoryError(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	mockRiftRepo.On("GetAllRifts").Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetAllRifts(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRiftRepo.AssertExpectations(t)
}

// Test GetAllRifts - Expedition Repository Error
func TestRiftService_GetAllRifts_ExpeditionRepositoryError(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rifts := []*repositories.RiftEntity{
		{ID: 1, Name: "Tutorial Rift", Difficulty: "tutorial", Description: "Test", WorldType: "neon", DurationMinutes: 10, WeakToElement: "fire", UnlockRequirementText: stringPtr("Always available"), Icon: "icon.png"},
	}

	mockRiftRepo.On("GetAllRifts").Return(rifts, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(0, errors.New("database error"))

	// Act
	result, err := service.GetAllRifts(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test GetRiftById - Success
func TestRiftService_GetRiftById_Success(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:                    1,
		Name:                  "Tutorial Rift",
		Description:           "Test description",
		WorldType:             "neon",
		DurationMinutes:       10,
		Difficulty:            "tutorial",
		WeakToElement:         "fire",
		UnlockRequirementText: stringPtr("Always available"),
		Icon:                  "icon.png",
	}

	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)

	// Act
	result, err := service.GetRiftById(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Tutorial Rift", result.Name)
	assert.True(t, result.IsUnlocked) // Always unlocked when querying by ID
	mockRiftRepo.AssertExpectations(t)
}

// Test GetRiftById - Not Found
func TestRiftService_GetRiftById_NotFound(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	mockRiftRepo.On("GetRiftById", int64(999)).Return(nil, errors.New("rift not found"))

	// Act
	result, err := service.GetRiftById(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRiftRepo.AssertExpectations(t)
}

// Test GetRiftById - Repository Error
func TestRiftService_GetRiftById_RepositoryError(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	mockRiftRepo.On("GetRiftById", int64(1)).Return(nil, errors.New("database connection failed"))

	// Act
	result, err := service.GetRiftById(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRiftRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - No Requirement (Always Unlocked)
func TestRiftService_IsRiftUnlockedForUser_NoRequirement(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:         1,
		Difficulty: "tutorial",
	}

	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(0, nil)

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 1)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - Requirement Met
func TestRiftService_IsRiftUnlockedForUser_RequirementMet(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:         2,
		Difficulty: "medium",
	}

	mockRiftRepo.On("GetRiftById", int64(2)).Return(rift, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(10, nil) // 10 > 5

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 2)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - Requirement Not Met
func TestRiftService_IsRiftUnlockedForUser_RequirementNotMet(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:         3,
		Difficulty: "hard",
	}

	mockRiftRepo.On("GetRiftById", int64(3)).Return(rift, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(5, nil) // 5 < 15

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 3)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - Rift Not Found
func TestRiftService_IsRiftUnlockedForUser_RiftNotFound(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	mockRiftRepo.On("GetRiftById", int64(999)).Return(nil, errors.New("rift not found"))

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 999)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	mockRiftRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - Expedition Count Error
func TestRiftService_IsRiftUnlockedForUser_ExpeditionCountError(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:         2,
		Difficulty: "medium",
	}

	mockRiftRepo.On("GetRiftById", int64(2)).Return(rift, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(0, errors.New("database error"))

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 2)

	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

// Test IsRiftUnlockedForUser - Legendary Difficulty
func TestRiftService_IsRiftUnlockedForUser_LegendaryDifficulty(t *testing.T) {
	// Arrange
	mockRiftRepo := new(MockRiftRepository)
	mockExpeditionRepo := new(MockExpeditionRepository)
	service := NewRiftService(mockRiftRepo, mockExpeditionRepo)

	rift := &repositories.RiftEntity{
		ID:         4,
		Difficulty: "legendary",
	}

	mockRiftRepo.On("GetRiftById", int64(4)).Return(rift, nil)
	mockExpeditionRepo.On("GetCompletedExpeditionsCount", int64(1)).Return(35, nil) // 35 >= 30

	// Act
	result, err := service.IsRiftUnlockedForUser(1, 4)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}
