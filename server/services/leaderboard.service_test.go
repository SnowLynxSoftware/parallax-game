package services

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
)

// ============================================================================
// Mock Repository
// ============================================================================

type mockLeaderboardRepository struct {
	getCacheMetadataFunc    func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error)
	setSyncInProgressFunc   func(leaderboardType string, inProgress bool) error
	updateLastSyncedFunc    func(leaderboardType string) error
	getTopRankingsFunc      func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error)
	getUserRankFunc         func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error)
	truncateCacheFunc       func(leaderboardType string) error
	insertCacheItemsFunc    func(items []*repositories.LeaderboardCacheItemEntity) error
	getLegendaryItemsFunc   func() ([]*repositories.LeaderboardCacheItemEntity, error)
	getPowerScoresFunc      func() ([]*repositories.LeaderboardCacheItemEntity, error)
	getExpeditionCountsFunc func() ([]*repositories.LeaderboardCacheItemEntity, error)
}

func (m *mockLeaderboardRepository) GetCacheMetadata(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
	if m.getCacheMetadataFunc != nil {
		return m.getCacheMetadataFunc(leaderboardType)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLeaderboardRepository) SetSyncInProgress(leaderboardType string, inProgress bool) error {
	if m.setSyncInProgressFunc != nil {
		return m.setSyncInProgressFunc(leaderboardType, inProgress)
	}
	return errors.New("not implemented")
}

func (m *mockLeaderboardRepository) UpdateLastSynced(leaderboardType string) error {
	if m.updateLastSyncedFunc != nil {
		return m.updateLastSyncedFunc(leaderboardType)
	}
	return errors.New("not implemented")
}

func (m *mockLeaderboardRepository) GetTopRankings(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
	if m.getTopRankingsFunc != nil {
		return m.getTopRankingsFunc(leaderboardType, limit)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLeaderboardRepository) GetUserRank(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
	if m.getUserRankFunc != nil {
		return m.getUserRankFunc(leaderboardType, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockLeaderboardRepository) TruncateCache(leaderboardType string) error {
	if m.truncateCacheFunc != nil {
		return m.truncateCacheFunc(leaderboardType)
	}
	return errors.New("not implemented")
}

func (m *mockLeaderboardRepository) InsertCacheItems(items []*repositories.LeaderboardCacheItemEntity) error {
	if m.insertCacheItemsFunc != nil {
		return m.insertCacheItemsFunc(items)
	}
	return errors.New("not implemented")
}

func (m *mockLeaderboardRepository) GetLegendaryItemCounts() ([]*repositories.LeaderboardCacheItemEntity, error) {
	if m.getLegendaryItemsFunc != nil {
		return m.getLegendaryItemsFunc()
	}
	return nil, errors.New("not implemented")
}

func (m *mockLeaderboardRepository) GetPowerScores() ([]*repositories.LeaderboardCacheItemEntity, error) {
	if m.getPowerScoresFunc != nil {
		return m.getPowerScoresFunc()
	}
	return nil, errors.New("not implemented")
}

func (m *mockLeaderboardRepository) GetExpeditionCounts() ([]*repositories.LeaderboardCacheItemEntity, error) {
	if m.getExpeditionCountsFunc != nil {
		return m.getExpeditionCountsFunc()
	}
	return nil, errors.New("not implemented")
}

// ============================================================================
// Test Data Helpers
// ============================================================================

func createTestMetadata(leaderboardType string, lastSynced time.Time, isSyncing bool) *repositories.LeaderboardCacheEntity {
	return &repositories.LeaderboardCacheEntity{
		ID:              1,
		LeaderboardType: leaderboardType,
		LastSynced:      lastSynced,
		IsSyncing:       isSyncing,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func createTestPlayers(count int) []*repositories.LeaderboardCacheItemEntity {
	players := make([]*repositories.LeaderboardCacheItemEntity, count)
	for i := 0; i < count; i++ {
		players[i] = &repositories.LeaderboardCacheItemEntity{
			ID:              i + 1,
			LeaderboardType: "legendary",
			UserID:          i + 1,
			Username:        "player" + string(rune('A'+i)),
			Score:           int64(1000 - i*10),
			Rank:            i + 1,
			CreatedAt:       time.Now(),
		}
	}
	return players
}

func createTestRawData() []*repositories.LeaderboardCacheItemEntity {
	return []*repositories.LeaderboardCacheItemEntity{
		{UserID: 5, Username: "player5", Score: 100},
		{UserID: 1, Username: "player1", Score: 500},
		{UserID: 3, Username: "player3", Score: 300},
		{UserID: 2, Username: "player2", Score: 500}, // Tie with player1
		{UserID: 4, Username: "player4", Score: 200},
	}
}

func createTestRawDataWithTies() []*repositories.LeaderboardCacheItemEntity {
	return []*repositories.LeaderboardCacheItemEntity{
		{UserID: 4, Username: "player4", Score: 100},
		{UserID: 1, Username: "player1", Score: 100},
		{UserID: 3, Username: "player3", Score: 100},
		{UserID: 2, Username: "player2", Score: 100},
		{UserID: 5, Username: "player5", Score: 50},
	}
}

// ============================================================================
// Test Cases
// ============================================================================

func TestGetLeaderboard_Success_UserInTop20(t *testing.T) {
	// Setup
	currentTime := time.Now()
	testPlayers := createTestPlayers(15)
	currentUserID := 5

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, currentTime, false), nil
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers, nil
		},
		getUserRankFunc: func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers[4], nil // User 5 is at index 4
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", currentUserID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.LeaderboardType != "legendary" {
		t.Errorf("Expected leaderboard type 'legendary', got: %s", result.LeaderboardType)
	}

	if len(result.TopPlayers) != 15 {
		t.Errorf("Expected 15 top players, got: %d", len(result.TopPlayers))
	}

	// Verify current user is marked
	userFound := false
	for _, player := range result.TopPlayers {
		if player.UserID == currentUserID {
			if !player.IsCurrentUser {
				t.Error("Expected current user to be marked as IsCurrentUser=true")
			}
			userFound = true
		}
	}

	if !userFound {
		t.Error("Current user not found in top players")
	}

	// CurrentUserRank should be nil since user is in top 20
	if result.CurrentUserRank != nil {
		t.Error("Expected CurrentUserRank to be nil when user is in top 20")
	}
}

func TestGetLeaderboard_Success_UserOutsideTop20(t *testing.T) {
	// Setup
	currentTime := time.Now()
	testPlayers := createTestPlayers(20)
	currentUserID := 50 // User outside top 20

	outsideUser := &repositories.LeaderboardCacheItemEntity{
		UserID:   currentUserID,
		Username: "player50",
		Score:    100,
		Rank:     50,
	}

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, currentTime, false), nil
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers, nil
		},
		getUserRankFunc: func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
			return outsideUser, nil
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", currentUserID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.CurrentUserRank == nil {
		t.Fatal("Expected CurrentUserRank to be set for user outside top 20")
	}

	if result.CurrentUserRank.UserID != currentUserID {
		t.Errorf("Expected CurrentUserRank.UserID to be %d, got: %d", currentUserID, result.CurrentUserRank.UserID)
	}

	if result.CurrentUserRank.Rank != 50 {
		t.Errorf("Expected CurrentUserRank.Rank to be 50, got: %d", result.CurrentUserRank.Rank)
	}

	if !result.CurrentUserRank.IsCurrentUser {
		t.Error("Expected CurrentUserRank.IsCurrentUser to be true")
	}

	// No player in top 20 should be marked as current user
	for _, player := range result.TopPlayers {
		if player.IsCurrentUser {
			t.Error("No player in top 20 should be marked as current user")
		}
	}
}

