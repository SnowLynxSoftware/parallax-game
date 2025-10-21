package services

import (
	"fmt"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type ITeamService interface {
	GetUserTeams(userId int64) ([]*models.TeamResponseDTO, error)
	GetTeamById(teamId int64) (*models.TeamResponseDTO, error)
	EquipItemToTeam(userId, teamId int64, slot string, inventoryId int64) (*models.TeamResponseDTO, error)
	UnequipItemFromTeam(userId, teamId int64, slot string) (*models.TeamResponseDTO, error)
	ConsumeItemOnTeam(userId, teamId, inventoryId int64) (*models.TeamResponseDTO, error)
	UnlockTeam(userId, teamId int64) error
}

type TeamService struct {
	teamRepository      repositories.ITeamRepository
	inventoryRepository repositories.IUserInventoryRepository
	lootItemRepository  repositories.ILootItemRepository
	gameCoreService     IGameCoreService
}

func NewTeamService(
	teamRepository repositories.ITeamRepository,
	inventoryRepository repositories.IUserInventoryRepository,
	lootItemRepository repositories.ILootItemRepository,
	gameCoreService IGameCoreService,
) ITeamService {
	return &TeamService{
		teamRepository:      teamRepository,
		inventoryRepository: inventoryRepository,
		lootItemRepository:  lootItemRepository,
		gameCoreService:     gameCoreService,
	}
}

func (s *TeamService) GetUserTeams(userId int64) ([]*models.TeamResponseDTO, error) {
	teams, err := s.teamRepository.GetTeamsByUserId(userId)
	if err != nil {
		return nil, err
	}

	response := make([]*models.TeamResponseDTO, len(teams))
	for i, team := range teams {
		teamDTO, err := s.mapTeamToDTO(team)
		if err != nil {
			return nil, err
		}
		response[i] = teamDTO
	}

	return response, nil
}

func (s *TeamService) GetTeamById(teamId int64) (*models.TeamResponseDTO, error) {
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}

	return s.mapTeamToDTO(team)
}

func (s *TeamService) EquipItemToTeam(userId, teamId int64, slot string, inventoryId int64) (*models.TeamResponseDTO, error) {
	// Validate team belongs to user
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}
	if team.UserID != userId {
		return nil, fmt.Errorf("team does not belong to user")
	}

	// Validate inventory item belongs to user
	invItem, err := s.inventoryRepository.GetInventoryById(inventoryId)
	if err != nil {
		return nil, err
	}
	if invItem.UserID != userId {
		return nil, fmt.Errorf("inventory item does not belong to user")
	}

	// Get loot item details to validate slot
	lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
	if err != nil {
		return nil, err
	}

	// Validate item type is equipment
	if lootItem.ItemType != "equipment" {
		return nil, fmt.Errorf("item is not equipment")
	}

	// Validate slot matches equipment slot
	if lootItem.EquipmentSlot == nil || *lootItem.EquipmentSlot != slot {
		return nil, fmt.Errorf("item does not match slot type")
	}

	// Check if item is already equipped on another team
	equippedTeam, equippedSlot, err := s.teamRepository.GetTeamsByUserIdWithSlot(userId, inventoryId)
	if err != nil {
		return nil, err
	}

	// Unequip from other team if necessary
	if equippedTeam != nil && equippedSlot != nil {
		err = s.teamRepository.UnequipItem(equippedTeam.ID, *equippedSlot)
		if err != nil {
			return nil, err
		}
	}

	// Equip to target team
	err = s.teamRepository.EquipItem(teamId, slot, &inventoryId)
	if err != nil {
		return nil, err
	}

	// Return updated team
	return s.GetTeamById(teamId)
}

func (s *TeamService) UnequipItemFromTeam(userId, teamId int64, slot string) (*models.TeamResponseDTO, error) {
	// Validate team belongs to user
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}
	if team.UserID != userId {
		return nil, fmt.Errorf("team does not belong to user")
	}

	// Unequip item
	err = s.teamRepository.UnequipItem(teamId, slot)
	if err != nil {
		return nil, err
	}

	// Return updated team
	return s.GetTeamById(teamId)
}

