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

// Mock ExpeditionRepository
type MockExpeditionRepositoryForExpedition struct {
	mock.Mock
}

func (m *MockExpeditionRepositoryForExpedition) CreateExpedition(userId, teamId, riftId int64, duration int) (*repositories.ExpeditionEntity, error) {
	args := m.Called(userId, teamId, riftId, duration)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) GetExpeditionById(id int64) (*repositories.ExpeditionEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) GetActiveExpeditionsByUserId(userId int64) ([]*repositories.ExpeditionEntity, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) GetCompletedExpeditionsByUserId(userId int64, limit int) ([]*repositories.ExpeditionEntity, error) {
	args := m.Called(userId, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) MarkClaimed(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

func (m *MockExpeditionRepositoryForExpedition) UpdateExpeditionCompletion(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

func (m *MockExpeditionRepositoryForExpedition) GetExpeditionsByUserIdAndRiftId(userId, riftId int64) ([]*repositories.ExpeditionEntity, error) {
	args := m.Called(userId, riftId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionEntity), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) CountCompletedExpeditionsByUserIdAndRiftId(userId, riftId int64) (int, error) {
	args := m.Called(userId, riftId)
	return args.Int(0), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) GetCompletedExpeditionsCount(userId int64) (int, error) {
	args := m.Called(userId)
	return args.Int(0), args.Error(1)
}

func (m *MockExpeditionRepositoryForExpedition) MarkCompleted(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

func (m *MockExpeditionRepositoryForExpedition) MarkProcessed(expeditionId int64) error {
	args := m.Called(expeditionId)
	return args.Error(0)
}

// Mock ExpeditionLootRepository
type MockExpeditionLootRepository struct {
	mock.Mock
}

func (m *MockExpeditionLootRepository) CreateExpeditionLoot(expeditionId, lootItemId int64, quantity int) error {
	args := m.Called(expeditionId, lootItemId, quantity)
	return args.Error(0)
}

func (m *MockExpeditionLootRepository) GetLootByExpeditionId(expeditionId int64) ([]*repositories.ExpeditionLootEntity, error) {
	args := m.Called(expeditionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.ExpeditionLootEntity), args.Error(1)
}

// Mock LootDropTableRepository
type MockLootDropTableRepository struct {
	mock.Mock
}

func (m *MockLootDropTableRepository) GetDropTablesByRiftId(riftId int64) ([]*repositories.LootDropTableEntity, error) {
	args := m.Called(riftId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.LootDropTableEntity), args.Error(1)
}

// Tests for StartExpedition
func TestExpeditionService_StartExpedition_Success(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockGameCoreService := new(MockGameCoreService)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil, // loot repo not needed for start
		mockTeamRepo,
		mockRiftRepo,
		mockInventoryRepo,
		mockLootItemRepo,
		nil, // drop table not needed for start
		mockGameCoreService,
	)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		IsUnlocked: true,
	}
	rift := &repositories.RiftEntity{
		ID:              1,
		Name:            "Test Rift",
		WorldType:       "desert",
		WeakToElement:   "water",
		DurationMinutes: 60,
	}
	stats := &models.TeamStatsDTO{Speed: 10.0, Luck: 5.0, Power: 20}
	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now(),
		DurationMinutes: 50,
		Completed:       false,
		Claimed:         false,
	}

	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockGameCoreService.On("CalculateTeamStats", team, mock.Anything).Return(stats)
	mockGameCoreService.On("CalculateExpeditionDuration", stats, 60).Return(50)
	mockExpeditionRepo.On("CreateExpedition", int64(1), int64(1), int64(1), 50).Return(expedition, nil)

	result, err := service.StartExpedition(1, 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test Rift", result.RiftName)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_StartExpedition_TeamDoesntBelongToUser(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     999, // Different user
		TeamNumber: 1,
		IsUnlocked: true,
	}

	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)

	result, err := service.StartExpedition(1, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
}

func TestExpeditionService_StartExpedition_TeamNotUnlocked(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		IsUnlocked: false, // Not unlocked
	}

	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)

	result, err := service.StartExpedition(1, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
}

func TestExpeditionService_StartExpedition_TeamRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mockTeamRepo.On("GetTeamById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.StartExpedition(1, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
}

