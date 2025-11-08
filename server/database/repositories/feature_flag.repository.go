package repositories

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/snowlynxsoftware/parallax-game/server/database"
)

type FeatureFlagEntity struct {
	ID          int64      `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at" db:"modified_at"`
	IsArchived  bool       `json:"is_archived" db:"is_archived"`
	Key         string     `json:"key" db:"key"`
	Enabled     bool       `json:"enabled" db:"enabled"`
	Description string     `json:"description" db:"description"`
}

type IFeatureFlagRepository interface {
	// GetAllFlags retrieves all non-archived feature flags
	GetAllFlags() ([]*FeatureFlagEntity, error)

	// GetFlagByKey retrieves a single feature flag by its key
	GetFlagByKey(key string) (*FeatureFlagEntity, error)

	// CreateFlag inserts a new feature flag (for future admin functionality)
	CreateFlag(key string, enabled bool, description string) (*FeatureFlagEntity, error)

	// UpdateFlag updates an existing feature flag (for future admin functionality)
	UpdateFlag(key string, enabled bool, description string) (*FeatureFlagEntity, error)

	// ArchiveFlag archives a feature flag (for future admin functionality)
	ArchiveFlag(key string) error
}

type FeatureFlagRepository struct {
	db *database.AppDataSource
}

func NewFeatureFlagRepository(db *database.AppDataSource) IFeatureFlagRepository {
	return &FeatureFlagRepository{
		db: db,
	}
}

// GetAllFlags retrieves all non-archived feature flags ordered by key
func (r *FeatureFlagRepository) GetAllFlags() ([]*FeatureFlagEntity, error) {
	var flags []*FeatureFlagEntity
	query := `SELECT * FROM feature_flags WHERE is_archived = false ORDER BY key ASC`
	err := r.db.DB.Select(&flags, query)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error retrieving all feature flags from database")
		return nil, err
	}
	return flags, nil
}

// GetFlagByKey retrieves a single feature flag by its key
// Returns nil, nil if flag doesn't exist (not an error condition)
func (r *FeatureFlagRepository) GetFlagByKey(key string) (*FeatureFlagEntity, error) {
	flag := &FeatureFlagEntity{}
	query := `SELECT * FROM feature_flags WHERE key = $1 AND is_archived = false`
	err := r.db.DB.Get(flag, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Flag not found - return nil, not error
		}
		log.Error().
			Err(err).
			Str("flag_key", key).
			Msg("Error retrieving feature flag from database")
		return nil, err
	}
	return flag, nil
}

// CreateFlag inserts a new feature flag
func (r *FeatureFlagRepository) CreateFlag(key string, enabled bool, description string) (*FeatureFlagEntity, error) {
	query := `INSERT INTO feature_flags (key, enabled, description)
			  VALUES ($1, $2, $3)
			  RETURNING *`

	flag := &FeatureFlagEntity{}
	err := r.db.DB.Get(flag, query, key, enabled, description)
	if err != nil {
		log.Error().
			Err(err).
			Str("flag_key", key).
			Bool("enabled", enabled).
			Msg("Failed to create feature flag")
		return nil, err
	}

	log.Info().
		Str("flag_key", key).
		Bool("enabled", enabled).
		Str("description", description).
		Msg("Feature flag created")

	return flag, nil
}

// UpdateFlag updates an existing feature flag
func (r *FeatureFlagRepository) UpdateFlag(key string, enabled bool, description string) (*FeatureFlagEntity, error) {
	query := `UPDATE feature_flags
			  SET enabled = $1, description = $2, modified_at = NOW()
			  WHERE key = $3 AND is_archived = false
			  RETURNING *`

	flag := &FeatureFlagEntity{}
	err := r.db.DB.Get(flag, query, enabled, description, key)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().
				Str("flag_key", key).
				Msg("Attempted to update non-existent feature flag")
			return nil, nil
		}
		log.Error().
			Err(err).
			Str("flag_key", key).
			Bool("enabled", enabled).
			Msg("Failed to update feature flag")
		return nil, err
	}

	log.Info().
		Str("flag_key", key).
		Bool("enabled", enabled).
		Str("description", description).
		Msg("Feature flag updated")

	return flag, nil
}

// ArchiveFlag archives a feature flag (soft delete)
func (r *FeatureFlagRepository) ArchiveFlag(key string) error {
	query := `UPDATE feature_flags
			  SET is_archived = true, modified_at = NOW()
			  WHERE key = $1`

	result, err := r.db.DB.Exec(query, key)
	if err != nil {
		log.Error().
			Err(err).
			Str("flag_key", key).
			Msg("Failed to archive feature flag")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		log.Warn().
			Str("flag_key", key).
			Msg("Attempted to archive non-existent feature flag")
	} else {
		log.Info().
			Str("flag_key", key).
			Msg("Feature flag archived")
	}

	return nil
}
