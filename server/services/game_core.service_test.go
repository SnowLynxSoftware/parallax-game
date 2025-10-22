package services

import (
	"testing"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLootItemRepository for GameCoreService tests
type MockLootItemRepository struct {
	mock.Mock
}

func (m *MockLootItemRepository) GetLootItemById(id int64) (*repositories.LootItemEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.LootItemEntity), args.Error(1)
}

func (m *MockLootItemRepository) GetAllLootItems() ([]*repositories.LootItemEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.LootItemEntity), args.Error(1)
}

func (m *MockLootItemRepository) GetLootItemsByWorldType(worldType string) ([]*repositories.LootItemEntity, error) {
	args := m.Called(worldType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.LootItemEntity), args.Error(1)
}

func (m *MockLootItemRepository) GetLootItemsByRarityAndWorldType(rarity string, worldType string) ([]*repositories.LootItemEntity, error) {
	args := m.Called(rarity, worldType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.LootItemEntity), args.Error(1)
}

// Test CalculateTeamStats - No Equipped Items
func TestGameCoreService_CalculateTeamStats_NoEquippedItems(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		SpeedBonus: 10.0,
		LuckBonus:  5.0,
		PowerBonus: 20,
	}

	equippedItems := make(map[string]*repositories.LootItemEntity)

	// Act
	result := service.CalculateTeamStats(team, equippedItems)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 10.0, result.Speed)
	assert.Equal(t, 5.0, result.Luck)
	assert.Equal(t, 20, result.Power)
}

// Test CalculateTeamStats - All Slots Filled
func TestGameCoreService_CalculateTeamStats_AllSlotsFilled(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		SpeedBonus: 10.0,
		LuckBonus:  5.0,
		PowerBonus: 20,
	}

	equippedItems := map[string]*repositories.LootItemEntity{
		"weapon": {
			ID:         1,
			Name:       "Sword",
			SpeedBonus: 5.0,
			LuckBonus:  2.0,
			PowerBonus: 15,
		},
		"armor": {
			ID:         2,
			Name:       "Shield",
			SpeedBonus: 3.0,
			LuckBonus:  1.0,
			PowerBonus: 10,
		},
		"accessory": {
			ID:         3,
			Name:       "Ring",
			SpeedBonus: 2.0,
			LuckBonus:  3.0,
			PowerBonus: 5,
		},
	}

	// Act
	result := service.CalculateTeamStats(team, equippedItems)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 20.0, result.Speed) // 10 + 5 + 3 + 2
	assert.Equal(t, 11.0, result.Luck)  // 5 + 2 + 1 + 3
	assert.Equal(t, 50, result.Power)   // 20 + 15 + 10 + 5
}

// Test CalculateTeamStats - Partial Equipment
func TestGameCoreService_CalculateTeamStats_PartialEquipment(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		SpeedBonus: 10.0,
		LuckBonus:  5.0,
		PowerBonus: 20,
	}

	equippedItems := map[string]*repositories.LootItemEntity{
		"weapon": {
			ID:         1,
			Name:       "Sword",
			SpeedBonus: 5.0,
			LuckBonus:  2.0,
			PowerBonus: 15,
		},
		"armor": nil, // Empty slot
	}

	// Act
	result := service.CalculateTeamStats(team, equippedItems)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 15.0, result.Speed) // 10 + 5
	assert.Equal(t, 7.0, result.Luck)   // 5 + 2
	assert.Equal(t, 35, result.Power)   // 20 + 15
}

// Test CalculateTeamStats - Nil Equipped Items Map
func TestGameCoreService_CalculateTeamStats_NilEquippedItemsMap(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	team := &repositories.TeamEntity{
		ID:         1,
		UserID:     1,
		TeamNumber: 1,
		SpeedBonus: 10.0,
		LuckBonus:  5.0,
		PowerBonus: 20,
	}

	// Act
	result := service.CalculateTeamStats(team, nil)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, 10.0, result.Speed)
	assert.Equal(t, 5.0, result.Luck)
	assert.Equal(t, 20, result.Power)
}

// Test CalculateExpeditionDuration - Zero Speed Bonus
func TestGameCoreService_CalculateExpeditionDuration_ZeroSpeedBonus(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	stats := &models.TeamStatsDTO{
		Speed: 0.0,
		Luck:  5.0,
		Power: 20,
	}

	baseDuration := 60 // 60 minutes

	// Act
	result := service.CalculateExpeditionDuration(stats, baseDuration)

	// Assert
	assert.Equal(t, 60, result) // No speed bonus, full duration
}

