package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type IExpeditionService interface {
	StartExpedition(userId, teamId, riftId int64) (*models.ExpeditionResponseDTO, error)
	GetActiveExpeditions(userId int64) ([]*models.ExpeditionResponseDTO, error)
	GetExpeditionHistory(userId int64, limit int) ([]*models.ExpeditionResponseDTO, error)
	ClaimExpeditionRewards(userId, expeditionId int64) (*models.ExpeditionRewardsDTO, error)
}

type ExpeditionService struct {
	expeditionRepository     repositories.IExpeditionRepository
	expeditionLootRepository repositories.IExpeditionLootRepository
	teamRepository           repositories.ITeamRepository
	riftRepository           repositories.IRiftRepository
	inventoryRepository      repositories.IUserInventoryRepository
	lootItemRepository       repositories.ILootItemRepository
	lootDropTableRepository  repositories.ILootDropTableRepository
	gameCoreService          IGameCoreService
}

func NewExpeditionService(
	expeditionRepository repositories.IExpeditionRepository,
	expeditionLootRepository repositories.IExpeditionLootRepository,
	teamRepository repositories.ITeamRepository,
	riftRepository repositories.IRiftRepository,
	inventoryRepository repositories.IUserInventoryRepository,
	lootItemRepository repositories.ILootItemRepository,
	lootDropTableRepository repositories.ILootDropTableRepository,
	gameCoreService IGameCoreService,
) IExpeditionService {
	return &ExpeditionService{
		expeditionRepository:     expeditionRepository,
		expeditionLootRepository: expeditionLootRepository,
		teamRepository:           teamRepository,
		riftRepository:           riftRepository,
		inventoryRepository:      inventoryRepository,
		lootItemRepository:       lootItemRepository,
		lootDropTableRepository:  lootDropTableRepository,
		gameCoreService:          gameCoreService,
	}
}

func (s *ExpeditionService) StartExpedition(userId, teamId, riftId int64) (*models.ExpeditionResponseDTO, error) {
	// Validate team belongs to user
	team, err := s.teamRepository.GetTeamById(teamId)
	if err != nil {
		return nil, err
	}
	if team.UserID != userId {
		return nil, fmt.Errorf("team does not belong to user")
	}
	if !team.IsUnlocked {
		return nil, fmt.Errorf("team is not unlocked")
	}

	// Get rift details
	rift, err := s.riftRepository.GetRiftById(riftId)
	if err != nil {
		return nil, err
	}

	// Get equipped items for team
	equippedItems := make(map[string]*repositories.LootItemEntity)
	s.loadEquippedItems(team, equippedItems)

	// Calculate team stats
	totalStats := s.gameCoreService.CalculateTeamStats(team, equippedItems)

	// Apply elemental bonus if relic matches rift weakness
	var relicAffinity string
	if relic, exists := equippedItems["relic"]; exists && relic != nil {
		relicAffinity = relic.ElementalAffinity
		bonus := s.gameCoreService.GetElementalBonus(relicAffinity, rift.WeakToElement)
		if bonus > 0 {
			totalStats.Power = int(float64(totalStats.Power) * (1.0 + bonus))
		}
	}

	// Calculate actual expedition duration
	duration := s.gameCoreService.CalculateExpeditionDuration(totalStats, rift.DurationMinutes)

	// Create expedition
	expedition, err := s.expeditionRepository.CreateExpedition(userId, teamId, riftId, duration)
	if err != nil {
		return nil, err
	}

	return s.mapExpeditionToDTO(expedition, rift.Name, team.TeamNumber, false)
}

func (s *ExpeditionService) GetActiveExpeditions(userId int64) ([]*models.ExpeditionResponseDTO, error) {
	expeditions, err := s.expeditionRepository.GetActiveExpeditionsByUserId(userId)
	if err != nil {
		return nil, err
	}

	response := make([]*models.ExpeditionResponseDTO, len(expeditions))
	for i, expedition := range expeditions {
		team, err := s.teamRepository.GetTeamById(expedition.TeamID)
		if err != nil {
			return nil, err
		}

		rift, err := s.riftRepository.GetRiftById(expedition.RiftID)
		if err != nil {
			return nil, err
		}

		dto, err := s.mapExpeditionToDTO(expedition, rift.Name, team.TeamNumber, false)
		if err != nil {
			return nil, err
		}
		response[i] = dto
	}

	return response, nil
}