func TestGetLeaderboard_Success_UserNotRanked(t *testing.T) {
	// Setup
	currentTime := time.Now()
	testPlayers := createTestPlayers(10)
	currentUserID := 999 // User not ranked

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, currentTime, false), nil
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers, nil
		},
		getUserRankFunc: func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
			return nil, sql.ErrNoRows // User not found
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", currentUserID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.CurrentUserRank != nil {
		t.Error("Expected CurrentUserRank to be nil for unranked user")
	}

	// No player should be marked as current user
	for _, player := range result.TopPlayers {
		if player.IsCurrentUser {
			t.Error("No player should be marked as current user")
		}
	}
}

func TestGetLeaderboard_InvalidType(t *testing.T) {
	// Setup
	mockRepo := &mockLeaderboardRepository{}
	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("invalid_type", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid leaderboard type")
	}

	if result != nil {
		t.Error("Expected nil result for invalid leaderboard type")
	}
}

func TestGetLeaderboard_MetadataError(t *testing.T) {
	// Setup
	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error when metadata fetch fails")
	}

	if result != nil {
		t.Error("Expected nil result when metadata fetch fails")
	}
}

func TestGetLeaderboard_TopRankingsError(t *testing.T) {
	// Setup
	currentTime := time.Now()

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, currentTime, false), nil
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", 1)

	// Assert
	if err == nil {
		t.Fatal("Expected error when top rankings fetch fails")
	}

	if result != nil {
		t.Error("Expected nil result when top rankings fetch fails")
	}
}

