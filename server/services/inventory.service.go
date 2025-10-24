package services

import (
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type IInventoryService interface {
	GetUserInventory(userId int64) (*models.InventoryResponseDTO, error)
	GetEquipment(userId int64) ([]models.InventoryItemResponseDTO, error)
	GetConsumables(userId int64) ([]models.InventoryItemResponseDTO, error)
}

type InventoryService struct {
	inventoryRepository repositories.IUserInventoryRepository
	lootItemRepository  repositories.ILootItemRepository
	teamRepository      repositories.ITeamRepository
}

func NewInventoryService(
	inventoryRepository repositories.IUserInventoryRepository,
	lootItemRepository repositories.ILootItemRepository,
	teamRepository repositories.ITeamRepository,
) IInventoryService {
	return &InventoryService{
		inventoryRepository: inventoryRepository,
		lootItemRepository:  lootItemRepository,
		teamRepository:      teamRepository,
	}
}

func (s *InventoryService) GetUserInventory(userId int64) (*models.InventoryResponseDTO, error) {
	equipment, err := s.GetEquipment(userId)
	if err != nil {
		return nil, err
	}

	consumables, err := s.GetConsumables(userId)
	if err != nil {
		return nil, err
	}

	return &models.InventoryResponseDTO{
		Equipment:   equipment,
		Consumables: consumables,
	}, nil
}

func (s *InventoryService) GetEquipment(userId int64) ([]models.InventoryItemResponseDTO, error) {
	inventoryItems, err := s.inventoryRepository.GetEquipmentByUserId(userId)
	if err != nil {
		return nil, err
	}

	response := make([]models.InventoryItemResponseDTO, len(inventoryItems))
	for i, invItem := range inventoryItems {
		lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
		if err != nil {
			return nil, err
		}

		// Check if equipped and get team number
		isEquipped, teamNumber := s.getEquippedTeamInfo(userId, invItem.ID)

		response[i] = models.InventoryItemResponseDTO{
			InventoryID:          invItem.ID,
			Quantity:             invItem.Quantity,
			AcquiredAt:           invItem.AcquiredAt.Format("2006-01-02T15:04:05Z"),
			IsEquipped:           isEquipped,
			EquippedByTeamNumber: teamNumber,
			LootItem:             s.mapLootItemDTO(lootItem),
		}
	}

	return response, nil
}

func (s *InventoryService) GetConsumables(userId int64) ([]models.InventoryItemResponseDTO, error) {
	inventoryItems, err := s.inventoryRepository.GetConsumablesByUserId(userId)
	if err != nil {
		return nil, err
	}

	response := make([]models.InventoryItemResponseDTO, len(inventoryItems))
	for i, invItem := range inventoryItems {
		lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
		if err != nil {
			return nil, err
		}

		response[i] = models.InventoryItemResponseDTO{
			InventoryID: invItem.ID,
			Quantity:    invItem.Quantity,
			AcquiredAt:  invItem.AcquiredAt.Format("2006-01-02T15:04:05Z"),
			IsEquipped:  false, // Consumables are never equipped
			LootItem:    s.mapLootItemDTO(lootItem),
		}
	}

	return response, nil
}

func (s *InventoryService) isItemEquipped(userId int64, inventoryId int64) bool {
	team, slot, err := s.teamRepository.GetTeamsByUserIdWithSlot(userId, inventoryId)
	if err != nil || team == nil || slot == nil {
		return false
	}
	return true
}

func (s *InventoryService) getEquippedTeamInfo(userId int64, inventoryId int64) (bool, *int) {
	team, slot, err := s.teamRepository.GetTeamsByUserIdWithSlot(userId, inventoryId)
	if err != nil || team == nil || slot == nil {
		return false, nil
	}
	return true, &team.TeamNumber
}

func (s *InventoryService) mapLootItemDTO(lootItem *repositories.LootItemEntity) models.LootItemResponseDTO {
	return models.LootItemResponseDTO{
		ID:                lootItem.ID,
		Name:              lootItem.Name,
		Description:       lootItem.Description,
		Rarity:            lootItem.Rarity,
		WorldType:         lootItem.WorldType,
		ItemType:          lootItem.ItemType,
		EquipmentSlot:     lootItem.EquipmentSlot,
		SpeedBonus:        lootItem.SpeedBonus,
		LuckBonus:         lootItem.LuckBonus,
		PowerBonus:        lootItem.PowerBonus,
		ElementalAffinity: lootItem.ElementalAffinity,
		PowerValue:        lootItem.PowerValue,
		Icon:              lootItem.Icon,
	}
}