func (s *TeamService) ConsumeItemOnTeam(userId, teamId, inventoryId int64) (*models.TeamResponseDTO, error) {
	// Validate team belongs to user
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}
	if team.UserID != userId {
		return nil, fmt.Errorf("team does not belong to user")
	}

	// Validate inventory item belongs to user
	invItem, err := s.inventoryRepository.GetInventoryById(inventoryId)
	if err != nil {
		return nil, err
	}
	if invItem.UserID != userId {
		return nil, fmt.Errorf("inventory item does not belong to user")
	}

	// Get loot item details
	lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
	if err != nil {
		return nil, err
	}

	// Validate item type is consumable
	if lootItem.ItemType != "consumable" {
		return nil, fmt.Errorf("item is not consumable")
	}

	// Update team stats
	err = s.teamRepository.UpdateTeamStats(teamId, lootItem.SpeedBonus, lootItem.LuckBonus, lootItem.PowerBonus)
	if err != nil {
		return nil, err
	}

	// Consume the item (decrements quantity or deletes)
	err = s.inventoryRepository.ConsumeLoot(inventoryId)
	if err != nil {
		return nil, err
	}

	// Return updated team
	return s.GetTeamById(teamId)
}

func (s *TeamService) UnlockTeam(userId, teamId int64) error {
	// Validate team belongs to user
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return err
	}
	if team.UserID != userId {
		return fmt.Errorf("team does not belong to user")
	}

	return s.teamRepository.UnlockTeam(teamId)
}

// mapTeamToDTO converts TeamEntity to TeamResponseDTO with all equipment details
func (s *TeamService) mapTeamToDTO(team *repositories.TeamEntity) (*models.TeamResponseDTO, error) {
	// Get equipped items
	equippedItems := make(map[string]*repositories.LootItemEntity)

	weapon, _ := s.getEquippedItem(team.EquippedWeaponSlot)
	armor, _ := s.getEquippedItem(team.EquippedArmorSlot)
	accessory, _ := s.getEquippedItem(team.EquippedAccessorySlot)
	artifact, _ := s.getEquippedItem(team.EquippedArtifactSlot)
	relic, _ := s.getEquippedItem(team.EquippedRelicSlot)

	if weapon != nil {
		equippedItems["weapon"] = weapon
	}
	if armor != nil {
		equippedItems["armor"] = armor
	}
	if accessory != nil {
		equippedItems["accessory"] = accessory
	}
	if artifact != nil {
		equippedItems["artifact"] = artifact
	}
	if relic != nil {
		equippedItems["relic"] = relic
	}

	// Calculate base stats
	baseStats := &models.TeamStatsDTO{
		Speed: team.SpeedBonus,
		Luck:  team.LuckBonus,
		Power: team.PowerBonus,
	}

	// Calculate total stats
	totalStats := s.gameCoreService.CalculateTeamStats(team, equippedItems)

	return &models.TeamResponseDTO{
		ID:                team.ID,
		TeamNumber:        team.TeamNumber,
		IsUnlocked:        team.IsUnlocked,
		BaseStats:         *baseStats,
		TotalStats:        *totalStats,
		EquippedWeapon:    s.mapEquippedItemDTO(team.EquippedWeaponSlot, weapon),
		EquippedArmor:     s.mapEquippedItemDTO(team.EquippedArmorSlot, armor),
		EquippedAccessory: s.mapEquippedItemDTO(team.EquippedAccessorySlot, accessory),
		EquippedArtifact:  s.mapEquippedItemDTO(team.EquippedArtifactSlot, artifact),
		EquippedRelic:     s.mapEquippedItemDTO(team.EquippedRelicSlot, relic),
	}, nil
}

func (s *TeamService) getEquippedItem(inventoryId *int64) (*repositories.LootItemEntity, error) {
	if inventoryId == nil {
		return nil, nil
	}

	invItem, err := s.inventoryRepository.GetInventoryById(*inventoryId)
	if err != nil {
		return nil, nil // Item not found, treat as unequipped
	}

	lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
	if err != nil {
		return nil, nil
	}

	return lootItem, nil
}

func (s *TeamService) mapEquippedItemDTO(inventoryId *int64, lootItem *repositories.LootItemEntity) *models.EquippedItemDTO {
	if inventoryId == nil || lootItem == nil {
		return &models.EquippedItemDTO{
			InventoryID: nil,
			LootItemID:  nil,
			Name:        nil,
			Icon:        nil,
			Rarity:      nil,
			SpeedBonus:  0,
			LuckBonus:   0,
			PowerBonus:  0,
		}
	}

	return &models.EquippedItemDTO{
		InventoryID:       inventoryId,
		LootItemID:        &lootItem.ID,
		Name:              &lootItem.Name,
		Icon:              &lootItem.Icon,
		Rarity:            &lootItem.Rarity,
		SpeedBonus:        lootItem.SpeedBonus,
		LuckBonus:         lootItem.LuckBonus,
		PowerBonus:        lootItem.PowerBonus,
		ElementalAffinity: &lootItem.ElementalAffinity,
	}
}