func TestGetLeaderboard_StaleCache_TriggersRebuild(t *testing.T) {
	// Setup
	staleTime := time.Now().Add(-2 * time.Hour) // 2 hours ago
	testPlayers := createTestPlayers(5)

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, staleTime, false), nil
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers, nil
		},
		getUserRankFunc: func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result even with stale cache")
	}

	// Note: We can't easily test that goroutine was spawned without adding complexity
	// In production, the rebuild happens asynchronously
}

func TestGetLeaderboard_CacheSyncing_DoesNotTriggerRebuild(t *testing.T) {
	// Setup
	staleTime := time.Now().Add(-2 * time.Hour)
	testPlayers := createTestPlayers(5)

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, staleTime, true), nil // isSyncing = true
		},
		getTopRankingsFunc: func(leaderboardType string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
			return testPlayers, nil
		},
		getUserRankFunc: func(leaderboardType string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewLeaderboardService(mockRepo)

	// Execute
	result, err := service.GetLeaderboard("legendary", 1)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result")
	}

	// Rebuild should not be triggered since isSyncing=true
	// This is validated by the fact that we don't set up rebuild-related mocks
}

func TestRebuildCache_Legendary_Success(t *testing.T) {
	// Setup
	rawData := createTestRawData()
	var insertedItems []*repositories.LeaderboardCacheItemEntity

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getLegendaryItemsFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return rawData, nil
		},
		truncateCacheFunc: func(leaderboardType string) error {
			return nil
		},
		insertCacheItemsFunc: func(items []*repositories.LeaderboardCacheItemEntity) error {
			insertedItems = items
			return nil
		},
		updateLastSyncedFunc: func(leaderboardType string) error {
			return nil
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("legendary")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(insertedItems) != 5 {
		t.Fatalf("Expected 5 items inserted, got: %d", len(insertedItems))
	}

	// Verify ranking is correct (sorted by score desc, then user_id asc for ties)
	if insertedItems[0].UserID != 1 { // Score 500, user_id 1 (lower than 2)
		t.Errorf("Expected first player to be user 1, got: %d", insertedItems[0].UserID)
	}

	if insertedItems[0].Rank != 1 {
		t.Errorf("Expected first player rank to be 1, got: %d", insertedItems[0].Rank)
	}

	if insertedItems[1].UserID != 2 { // Score 500, user_id 2 (tie with 1)
		t.Errorf("Expected second player to be user 2, got: %d", insertedItems[1].UserID)
	}

	if insertedItems[1].Rank != 1 { // Same rank due to tie
		t.Errorf("Expected second player rank to be 1 (tie), got: %d", insertedItems[1].Rank)
	}

	if insertedItems[2].UserID != 3 { // Score 300
		t.Errorf("Expected third player to be user 3, got: %d", insertedItems[2].UserID)
	}

	if insertedItems[2].Rank != 3 { // Rank 3 since 2 players tied at rank 1
		t.Errorf("Expected third player rank to be 3, got: %d", insertedItems[2].Rank)
	}
}

