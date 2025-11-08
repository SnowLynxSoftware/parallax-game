package services

import (
	"github.com/rs/zerolog/log"
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
)

type IFeatureFlagService interface {
	// IsEnabled checks if a feature flag is enabled (primary usage method)
	// Returns false if flag doesn't exist or is disabled/archived
	IsEnabled(key string) bool

	// GetFlag retrieves a single feature flag with full metadata
	// Returns nil if flag doesn't exist
	GetFlag(key string) (*repositories.FeatureFlagEntity, error)

	// GetAllFlags retrieves all non-archived feature flags as a map
	// Map key is the flag key, value is the entity
	GetAllFlags() (map[string]*repositories.FeatureFlagEntity, error)

	// GetAllFlagsList retrieves all non-archived feature flags as a slice
	// Useful for admin interfaces or logging
	GetAllFlagsList() ([]*repositories.FeatureFlagEntity, error)
}

type FeatureFlagService struct {
	featureFlagRepository repositories.IFeatureFlagRepository
}

func NewFeatureFlagService(featureFlagRepository repositories.IFeatureFlagRepository) IFeatureFlagService {
	return &FeatureFlagService{
		featureFlagRepository: featureFlagRepository,
	}
}

// IsEnabled checks if a feature flag is enabled
// Returns false if flag doesn't exist, is disabled, or there's an error
// This method never returns an error to provide a clean, simple API
func (s *FeatureFlagService) IsEnabled(key string) bool {
	flag, err := s.featureFlagRepository.GetFlagByKey(key)
	if err != nil {
		// Log error but return false (safe default)
		log.Error().
			Err(err).
			Str("flag_key", key).
			Msg("Error checking feature flag")
		return false
	}
	if flag == nil {
		// Flag doesn't exist - log at debug level to catch typos
		log.Debug().
			Str("flag_key", key).
			Msg("Feature flag does not exist, returning false")
		return false
	}

	// Log flag check at debug level for visibility
	log.Debug().
		Str("flag_key", key).
		Bool("enabled", flag.Enabled).
		Msg("Feature flag checked")

	return flag.Enabled
}

// GetFlag retrieves a single feature flag with full metadata
func (s *FeatureFlagService) GetFlag(key string) (*repositories.FeatureFlagEntity, error) {
	flag, err := s.featureFlagRepository.GetFlagByKey(key)
	if err != nil {
		log.Error().
			Err(err).
			Str("flag_key", key).
			Msg("Error retrieving feature flag")
		return nil, err
	}
	return flag, nil
}

// GetAllFlags retrieves all non-archived feature flags as a map
// Map key is the flag key, value is the entity
func (s *FeatureFlagService) GetAllFlags() (map[string]*repositories.FeatureFlagEntity, error) {
	flags, err := s.featureFlagRepository.GetAllFlags()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error retrieving all feature flags")
		return nil, err
	}

	flagMap := make(map[string]*repositories.FeatureFlagEntity)
	for _, flag := range flags {
		flagMap[flag.Key] = flag
	}

	log.Debug().
		Int("flag_count", len(flagMap)).
		Msg("Retrieved all feature flags")

	return flagMap, nil
}

// GetAllFlagsList retrieves all non-archived feature flags as a slice
func (s *FeatureFlagService) GetAllFlagsList() ([]*repositories.FeatureFlagEntity, error) {
	flags, err := s.featureFlagRepository.GetAllFlags()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error retrieving feature flags list")
		return nil, err
	}

	log.Debug().
		Int("flag_count", len(flags)).
		Msg("Retrieved feature flags list")

	return flags, nil
}
