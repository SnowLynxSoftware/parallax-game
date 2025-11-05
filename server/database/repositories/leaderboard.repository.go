package repositories

import (
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type LeaderboardCacheEntity struct {
	ID              int       `db:"id"`
	LeaderboardType string    `db:"leaderboard_type"`
	LastSynced      time.Time `db:"last_synced"`
	IsSyncing       bool      `db:"is_syncing"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type LeaderboardCacheItemEntity struct {
	ID              int       `db:"id"`
	LeaderboardType string    `db:"leaderboard_type"`
	UserID          int       `db:"user_id"`
	Username        string    `db:"username"`
	Score           int64     `db:"score"`
	Rank            int       `db:"rank"`
	CreatedAt       time.Time `db:"created_at"`
}

type ILeaderboardRepository interface {
	// Cache metadata
	GetCacheMetadata(leaderboardType string) (*LeaderboardCacheEntity, error)
	SetSyncInProgress(leaderboardType string, inProgress bool) error
	UpdateLastSynced(leaderboardType string) error

	// Cache items
	GetTopRankings(leaderboardType string, limit int) ([]*LeaderboardCacheItemEntity, error)
	GetUserRank(leaderboardType string, userID int) (*LeaderboardCacheItemEntity, error)
	TruncateCache(leaderboardType string) error
	InsertCacheItems(items []*LeaderboardCacheItemEntity) error

	// Raw data queries (for building cache)
	GetLegendaryItemCounts() ([]*LeaderboardCacheItemEntity, error)
	GetPowerScores() ([]*LeaderboardCacheItemEntity, error)
	GetExpeditionCounts() ([]*LeaderboardCacheItemEntity, error)
}

type LeaderboardRepository struct {
	db *database.AppDataSource
}

func NewLeaderboardRepository(db *database.AppDataSource) ILeaderboardRepository {
	return &LeaderboardRepository{
		db: db,
	}
}

// ============================================================================
// Cache Metadata Operations
// ============================================================================

// GetCacheMetadata retrieves the metadata for a specific leaderboard type
func (r *LeaderboardRepository) GetCacheMetadata(leaderboardType string) (*LeaderboardCacheEntity, error) {
	entity := &LeaderboardCacheEntity{}
	sql := `SELECT id, leaderboard_type, last_synced, is_syncing, created_at, updated_at
	        FROM leaderboard_cache
	        WHERE leaderboard_type = $1`
	err := r.db.DB.Get(entity, sql, leaderboardType)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// SetSyncInProgress sets or clears the is_syncing flag for a leaderboard type
func (r *LeaderboardRepository) SetSyncInProgress(leaderboardType string, inProgress bool) error {
	sql := `UPDATE leaderboard_cache
	        SET is_syncing = $1, updated_at = NOW()
	        WHERE leaderboard_type = $2`
	_, err := r.db.DB.Exec(sql, inProgress, leaderboardType)
	return err
}

// UpdateLastSynced updates the last_synced timestamp for a leaderboard type
func (r *LeaderboardRepository) UpdateLastSynced(leaderboardType string) error {
	sql := `UPDATE leaderboard_cache
	        SET last_synced = NOW(), updated_at = NOW()
	        WHERE leaderboard_type = $1`
	_, err := r.db.DB.Exec(sql, leaderboardType)
	return err
}

// ============================================================================
// Cache Items Operations
// ============================================================================

// GetTopRankings retrieves the top N ranked players for a leaderboard type
func (r *LeaderboardRepository) GetTopRankings(leaderboardType string, limit int) ([]*LeaderboardCacheItemEntity, error) {
	items := []*LeaderboardCacheItemEntity{}
	sql := `SELECT id, leaderboard_type, user_id, username, score, rank, created_at
	        FROM leaderboard_cache_items
	        WHERE leaderboard_type = $1
	        ORDER BY rank ASC
	        LIMIT $2`
	err := r.db.DB.Select(&items, sql, leaderboardType, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetUserRank retrieves the rank entry for a specific user on a leaderboard type
func (r *LeaderboardRepository) GetUserRank(leaderboardType string, userID int) (*LeaderboardCacheItemEntity, error) {
	entity := &LeaderboardCacheItemEntity{}
	sql := `SELECT id, leaderboard_type, user_id, username, score, rank, created_at
	        FROM leaderboard_cache_items
	        WHERE leaderboard_type = $1 AND user_id = $2`
	err := r.db.DB.Get(entity, sql, leaderboardType, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// TruncateCache removes all cache items for a specific leaderboard type
func (r *LeaderboardRepository) TruncateCache(leaderboardType string) error {
	sql := `DELETE FROM leaderboard_cache_items WHERE leaderboard_type = $1`
	_, err := r.db.DB.Exec(sql, leaderboardType)
	return err
}

// InsertCacheItems inserts multiple cache items in a batch operation
func (r *LeaderboardRepository) InsertCacheItems(items []*LeaderboardCacheItemEntity) error {
	if len(items) == 0 {
		return nil
	}

	sql := `INSERT INTO leaderboard_cache_items (leaderboard_type, user_id, username, score, rank)
	        VALUES (:leaderboard_type, :user_id, :username, :score, :rank)
	        ON CONFLICT (leaderboard_type, user_id) DO UPDATE SET
	            score = EXCLUDED.score,
	            rank = EXCLUDED.rank`

	_, err := r.db.DB.NamedExec(sql, items)
	return err
}

// ============================================================================
// Raw Data Queries (for building cache)
// ============================================================================

// GetLegendaryItemCounts calculates total legendary items per user
// Returns users sorted by count descending
func (r *LeaderboardRepository) GetLegendaryItemCounts() ([]*LeaderboardCacheItemEntity, error) {
	items := []*LeaderboardCacheItemEntity{}
	sql := `SELECT
	        u.id as user_id,
	        u.display_name as username,
	        COALESCE(SUM(ui.quantity), 0)::bigint as score
	        FROM users u
	        LEFT JOIN user_inventory ui ON ui.user_id = u.id
	        LEFT JOIN loot_items li ON li.id = ui.loot_item_id AND li.rarity = 'legendary'
	        WHERE u.is_archived = false
	        GROUP BY u.id, u.display_name
	        HAVING COALESCE(SUM(ui.quantity), 0) > 0
	        ORDER BY score DESC`

	err := r.db.DB.Select(&items, sql)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetPowerScores calculates weighted inventory value per user
// Rarity weights: Common=1, Uncommon=5, Rare=25, Epic=125, Legendary=1000
// Returns users sorted by score descending
func (r *LeaderboardRepository) GetPowerScores() ([]*LeaderboardCacheItemEntity, error) {
	items := []*LeaderboardCacheItemEntity{}
	sql := `SELECT
	        u.id as user_id,
	        u.display_name as username,
	        COALESCE(SUM(
	            ui.quantity *
	            CASE li.rarity
	                WHEN 'common' THEN 1
	                WHEN 'uncommon' THEN 5
	                WHEN 'rare' THEN 25
	                WHEN 'epic' THEN 125
	                WHEN 'legendary' THEN 1000
	                ELSE 0
	            END
	        ), 0)::bigint as score
	        FROM users u
	        LEFT JOIN user_inventory ui ON ui.user_id = u.id
	        LEFT JOIN loot_items li ON li.id = ui.loot_item_id
	        WHERE u.is_archived = false
	        GROUP BY u.id, u.display_name
	        HAVING COALESCE(SUM(ui.quantity), 0) > 0
	        ORDER BY score DESC`

	err := r.db.DB.Select(&items, sql)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetExpeditionCounts calculates total completed expeditions per user
// Returns users sorted by count descending
func (r *LeaderboardRepository) GetExpeditionCounts() ([]*LeaderboardCacheItemEntity, error) {
	items := []*LeaderboardCacheItemEntity{}
	sql := `SELECT
	        u.id as user_id,
	        u.display_name as username,
	        COUNT(e.id)::bigint as score
	        FROM users u
	        LEFT JOIN expeditions e ON e.user_id = u.id AND e.completed = true AND e.is_archived = false
	        WHERE u.is_archived = false
	        GROUP BY u.id, u.display_name
	        HAVING COUNT(e.id) > 0
	        ORDER BY score DESC`

	err := r.db.DB.Select(&items, sql)
	if err != nil {
		return nil, err
	}
	return items, nil
}