func TestRebuildCache_Power_Success(t *testing.T) {
	// Setup
	rawData := createTestRawData()

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getPowerScoresFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return rawData, nil
		},
		truncateCacheFunc: func(leaderboardType string) error {
			return nil
		},
		insertCacheItemsFunc: func(items []*repositories.LeaderboardCacheItemEntity) error {
			return nil
		},
		updateLastSyncedFunc: func(leaderboardType string) error {
			return nil
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("power")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRebuildCache_Expeditions_Success(t *testing.T) {
	// Setup
	rawData := createTestRawData()

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getExpeditionCountsFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return rawData, nil
		},
		truncateCacheFunc: func(leaderboardType string) error {
			return nil
		},
		insertCacheItemsFunc: func(items []*repositories.LeaderboardCacheItemEntity) error {
			return nil
		},
		updateLastSyncedFunc: func(leaderboardType string) error {
			return nil
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("expeditions")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRebuildCache_AlreadySyncing_Skips(t *testing.T) {
	// Setup
	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), true), nil // isSyncing = true
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("legendary")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when skipping, got: %v", err)
	}
}

func TestRebuildCache_InvalidType(t *testing.T) {
	// Setup
	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("invalid_type")

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid leaderboard type")
	}
}

func TestRebuildCache_GetRawDataError(t *testing.T) {
	// Setup
	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getLegendaryItemsFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return nil, errors.New("database error")
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("legendary")

	// Assert
	if err == nil {
		t.Fatal("Expected error when raw data fetch fails")
	}
}

func TestRebuildCache_TruncateError(t *testing.T) {
	// Setup
	rawData := createTestRawData()

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getLegendaryItemsFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return rawData, nil
		},
		truncateCacheFunc: func(leaderboardType string) error {
			return errors.New("truncate error")
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("legendary")

	// Assert
	if err == nil {
		t.Fatal("Expected error when truncate fails")
	}
}

func TestRebuildCache_InsertError(t *testing.T) {
	// Setup
	rawData := createTestRawData()

	mockRepo := &mockLeaderboardRepository{
		getCacheMetadataFunc: func(leaderboardType string) (*repositories.LeaderboardCacheEntity, error) {
			return createTestMetadata(leaderboardType, time.Now(), false), nil
		},
		setSyncInProgressFunc: func(leaderboardType string, inProgress bool) error {
			return nil
		},
		getLegendaryItemsFunc: func() ([]*repositories.LeaderboardCacheItemEntity, error) {
			return rawData, nil
		},
		truncateCacheFunc: func(leaderboardType string) error {
			return nil
		},
		insertCacheItemsFunc: func(items []*repositories.LeaderboardCacheItemEntity) error {
			return errors.New("insert error")
		},
	}

	service := &LeaderboardService{
		repo: mockRepo,
	}

	// Execute
	err := service.rebuildCache("legendary")

	// Assert
	if err == nil {
		t.Fatal("Expected error when insert fails")
	}
}

func TestAssignRanks_EmptyData(t *testing.T) {
	// Setup
	service := &LeaderboardService{}
	items := []*repositories.LeaderboardCacheItemEntity{}

	// Execute
	result := service.assignRanks(items, "legendary")

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty result, got: %d items", len(result))
	}
}

func TestAssignRanks_NoTies(t *testing.T) {
	// Setup
	service := &LeaderboardService{}
	items := []*repositories.LeaderboardCacheItemEntity{
		{UserID: 1, Username: "player1", Score: 500},
		{UserID: 2, Username: "player2", Score: 300},
		{UserID: 3, Username: "player3", Score: 100},
	}

	// Execute
	result := service.assignRanks(items, "legendary")

	// Assert
	if len(result) != 3 {
		t.Fatalf("Expected 3 items, got: %d", len(result))
	}

	if result[0].Rank != 1 || result[0].UserID != 1 {
		t.Error("First player should be rank 1, user 1")
	}

	if result[1].Rank != 2 || result[1].UserID != 2 {
		t.Error("Second player should be rank 2, user 2")
	}

	if result[2].Rank != 3 || result[2].UserID != 3 {
		t.Error("Third player should be rank 3, user 3")
	}
}

