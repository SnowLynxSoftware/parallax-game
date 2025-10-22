package services

import (
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserInventoryRepository for InventoryService tests
type MockUserInventoryRepository struct {
	mock.Mock
}

func (m *MockUserInventoryRepository) GetInventoryByUserId(userId int64) ([]*repositories.UserInventoryEntity, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) GetInventoryById(inventoryId int64) (*repositories.UserInventoryEntity, error) {
	args := m.Called(inventoryId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) GetEquipmentByUserId(userId int64) ([]*repositories.UserInventoryEntity, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) GetConsumablesByUserId(userId int64) ([]*repositories.UserInventoryEntity, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) GetInventoryByUserAndItem(userId int64, lootItemId int64) (*repositories.UserInventoryEntity, error) {
	args := m.Called(userId, lootItemId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) AddLoot(userId, lootItemId int64, itemType string) (*repositories.UserInventoryEntity, error) {
	args := m.Called(userId, lootItemId, itemType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repositories.UserInventoryEntity), args.Error(1)
}

func (m *MockUserInventoryRepository) ConsumeLoot(inventoryId int64) error {
	args := m.Called(inventoryId)
	return args.Error(0)
}

// Test GetUserInventory - Success With Both Types
func TestInventoryService_GetUserInventory_Success(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	equipment := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 10, Quantity: 1, AcquiredAt: now},
	}
	consumables := []*repositories.UserInventoryEntity{
		{ID: 2, UserID: 1, LootItemID: 20, Quantity: 5, AcquiredAt: now},
	}

	lootItem1 := &repositories.LootItemEntity{ID: 10, Name: "Sword", ItemType: "equipment", EquipmentSlot: stringPtr("weapon")}
	lootItem2 := &repositories.LootItemEntity{ID: 20, Name: "Potion", ItemType: "consumable"}

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(equipment, nil)
	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return(consumables, nil)
	mockLootItemRepo.On("GetLootItemById", int64(10)).Return(lootItem1, nil)
	mockLootItemRepo.On("GetLootItemById", int64(20)).Return(lootItem2, nil)
	mockTeamRepo.On("GetTeamsByUserIdWithSlot", int64(1), int64(1)).Return(nil, nil, errors.New("not equipped"))

	// Act
	result, err := service.GetUserInventory(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Equipment))
	assert.Equal(t, 1, len(result.Consumables))
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
}