func TestExpeditionService_StartExpedition_RiftRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		nil,
		nil,
		nil,
		nil,
	)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		IsUnlocked: true,
	}

	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.StartExpedition(1, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
}

func TestExpeditionService_StartExpedition_CreateExpeditionError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockGameCoreService := new(MockGameCoreService)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		mockInventoryRepo,
		mockLootItemRepo,
		nil,
		mockGameCoreService,
	)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		IsUnlocked: true,
	}
	rift := &repositories.RiftEntity{
		ID:              1,
		Name:            "Test Rift",
		WorldType:       "desert",
		WeakToElement:   "water",
		DurationMinutes: 60,
	}
	stats := &models.TeamStatsDTO{Speed: 10.0, Luck: 5.0, Power: 20}

	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockGameCoreService.On("CalculateTeamStats", team, mock.Anything).Return(stats)
	mockGameCoreService.On("CalculateExpeditionDuration", stats, 60).Return(50)
	mockExpeditionRepo.On("CreateExpedition", int64(1), int64(1), int64(1), 50).Return(nil, errors.New("database error"))

	result, err := service.StartExpedition(1, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

// Tests for GetActiveExpeditions
func TestExpeditionService_GetActiveExpeditions_Success(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		nil,
		nil,
		nil,
		nil,
	)

	expeditions := []*repositories.ExpeditionEntity{
		{
			ID:              1,
			UserID:          1,
			TeamID:          1,
			RiftID:          1,
			StartTime:       time.Now(),
			DurationMinutes: 50,
			Completed:       false,
			Claimed:         false,
		},
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}
	rift := &repositories.RiftEntity{ID: 1, Name: "Test Rift"}

	mockExpeditionRepo.On("GetActiveExpeditionsByUserId", int64(1)).Return(expeditions, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)

	result, err := service.GetActiveExpeditions(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Rift", result[0].RiftName)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
}

func TestExpeditionService_GetActiveExpeditions_EmptyList(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expeditions := []*repositories.ExpeditionEntity{}

	mockExpeditionRepo.On("GetActiveExpeditionsByUserId", int64(1)).Return(expeditions, nil)

	result, err := service.GetActiveExpeditions(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_GetActiveExpeditions_RepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mockExpeditionRepo.On("GetActiveExpeditionsByUserId", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.GetActiveExpeditions(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_GetActiveExpeditions_TeamRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expeditions := []*repositories.ExpeditionEntity{
		{
			ID:              1,
			UserID:          1,
			TeamID:          1,
			RiftID:          1,
			StartTime:       time.Now(),
			DurationMinutes: 50,
			Completed:       false,
			Claimed:         false,
		},
	}

	mockExpeditionRepo.On("GetActiveExpeditionsByUserId", int64(1)).Return(expeditions, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.GetActiveExpeditions(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestExpeditionService_GetActiveExpeditions_RiftRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		nil,
		nil,
		nil,
		nil,
	)

	expeditions := []*repositories.ExpeditionEntity{
		{
			ID:              1,
			UserID:          1,
			TeamID:          1,
			RiftID:          1,
			StartTime:       time.Now(),
			DurationMinutes: 50,
			Completed:       false,
			Claimed:         false,
		},
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}

	mockExpeditionRepo.On("GetActiveExpeditionsByUserId", int64(1)).Return(expeditions, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.GetActiveExpeditions(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
}

// Tests for GetExpeditionHistory
func TestExpeditionService_GetExpeditionHistory_Success(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockExpeditionLootRepo := new(MockExpeditionLootRepository)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)
	mockLootItemRepo := new(MockLootItemRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		mockExpeditionLootRepo,
		mockTeamRepo,
		mockRiftRepo,
		nil,
		mockLootItemRepo,
		nil,
		nil,
	)

	expeditions := []*repositories.ExpeditionEntity{
		{
			ID:              1,
			UserID:          1,
			TeamID:          1,
			RiftID:          1,
			StartTime:       time.Now().Add(-2 * time.Hour),
			DurationMinutes: 50,
			Completed:       true,
			Claimed:         true,
		},
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}
	rift := &repositories.RiftEntity{ID: 1, Name: "Test Rift"}

	lootEntities := []*repositories.ExpeditionLootEntity{
		{ID: 1, ExpeditionID: 1, LootItemID: 100, Quantity: 1},
	}
	lootItem := &repositories.LootItemEntity{ID: 100, Name: "Test Item", Rarity: "common"}

	mockExpeditionRepo.On("GetCompletedExpeditionsByUserId", int64(1), 10).Return(expeditions, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockExpeditionLootRepo.On("GetLootByExpeditionId", int64(1)).Return(lootEntities, nil)
	mockLootItemRepo.On("GetLootItemById", int64(100)).Return(lootItem, nil)

	result, err := service.GetExpeditionHistory(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.NotNil(t, result[0].Loot)
	assert.Len(t, *result[0].Loot, 1)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
	mockExpeditionLootRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
}

func TestExpeditionService_GetExpeditionHistory_RepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mockExpeditionRepo.On("GetCompletedExpeditionsByUserId", int64(1), 10).Return(nil, errors.New("database error"))

	result, err := service.GetExpeditionHistory(1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

// Tests for ClaimExpeditionRewards
func TestExpeditionService_ClaimExpeditionRewards_ExpeditionDoesntBelongToUser(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          999, // Different user
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         false,
	}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_AlreadyClaimed(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         true, // Already claimed
	}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_NotYetComplete(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now(), // Just started
		DurationMinutes: 50,
		Completed:       false,
		Claimed:         false,
	}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_ExpeditionRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_TeamRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         false,
	}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_RiftRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		nil,
		nil,
		nil,
		nil,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         false,
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_DropTableRepositoryError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)
	mockDropTableRepo := new(MockLootDropTableRepository)
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockGameCoreService := new(MockGameCoreService)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		mockInventoryRepo,
		mockLootItemRepo,
		mockDropTableRepo,
		mockGameCoreService,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         false,
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}
	rift := &repositories.RiftEntity{ID: 1, Name: "Test Rift", WorldType: "desert"}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockDropTableRepo.On("GetDropTablesByRiftId", int64(1)).Return(nil, errors.New("database error"))

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockRiftRepo.AssertExpectations(t)
	mockDropTableRepo.AssertExpectations(t)
}

func TestExpeditionService_ClaimExpeditionRewards_MarkClaimedError(t *testing.T) {
	mockExpeditionRepo := new(MockExpeditionRepositoryForExpedition)
	mockTeamRepo := new(MockTeamRepository)
	mockRiftRepo := new(MockRiftRepository)
	mockDropTableRepo := new(MockLootDropTableRepository)
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockGameCoreService := new(MockGameCoreService)

	service := NewExpeditionService(
		mockExpeditionRepo,
		nil,
		mockTeamRepo,
		mockRiftRepo,
		mockInventoryRepo,
		mockLootItemRepo,
		mockDropTableRepo,
		mockGameCoreService,
	)

	expedition := &repositories.ExpeditionEntity{
		ID:              1,
		UserID:          1,
		TeamID:          1,
		RiftID:          1,
		StartTime:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 50,
		Completed:       true,
		Claimed:         false,
	}
	team := &repositories.TeamEntity{ID: 1, TeamNumber: 1}
	rift := &repositories.RiftEntity{ID: 1, Name: "Test Rift", WorldType: "desert"}
	dropTables := []*repositories.LootDropTableEntity{} // Empty drop tables for simplicity
	stats := &models.TeamStatsDTO{Speed: 10.0, Luck: 5.0, Power: 20}

	mockExpeditionRepo.On("GetExpeditionById", int64(1)).Return(expedition, nil)
	mockTeamRepo.On("GetTeamById", int64(1)).Return(team, nil)
	mockRiftRepo.On("GetRiftById", int64(1)).Return(rift, nil)
	mockDropTableRepo.On("GetDropTablesByRiftId", int64(1)).Return(dropTables, nil)
	mockGameCoreService.On("CalculateTeamStats", team, mock.Anything).Return(stats)
	mockExpeditionRepo.On("MarkClaimed", int64(1)).Return(errors.New("database error"))

	result, err := service.ClaimExpeditionRewards(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockExpeditionRepo.AssertExpectations(t)
}
