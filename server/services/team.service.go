package services

import (
	"fmt"
	"time"

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
	teamRepository       repositories.ITeamRepository
	inventoryRepository  repositories.IUserInventoryRepository
	lootItemRepository   repositories.ILootItemRepository
	expeditionRepository repositories.IExpeditionRepository
	riftRepository       repositories.IRiftRepository
	gameCoreService      IGameCoreService
}

func NewTeamService(
	teamRepository repositories.ITeamRepository,
	inventoryRepository repositories.IUserInventoryRepository,
	lootItemRepository repositories.ILootItemRepository,
	expeditionRepository repositories.IExpeditionRepository,
	riftRepository repositories.IRiftRepository,
	gameCoreService IGameCoreService,
) ITeamService {
	return &TeamService{
		teamRepository:       teamRepository,
		inventoryRepository:  inventoryRepository,
		lootItemRepository:   lootItemRepository,
		expeditionRepository: expeditionRepository,
		riftRepository:       riftRepository,
		gameCoreService:      gameCoreService,
	}
}

func (s *TeamService) GetUserTeams(userId int64) ([]*models.TeamResponseDTO, error) {
	teams, err := s.teamRepository.GetTeamsByUserId(userId)
	if err != nil {
		return nil, err
	}

	// Get active expeditions for this user
	activeExpeditions, err := s.expeditionRepository.GetActiveExpeditionsByUserId(userId)
	if err != nil {
		return nil, err
	}

	// Create map of team_id -> expedition for quick lookup
	expeditionByTeam := make(map[int64]*repositories.ExpeditionEntity)
	for _, expedition := range activeExpeditions {
		expeditionByTeam[expedition.TeamID] = expedition
	}

	// Get total completed expeditions for unlock requirements
	completedCount, err := s.expeditionRepository.GetCompletedExpeditionsCount(userId)
	if err != nil {
		completedCount = 0 // Default to 0 if error
	}

	response := make([]*models.TeamResponseDTO, len(teams))
	for i, team := range teams {
		teamDTO, err := s.mapTeamToDTO(team)
		if err != nil {
			return nil, err
		}

		// Check if team is on expedition
		if expedition, exists := expeditionByTeam[team.ID]; exists {
			teamDTO.OnExpedition = true

			// Get rift details
			rift, err := s.riftRepository.GetRiftById(expedition.RiftID)
			if err == nil {
				// Map expedition to DTO
				teamDTO.ExpeditionData = s.mapExpeditionToDTO(expedition, rift.Name, team.TeamNumber)
			}
		}

		// Set unlock requirement for locked teams
		if !team.IsUnlocked {
			requirement := s.getUnlockRequirement(team.TeamNumber, completedCount)
			teamDTO.UnlockRequirement = &requirement
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

func (s *TeamService) mapExpeditionToDTO(expedition *repositories.ExpeditionEntity, riftName string, teamNumber int) *models.ExpeditionResponseDTO {
	completionTime := expedition.StartTime.Add(time.Duration(expedition.DurationMinutes) * time.Minute)

	var timeRemaining *int
	if !expedition.Completed {
		remaining := int(time.Until(completionTime).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		timeRemaining = &remaining
	}

	return &models.ExpeditionResponseDTO{
		ID:              expedition.ID,
		TeamID:          expedition.TeamID,
		TeamNumber:      teamNumber,
		RiftID:          expedition.RiftID,
		RiftName:        riftName,
		StartTime:       expedition.StartTime.Format("2006-01-02T15:04:05Z"),
		DurationMinutes: expedition.DurationMinutes,
		CompletionTime:  completionTime.Format("2006-01-02T15:04:05Z"),
		TimeRemaining:   timeRemaining,
		IsCompleted:     expedition.Completed,
		IsClaimed:       expedition.Claimed,
	}
}

func (s *TeamService) getUnlockRequirement(teamNumber int, completedExpeditions int) string {
	switch teamNumber {
	case 2:
		return "Complete 1 expedition to unlock"
	case 3:
		return "Complete 3 expeditions to unlock"
	case 4:
		return "Complete 25 expeditions to unlock"
	case 5:
		return "Complete 50 expeditions to unlock"
	default:
		return "Locked"
	}
}