func (s *ExpeditionService) GetExpeditionHistory(userId int64, limit int) ([]*models.ExpeditionResponseDTO, error) {
	expeditions, err := s.expeditionRepository.GetCompletedExpeditionsByUserId(userId, limit)
	if err != nil {
		return nil, err
	}

	response := make([]*models.ExpeditionResponseDTO, len(expeditions))
	for i, expedition := range expeditions {
		team, err := s.teamRepository.GetTeamById(expedition.TeamID)
		if err != nil {
			return nil, err
		}

		rift, err := s.riftRepository.GetRiftById(expedition.RiftID)
		if err != nil {
			return nil, err
		}

		// Include loot if claimed
		includeLoot := expedition.Claimed
		dto, err := s.mapExpeditionToDTO(expedition, rift.Name, team.TeamNumber, includeLoot)
		if err != nil {
			return nil, err
		}
		response[i] = dto
	}

	return response, nil
}

func (s *ExpeditionService) ClaimExpeditionRewards(userId, expeditionId int64) (*models.ExpeditionRewardsDTO, error) {
	// Get expedition
	expedition, err := s.expeditionRepository.GetExpeditionById(expeditionId)
	if err != nil {
		return nil, err
	}
	if expedition.UserID != userId {
		return nil, fmt.Errorf("expedition does not belong to user")
	}

	// Check if already claimed
	if expedition.Claimed {
		return nil, fmt.Errorf("rewards already claimed")
	}

	// Check if expedition is complete
	completionTime := expedition.StartTime.Add(time.Duration(expedition.DurationMinutes) * time.Minute)
	if time.Now().Before(completionTime) {
		return nil, fmt.Errorf("expedition not yet complete")
	}

	// Get team and rift for loot generation
	team, err := s.teamRepository.GetTeamById(expedition.TeamID)
	if err != nil {
		return nil, err
	}

	rift, err := s.riftRepository.GetRiftById(expedition.RiftID)
	if err != nil {
		return nil, err
	}

	// Generate loot
	loot, err := s.generateLoot(expedition, team, rift)
	if err != nil {
		return nil, err
	}

	// Mark expedition as claimed
	err = s.expeditionRepository.MarkClaimed(expeditionId)
	if err != nil {
		return nil, err
	}

	// Map loot to response DTOs
	lootDTOs := make([]models.LootItemResponseDTO, len(loot))
	for i, item := range loot {
		lootDTOs[i] = models.LootItemResponseDTO{
			ID:                item.ID,
			Name:              item.Name,
			Description:       item.Description,
			Rarity:            item.Rarity,
			WorldType:         item.WorldType,
			ItemType:          item.ItemType,
			EquipmentSlot:     item.EquipmentSlot,
			SpeedBonus:        item.SpeedBonus,
			LuckBonus:         item.LuckBonus,
			PowerBonus:        item.PowerBonus,
			ElementalAffinity: item.ElementalAffinity,
			PowerValue:        item.PowerValue,
			Icon:              item.Icon,
		}
	}

	return &models.ExpeditionRewardsDTO{
		ExpeditionID: expeditionId,
		Loot:         lootDTOs,
	}, nil
}

