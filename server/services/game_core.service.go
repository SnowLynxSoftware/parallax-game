package services

import (
	"math"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type IGameCoreService interface {
	CalculateTeamStats(team *repositories.TeamEntity, equippedItems map[string]*repositories.LootItemEntity) *models.TeamStatsDTO
	CalculateExpeditionDuration(baseStats *models.TeamStatsDTO, baseDuration int) int
	GetElementalBonus(relicAffinity string, riftWeakness string) float64
}

type GameCoreService struct {
	lootItemRepository repositories.ILootItemRepository
}

func NewGameCoreService(lootItemRepository repositories.ILootItemRepository) IGameCoreService {
	return &GameCoreService{
		lootItemRepository: lootItemRepository,
	}
}

// CalculateTeamStats computes total stats from base stats + equipped items
func (s *GameCoreService) CalculateTeamStats(team *repositories.TeamEntity, equippedItems map[string]*repositories.LootItemEntity) *models.TeamStatsDTO {
	stats := &models.TeamStatsDTO{
		Speed: team.SpeedBonus,
		Luck:  team.LuckBonus,
		Power: team.PowerBonus,
	}

	// Add stats from all equipped items
	for _, item := range equippedItems {
		if item != nil {
			stats.Speed += item.SpeedBonus
			stats.Luck += item.LuckBonus
			stats.Power += item.PowerBonus
		}
	}

	return stats
}

// CalculateExpeditionDuration applies speed bonus reduction (min 5 minutes)
func (s *GameCoreService) CalculateExpeditionDuration(totalStats *models.TeamStatsDTO, baseDuration int) int {
	speedModifier := totalStats.Speed / 100.0
	actualDuration := float64(baseDuration) * (1.0 - speedModifier)

	// Floor at 5 minutes
	finalDuration := int(math.Max(5.0, actualDuration))

	return finalDuration
}

// GetElementalBonus returns 0.20 if relic matches rift weakness, else 0
func (s *GameCoreService) GetElementalBonus(relicAffinity string, riftWeakness string) float64 {
	if relicAffinity != "" && relicAffinity != "none" && relicAffinity == riftWeakness {
		return 0.20 // 20% bonus
	}
	return 0.0
}