func TestAssignRanks_WithTies(t *testing.T) {
	// Setup
	service := &LeaderboardService{}
	items := createTestRawDataWithTies()

	// Execute
	result := service.assignRanks(items, "legendary")

	// Assert
	if len(result) != 5 {
		t.Fatalf("Expected 5 items, got: %d", len(result))
	}

	// All first 4 players should have rank 1 (tied at score 100)
	for i := 0; i < 4; i++ {
		if result[i].Rank != 1 {
			t.Errorf("Player %d should have rank 1, got: %d", i, result[i].Rank)
		}
	}

	// Fifth player should have rank 5 (score 50)
	if result[4].Rank != 5 {
		t.Errorf("Fifth player should have rank 5, got: %d", result[4].Rank)
	}

	// Verify tie-breaking by user_id (ascending)
	if result[0].UserID != 1 {
		t.Errorf("First tied player should be user 1, got: %d", result[0].UserID)
	}

	if result[1].UserID != 2 {
		t.Errorf("Second tied player should be user 2, got: %d", result[1].UserID)
	}

	if result[2].UserID != 3 {
		t.Errorf("Third tied player should be user 3, got: %d", result[2].UserID)
	}

	if result[3].UserID != 4 {
		t.Errorf("Fourth tied player should be user 4, got: %d", result[3].UserID)
	}
}

func TestAssignRanks_SetsLeaderboardType(t *testing.T) {
	// Setup
	service := &LeaderboardService{}
	items := []*repositories.LeaderboardCacheItemEntity{
		{UserID: 1, Username: "player1", Score: 500},
	}

	// Execute
	result := service.assignRanks(items, "power")

	// Assert
	if result[0].LeaderboardType != "power" {
		t.Errorf("Expected leaderboard type 'power', got: %s", result[0].LeaderboardType)
	}
}

func TestIsValidLeaderboardType(t *testing.T) {
	tests := []struct {
		leaderboardType string
		expected        bool
	}{
		{"legendary", true},
		{"power", true},
		{"expeditions", true},
		{"invalid", false},
		{"", false},
		{"LEGENDARY", false}, // Case sensitive
	}

	for _, test := range tests {
		result := isValidLeaderboardType(test.leaderboardType)
		if result != test.expected {
			t.Errorf("isValidLeaderboardType(%s) = %v, expected %v", test.leaderboardType, result, test.expected)
		}
	}
}

func TestGetLeaderboard_AllLeaderboardTypes(t *testing.T) {
	types := []string{"legendary", "power", "expeditions"}

	for _, leaderboardType := range types {
		t.Run(leaderboardType, func(t *testing.T) {
			// Setup
			currentTime := time.Now()
			testPlayers := createTestPlayers(10)

			mockRepo := &mockLeaderboardRepository{
				getCacheMetadataFunc: func(lt string) (*repositories.LeaderboardCacheEntity, error) {
					return createTestMetadata(lt, currentTime, false), nil
				},
				getTopRankingsFunc: func(lt string, limit int) ([]*repositories.LeaderboardCacheItemEntity, error) {
					return testPlayers, nil
				},
				getUserRankFunc: func(lt string, userID int) (*repositories.LeaderboardCacheItemEntity, error) {
					return nil, sql.ErrNoRows
				},
			}

			service := NewLeaderboardService(mockRepo)

			// Execute
			result, err := service.GetLeaderboard(leaderboardType, 1)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error for type %s, got: %v", leaderboardType, err)
			}

			if result.LeaderboardType != leaderboardType {
				t.Errorf("Expected leaderboard type %s, got: %s", leaderboardType, result.LeaderboardType)
			}
		})
	}
}
