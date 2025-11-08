package services

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

const (
	// CacheTTL is the duration after which the cache is considered stale
	CacheTTL = 1 * time.Hour

	// TopPlayersLimit is the number of top players to return
	TopPlayersLimit = 20

	// LeaderboardTypes
	LeaderboardTypeLegendary   = "legendary"
	LeaderboardTypePower       = "power"
	LeaderboardTypeExpeditions = "expeditions"
)

// ILeaderboardService defines the interface for leaderboard operations
type ILeaderboardService interface {
	GetLeaderboard(leaderboardType string, currentUserID int) (*models.LeaderboardResponse, error)
}

// LeaderboardService implements leaderboard business logic
type LeaderboardService struct {
	repo  repositories.ILeaderboardRepository
	mutex sync.Mutex // Prevents concurrent cache rebuilds for same type
}

// NewLeaderboardService creates a new leaderboard service
func NewLeaderboardService(repo repositories.ILeaderboardRepository) ILeaderboardService {
	return &LeaderboardService{
		repo:  repo,
		mutex: sync.Mutex{},
	}
}

// GetLeaderboard retrieves the leaderboard for a specific type and marks the current user
func (s *LeaderboardService) GetLeaderboard(leaderboardType string, currentUserID int) (*models.LeaderboardResponse, error) {
	// Validate leaderboard type
	if !isValidLeaderboardType(leaderboardType) {
		return nil, fmt.Errorf("invalid leaderboard type: %s", leaderboardType)
	}

	// Get cache metadata
	metadata, err := s.repo.GetCacheMetadata(leaderboardType)
	if err != nil {
		util.LogError(err)
		return nil, fmt.Errorf("failed to get cache metadata: %w", err)
	}

	// Check if cache is stale (older than 1 hour)
	cacheAge := time.Since(metadata.LastSynced)
	if cacheAge > CacheTTL && !metadata.IsSyncing {
		// Spawn goroutine to rebuild cache asynchronously
		go func() {
			if err := s.rebuildCache(leaderboardType); err != nil {
				util.LogError(err)
			}
		}()
	}

	// Fetch top 20 players from cache
	topPlayers, err := s.repo.GetTopRankings(leaderboardType, TopPlayersLimit)
	if err != nil {
		util.LogError(err)
		return nil, fmt.Errorf("failed to get top rankings: %w", err)
	}

	// Fetch current user's rank
	userRank, err := s.repo.GetUserRank(leaderboardType, currentUserID)
	if err != nil && err != sql.ErrNoRows {
		util.LogError(err)
		return nil, fmt.Errorf("failed to get user rank: %w", err)
	}

	// Build response
	response := &models.LeaderboardResponse{
		LeaderboardType: leaderboardType,
		LastSynced:      metadata.LastSynced.Format("2006-01-02T15:04:05Z"),
		TopPlayers:      make([]models.LeaderboardEntry, 0, len(topPlayers)),
		CurrentUserRank: nil,
	}

	// Map top players to DTOs and check if current user is in top 20
	userInTopPlayers := false
	for _, player := range topPlayers {
		isCurrentUser := player.UserID == currentUserID
		if isCurrentUser {
			userInTopPlayers = true
		}

		response.TopPlayers = append(response.TopPlayers, models.LeaderboardEntry{
			Rank:          player.Rank,
			UserID:        player.UserID,
			Username:      player.Username,
			Score:         player.Score,
			IsCurrentUser: isCurrentUser,
		})
	}

	// If user is not in top 20 but has a rank, add to CurrentUserRank
	if !userInTopPlayers && userRank != nil {
		response.CurrentUserRank = &models.LeaderboardEntry{
			Rank:          userRank.Rank,
			UserID:        userRank.UserID,
			Username:      userRank.Username,
			Score:         userRank.Score,
			IsCurrentUser: true,
		}
	}

	return response, nil
}

// rebuildCache rebuilds the cache for a specific leaderboard type
// This method is called in a goroutine and should not block the main request
func (s *LeaderboardService) rebuildCache(leaderboardType string) error {
	// Acquire mutex to prevent concurrent rebuilds
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Double-check if another goroutine already started syncing
	metadata, err := s.repo.GetCacheMetadata(leaderboardType)
	if err != nil {
		return fmt.Errorf("failed to get cache metadata: %w", err)
	}

	if metadata.IsSyncing {
		util.LogInfo("Cache rebuild already in progress, skipping")
		return nil
	}

	// Set syncing flag
	if err := s.repo.SetSyncInProgress(leaderboardType, true); err != nil {
		return fmt.Errorf("failed to set sync in progress: %w", err)
	}
	defer func() {
		// Always clear syncing flag even if rebuild fails
		if err := s.repo.SetSyncInProgress(leaderboardType, false); err != nil {
			util.LogError(err)
		}
	}()

	// Fetch raw data based on leaderboard type
	var rawData []*repositories.LeaderboardCacheItemEntity
	switch leaderboardType {
	case LeaderboardTypeLegendary:
		rawData, err = s.repo.GetLegendaryItemCounts()
	case LeaderboardTypePower:
		rawData, err = s.repo.GetPowerScores()
	case LeaderboardTypeExpeditions:
		rawData, err = s.repo.GetExpeditionCounts()
	default:
		return fmt.Errorf("invalid leaderboard type: %s", leaderboardType)
	}

	if err != nil {
		return fmt.Errorf("failed to get raw data: %w", err)
	}

	// Assign ranks with deterministic tie-breaking (by user_id)
	rankedData := s.assignRanks(rawData, leaderboardType)

	// Truncate old cache
	if err := s.repo.TruncateCache(leaderboardType); err != nil {
		return fmt.Errorf("failed to truncate cache: %w", err)
	}

	// Insert new cache items
	if err := s.repo.InsertCacheItems(rankedData); err != nil {
		return fmt.Errorf("failed to insert cache items: %w", err)
	}

	// Update last synced timestamp
	if err := s.repo.UpdateLastSynced(leaderboardType); err != nil {
		return fmt.Errorf("failed to update last synced: %w", err)
	}

	util.LogInfo("Cache rebuilt successfully")

	return nil
}

// assignRanks assigns ranks to the raw data with deterministic tie-breaking
// Tied players are sorted by user_id for consistent ordering
func (s *LeaderboardService) assignRanks(items []*repositories.LeaderboardCacheItemEntity, leaderboardType string) []*repositories.LeaderboardCacheItemEntity {
	if len(items) == 0 {
		return items
	}

	// Sort by score descending, then by user_id ascending for tie-breaking
	sort.Slice(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].UserID < items[j].UserID
		}
		return items[i].Score > items[j].Score
	})

	// Assign ranks (ties get the same rank, next player gets incremented rank)
	currentRank := 1
	for i := range items {
		if i > 0 && items[i].Score != items[i-1].Score {
			currentRank = i + 1
		}
		items[i].Rank = currentRank
		items[i].LeaderboardType = leaderboardType
	}

	return items
}

// isValidLeaderboardType checks if the leaderboard type is valid
func isValidLeaderboardType(leaderboardType string) bool {
	validTypes := map[string]bool{
		LeaderboardTypeLegendary:   true,
		LeaderboardTypePower:       true,
		LeaderboardTypeExpeditions: true,
	}
	return validTypes[leaderboardType]
}