// Test CalculateExpeditionDuration - 20% Speed Bonus
func TestGameCoreService_CalculateExpeditionDuration_TwentyPercentSpeedBonus(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	stats := &models.TeamStatsDTO{
		Speed: 20.0, // 20% speed bonus
		Luck:  5.0,
		Power: 20,
	}

	baseDuration := 60 // 60 minutes

	// Act
	result := service.CalculateExpeditionDuration(stats, baseDuration)

	// Assert
	assert.Equal(t, 48, result) // 60 * (1 - 0.20) = 48
}

// Test CalculateExpeditionDuration - 50% Speed Bonus
func TestGameCoreService_CalculateExpeditionDuration_FiftyPercentSpeedBonus(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	stats := &models.TeamStatsDTO{
		Speed: 50.0, // 50% speed bonus
		Luck:  5.0,
		Power: 20,
	}

	baseDuration := 60 // 60 minutes

	// Act
	result := service.CalculateExpeditionDuration(stats, baseDuration)

	// Assert
	assert.Equal(t, 30, result) // 60 * (1 - 0.50) = 30
}

// Test CalculateExpeditionDuration - 100% Speed Bonus (Should Cap at 5 Minutes)
func TestGameCoreService_CalculateExpeditionDuration_HundredPercentSpeedBonus(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	stats := &models.TeamStatsDTO{
		Speed: 100.0, // 100% speed bonus
		Luck:  5.0,
		Power: 20,
	}

	baseDuration := 60 // 60 minutes

	// Act
	result := service.CalculateExpeditionDuration(stats, baseDuration)

	// Assert
	assert.Equal(t, 5, result) // Should be capped at minimum 5 minutes
}

// Test CalculateExpeditionDuration - Speed Bonus That Would Go Below 5 Minutes
func TestGameCoreService_CalculateExpeditionDuration_BelowMinimum(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	stats := &models.TeamStatsDTO{
		Speed: 95.0, // 95% speed bonus
		Luck:  5.0,
		Power: 20,
	}

	baseDuration := 10 // 10 minutes

	// Act
	result := service.CalculateExpeditionDuration(stats, baseDuration)

	// Assert
	assert.Equal(t, 5, result) // 10 * (1 - 0.95) = 0.5, but capped at 5
}

// Test CalculateExpeditionDuration - Nil Stats
func TestGameCoreService_CalculateExpeditionDuration_NilStats(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	baseDuration := 60

	// Act & Assert - Will panic, but that's expected behavior for nil pointer
	assert.Panics(t, func() {
		service.CalculateExpeditionDuration(nil, baseDuration)
	})
}

// Test GetElementalBonus - Matching Affinity
func TestGameCoreService_GetElementalBonus_MatchingAffinity(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	relicAffinity := "fire"
	riftWeakness := "fire"

	// Act
	result := service.GetElementalBonus(relicAffinity, riftWeakness)

	// Assert
	assert.Equal(t, 0.20, result) // 20% bonus
}

// Test GetElementalBonus - Non-Matching Affinity
func TestGameCoreService_GetElementalBonus_NonMatchingAffinity(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	relicAffinity := "fire"
	riftWeakness := "ice"

	// Act
	result := service.GetElementalBonus(relicAffinity, riftWeakness)

	// Assert
	assert.Equal(t, 0.0, result) // No bonus
}

// Test GetElementalBonus - Empty Relic Affinity
func TestGameCoreService_GetElementalBonus_EmptyRelicAffinity(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	relicAffinity := ""
	riftWeakness := "fire"

	// Act
	result := service.GetElementalBonus(relicAffinity, riftWeakness)

	// Assert
	assert.Equal(t, 0.0, result) // No bonus
}

// Test GetElementalBonus - None Relic Affinity
func TestGameCoreService_GetElementalBonus_NoneRelicAffinity(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	relicAffinity := "none"
	riftWeakness := "none"

	// Act
	result := service.GetElementalBonus(relicAffinity, riftWeakness)

	// Assert
	assert.Equal(t, 0.0, result) // No bonus for "none"
}

// Test GetElementalBonus - Empty Rift Weakness
func TestGameCoreService_GetElementalBonus_EmptyRiftWeakness(t *testing.T) {
	// Arrange
	mockRepo := new(MockLootItemRepository)
	service := NewGameCoreService(mockRepo)

	relicAffinity := "fire"
	riftWeakness := ""

	// Act
	result := service.GetElementalBonus(relicAffinity, riftWeakness)

	// Assert
	assert.Equal(t, 0.0, result) // No bonus
}
