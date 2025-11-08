package services

import (
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
)

type IRiftService interface {
	GetAllRifts(userId int64) ([]*models.RiftResponseDTO, error)
	GetRiftById(riftId int64) (*models.RiftResponseDTO, error)
	IsRiftUnlockedForUser(userId, riftId int64) (bool, error)
}

type RiftService struct {
	riftRepository       repositories.IRiftRepository
	expeditionRepository repositories.IExpeditionRepository
}

func NewRiftService(riftRepository repositories.IRiftRepository, expeditionRepository repositories.IExpeditionRepository) IRiftService {
	return &RiftService{
		riftRepository:       riftRepository,
		expeditionRepository: expeditionRepository,
	}
}

func (s *RiftService) GetAllRifts(userId int64) ([]*models.RiftResponseDTO, error) {
	rifts, err := s.riftRepository.GetAllRifts()
	if err != nil {
		return nil, err
	}

	completedCount, err := s.expeditionRepository.GetCompletedExpeditionsCount(userId)
	if err != nil {
		return nil, err
	}

	response := make([]*models.RiftResponseDTO, len(rifts))
	for i, rift := range rifts {
		isUnlocked := s.checkRiftUnlock(rift, completedCount)
		response[i] = &models.RiftResponseDTO{
			ID:                    rift.ID,
			Name:                  rift.Name,
			Description:           rift.Description,
			WorldType:             rift.WorldType,
			DurationMinutes:       rift.DurationMinutes,
			Difficulty:            rift.Difficulty,
			WeakToElement:         rift.WeakToElement,
			UnlockRequirementText: rift.UnlockRequirementText,
			Icon:                  rift.Icon,
			IsUnlocked:            isUnlocked,
		}
	}

	return response, nil
}

func (s *RiftService) GetRiftById(riftId int64) (*models.RiftResponseDTO, error) {
	rift, err := s.riftRepository.GetRiftById(riftId)
	if err != nil {
		return nil, err
	}

	return &models.RiftResponseDTO{
		ID:                    rift.ID,
		Name:                  rift.Name,
		Description:           rift.Description,
		WorldType:             rift.WorldType,
		DurationMinutes:       rift.DurationMinutes,
		Difficulty:            rift.Difficulty,
		WeakToElement:         rift.WeakToElement,
		UnlockRequirementText: rift.UnlockRequirementText,
		Icon:                  rift.Icon,
		IsUnlocked:            true, // Assume unlocked if querying by ID
	}, nil
}

func (s *RiftService) IsRiftUnlockedForUser(userId, riftId int64) (bool, error) {
	rift, err := s.riftRepository.GetRiftById(riftId)
	if err != nil {
		return false, err
	}

	completedCount, err := s.expeditionRepository.GetCompletedExpeditionsCount(userId)
	if err != nil {
		return false, err
	}

	return s.checkRiftUnlock(rift, completedCount), nil
}

// checkRiftUnlock determines if a rift is unlocked based on completed expeditions
func (s *RiftService) checkRiftUnlock(rift *repositories.RiftEntity, completedCount int) bool {
	// Tutorial and first rift always unlocked
	if rift.Difficulty == "tutorial" || rift.ID == 1 {
		return true
	}

	// Based on unlock requirements in seed data:
	// Easy (Crimson Wastes): Tutorial completed (implicitly true if user exists)
	// Medium (Frozen/Neon): 5 completions
	// Hard (Verdant): 15 completions
	// Legendary (Void/Light): 30 completions

	switch rift.Difficulty {
	case "easy":
		return true // Always unlocked after tutorial
	case "medium":
		return completedCount >= 5
	case "hard":
		return completedCount >= 15
	case "legendary":
		return completedCount >= 30
	default:
		return false
	}
}