// generateLoot creates loot based on drop tables and team stats
func (s *ExpeditionService) generateLoot(expedition *repositories.ExpeditionEntity, team *repositories.TeamEntity, rift *repositories.RiftEntity) ([]*repositories.LootItemEntity, error) {
	// Get drop tables for this rift
	dropTables, err := s.lootDropTableRepository.GetDropTablesByRiftId(rift.ID)
	if err != nil {
		return nil, err
	}

	// Calculate team luck
	equippedItems := make(map[string]*repositories.LootItemEntity)
	s.loadEquippedItems(team, equippedItems)
	totalStats := s.gameCoreService.CalculateTeamStats(team, equippedItems)
	// TODO: Apply luck modifier to shift drop rates in future enhancement
	_ = totalStats.Luck

	var generatedLoot []*repositories.LootItemEntity

	// Roll for each rarity tier
	for _, dropTable := range dropTables {
		// Base drop rate (luck modifier could be applied here)
		adjustedRate := dropTable.DropRatePercent

		// Roll 0-100
		roll := rand.Float64() * 100.0

		if roll < adjustedRate {
			// This rarity drops!
			quantity := dropTable.MinQuantity
			if dropTable.MaxQuantity > dropTable.MinQuantity {
				quantity = dropTable.MinQuantity + rand.Intn(dropTable.MaxQuantity-dropTable.MinQuantity+1)
			}

			// Roll specific items for this rarity
			for i := 0; i < quantity; i++ {
				items, err := s.lootItemRepository.GetLootItemsByRarityAndWorldType(dropTable.Rarity, rift.WorldType)
				if err != nil || len(items) == 0 {
					continue
				}

				// Pick random item from this rarity/world
				item := items[rand.Intn(len(items))]
				generatedLoot = append(generatedLoot, item)

				// Add to user inventory
				_, err = s.inventoryRepository.AddLoot(expedition.UserID, item.ID, item.ItemType)
				if err != nil {
					return nil, err
				}

				// Create expedition_loot audit record
				err = s.expeditionLootRepository.CreateExpeditionLoot(expedition.ID, item.ID, 1)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return generatedLoot, nil
}

func (s *ExpeditionService) loadEquippedItems(team *repositories.TeamEntity, equippedItems map[string]*repositories.LootItemEntity) {
	slots := []struct {
		id   *int64
		name string
	}{
		{team.EquippedWeaponSlot, "weapon"},
		{team.EquippedArmorSlot, "armor"},
		{team.EquippedAccessorySlot, "accessory"},
		{team.EquippedArtifactSlot, "artifact"},
		{team.EquippedRelicSlot, "relic"},
	}

	for _, slot := range slots {
		if slot.id != nil {
			invItem, err := s.inventoryRepository.GetInventoryById(*slot.id)
			if err == nil {
				lootItem, err := s.lootItemRepository.GetLootItemById(invItem.LootItemID)
				if err == nil {
					equippedItems[slot.name] = lootItem
				}
			}
		}
	}
}

func (s *ExpeditionService) mapExpeditionToDTO(expedition *repositories.ExpeditionEntity, riftName string, teamNumber int, includeLoot bool) (*models.ExpeditionResponseDTO, error) {
	completionTime := expedition.StartTime.Add(time.Duration(expedition.DurationMinutes) * time.Minute)

	var timeRemaining *int
	if !expedition.Completed {
		remaining := int(time.Until(completionTime).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		timeRemaining = &remaining
	}

	dto := &models.ExpeditionResponseDTO{
		ID:              expedition.ID,
		TeamID:          expedition.TeamID,
		TeamNumber:      teamNumber,
		RiftID:          expedition.RiftID,
		RiftName:        riftName,
		StartTime:       expedition.StartTime.Format("2006-01-02T15:04:05Z"),
		DurationMinutes: expedition.DurationMinutes,
		CompletionTime:  completionTime.Format("2006-01-02T15:04:05Z"),
		TimeRemaining:   timeRemaining,
		IsCompleted:     expedition.Completed || time.Now().After(completionTime),
		IsClaimed:       expedition.Claimed,
	}

	if includeLoot && expedition.Claimed {
		lootEntities, err := s.expeditionLootRepository.GetLootByExpeditionId(expedition.ID)
		if err == nil {
			lootDTOs := make([]models.LootItemResponseDTO, 0)
			for _, lootEntity := range lootEntities {
				lootItem, err := s.lootItemRepository.GetLootItemById(lootEntity.LootItemID)
				if err == nil {
					lootDTOs = append(lootDTOs, models.LootItemResponseDTO{
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
					})
				}
			}
			dto.Loot = &lootDTOs
		}
	}

	return dto, nil
}