// Test GetUserInventory - Empty Inventory
func TestInventoryService_GetUserInventory_EmptyInventory(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return([]*repositories.UserInventoryEntity{}, nil)
	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return([]*repositories.UserInventoryEntity{}, nil)

	// Act
	result, err := service.GetUserInventory(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Equipment))
	assert.Equal(t, 0, len(result.Consumables))
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetUserInventory - Equipment Error
func TestInventoryService_GetUserInventory_EquipmentError(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetUserInventory(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetUserInventory - Consumables Error
func TestInventoryService_GetUserInventory_ConsumablesError(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return([]*repositories.UserInventoryEntity{}, nil)
	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetUserInventory(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetEquipment - Items Equipped
func TestInventoryService_GetEquipment_ItemsEquipped(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	equipment := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 10, Quantity: 1, AcquiredAt: now},
	}
	lootItem := &repositories.LootItemEntity{ID: 10, Name: "Sword", ItemType: "equipment", EquipmentSlot: stringPtr("weapon")}

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(equipment, nil)
	mockLootItemRepo.On("GetLootItemById", int64(10)).Return(lootItem, nil)
	slot := "weapon"
	mockTeamRepo.On("GetTeamsByUserIdWithSlot", int64(1), int64(1)).Return(&repositories.TeamEntity{ID: 1}, &slot, nil)

	// Act
	result, err := service.GetEquipment(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.True(t, result[0].IsEquipped)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

// Test GetEquipment - Items Not Equipped
func TestInventoryService_GetEquipment_ItemsNotEquipped(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	equipment := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 10, Quantity: 1, AcquiredAt: now},
	}
	lootItem := &repositories.LootItemEntity{ID: 10, Name: "Sword", ItemType: "equipment", EquipmentSlot: stringPtr("weapon")}

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(equipment, nil)
	mockLootItemRepo.On("GetLootItemById", int64(10)).Return(lootItem, nil)
	mockTeamRepo.On("GetTeamsByUserIdWithSlot", int64(1), int64(1)).Return(nil, nil, errors.New("not equipped"))

	// Act
	result, err := service.GetEquipment(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.False(t, result[0].IsEquipped)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

// Test GetEquipment - Mixed Equipped Status
func TestInventoryService_GetEquipment_MixedEquippedStatus(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	equipment := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 10, Quantity: 1, AcquiredAt: now},
		{ID: 2, UserID: 1, LootItemID: 11, Quantity: 1, AcquiredAt: now},
	}
	lootItem1 := &repositories.LootItemEntity{ID: 10, Name: "Sword", ItemType: "equipment", EquipmentSlot: stringPtr("weapon")}
	lootItem2 := &repositories.LootItemEntity{ID: 11, Name: "Shield", ItemType: "equipment", EquipmentSlot: stringPtr("armor")}

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(equipment, nil)
	mockLootItemRepo.On("GetLootItemById", int64(10)).Return(lootItem1, nil)
	mockLootItemRepo.On("GetLootItemById", int64(11)).Return(lootItem2, nil)
	slot := "weapon"
	mockTeamRepo.On("GetTeamsByUserIdWithSlot", int64(1), int64(1)).Return(&repositories.TeamEntity{ID: 1}, &slot, nil)
	mockTeamRepo.On("GetTeamsByUserIdWithSlot", int64(1), int64(2)).Return(nil, nil, errors.New("not equipped"))

	// Act
	result, err := service.GetEquipment(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.True(t, result[0].IsEquipped)
	assert.False(t, result[1].IsEquipped)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
}

// Test GetEquipment - Inventory Repository Error
func TestInventoryService_GetEquipment_InventoryRepositoryError(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetEquipment(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetEquipment - Loot Item Not Found
func TestInventoryService_GetEquipment_LootItemNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	equipment := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 10, Quantity: 1, AcquiredAt: now},
	}

	mockInventoryRepo.On("GetEquipmentByUserId", int64(1)).Return(equipment, nil)
	mockLootItemRepo.On("GetLootItemById", int64(10)).Return(nil, errors.New("item not found"))

	// Act
	result, err := service.GetEquipment(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
}

// Test GetConsumables - Various Quantities
func TestInventoryService_GetConsumables_VariousQuantities(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	consumables := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 20, Quantity: 5, AcquiredAt: now},
		{ID: 2, UserID: 1, LootItemID: 21, Quantity: 10, AcquiredAt: now},
	}
	lootItem1 := &repositories.LootItemEntity{ID: 20, Name: "Potion", ItemType: "consumable"}
	lootItem2 := &repositories.LootItemEntity{ID: 21, Name: "Elixir", ItemType: "consumable"}

	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return(consumables, nil)
	mockLootItemRepo.On("GetLootItemById", int64(20)).Return(lootItem1, nil)
	mockLootItemRepo.On("GetLootItemById", int64(21)).Return(lootItem2, nil)

	// Act
	result, err := service.GetConsumables(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 5, result[0].Quantity)
	assert.Equal(t, 10, result[1].Quantity)
	assert.False(t, result[0].IsEquipped) // Consumables never equipped
	assert.False(t, result[1].IsEquipped)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
}

// Test GetConsumables - Empty
func TestInventoryService_GetConsumables_Empty(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return([]*repositories.UserInventoryEntity{}, nil)

	// Act
	result, err := service.GetConsumables(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetConsumables - Inventory Repository Error
func TestInventoryService_GetConsumables_InventoryRepositoryError(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return(nil, errors.New("database error"))

	// Act
	result, err := service.GetConsumables(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
}

// Test GetConsumables - Loot Item Not Found
func TestInventoryService_GetConsumables_LootItemNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockUserInventoryRepository)
	mockLootItemRepo := new(MockLootItemRepository)
	mockTeamRepo := new(MockTeamRepository)
	service := NewInventoryService(mockInventoryRepo, mockLootItemRepo, mockTeamRepo)

	now := time.Now()
	consumables := []*repositories.UserInventoryEntity{
		{ID: 1, UserID: 1, LootItemID: 20, Quantity: 5, AcquiredAt: now},
	}

	mockInventoryRepo.On("GetConsumablesByUserId", int64(1)).Return(consumables, nil)
	mockLootItemRepo.On("GetLootItemById", int64(20)).Return(nil, errors.New("item not found"))

	// Act
	result, err := service.GetConsumables(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockInventoryRepo.AssertExpectations(t)
	mockLootItemRepo.AssertExpectations(t)
}
