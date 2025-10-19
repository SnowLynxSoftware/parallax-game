package services

import (
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
)

// Mock repository for testing
type mockFeatureFlagRepository struct {
	flags      map[string]*repositories.FeatureFlagEntity
	shouldFail bool
}

func newMockFeatureFlagRepository() *mockFeatureFlagRepository {
	now := time.Now()
	return &mockFeatureFlagRepository{
		flags: map[string]*repositories.FeatureFlagEntity{
			"test_enabled_flag": {
				ID:          1,
				CreatedAt:   now,
				ModifiedAt:  nil,
				IsArchived:  false,
				Key:         "test_enabled_flag",
				Enabled:     true,
				Description: "Test flag that is enabled",
			},
			"test_disabled_flag": {
				ID:          2,
				CreatedAt:   now,
				ModifiedAt:  nil,
				IsArchived:  false,
				Key:         "test_disabled_flag",
				Enabled:     false,
				Description: "Test flag that is disabled",
			},
			"test_archived_flag": {
				ID:          3,
				CreatedAt:   now,
				ModifiedAt:  &now,
				IsArchived:  true,
				Key:         "test_archived_flag",
				Enabled:     true,
				Description: "Test flag that is archived",
			},
		},
		shouldFail: false,
	}
}

func (m *mockFeatureFlagRepository) GetAllFlags() ([]*repositories.FeatureFlagEntity, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}

	var flags []*repositories.FeatureFlagEntity
	for _, flag := range m.flags {
		if !flag.IsArchived {
			flags = append(flags, flag)
		}
	}
	return flags, nil
}

func (m *mockFeatureFlagRepository) GetFlagByKey(key string) (*repositories.FeatureFlagEntity, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}

	flag, exists := m.flags[key]
	if !exists {
		return nil, nil
	}

	// Return nil if archived (simulating real repository behavior)
	if flag.IsArchived {
		return nil, nil
	}

	return flag, nil
}

func (m *mockFeatureFlagRepository) CreateFlag(key string, enabled bool, description string) (*repositories.FeatureFlagEntity, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}

	now := time.Now()
	flag := &repositories.FeatureFlagEntity{
		ID:          int64(len(m.flags) + 1),
		CreatedAt:   now,
		ModifiedAt:  nil,
		IsArchived:  false,
		Key:         key,
		Enabled:     enabled,
		Description: description,
	}
	m.flags[key] = flag
	return flag, nil
}

func (m *mockFeatureFlagRepository) UpdateFlag(key string, enabled bool, description string) (*repositories.FeatureFlagEntity, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}

	flag, exists := m.flags[key]
	if !exists || flag.IsArchived {
		return nil, nil
	}

	now := time.Now()
	flag.Enabled = enabled
	flag.Description = description
	flag.ModifiedAt = &now
	return flag, nil
}

func (m *mockFeatureFlagRepository) ArchiveFlag(key string) error {
	if m.shouldFail {
		return errors.New("database error")
	}

	flag, exists := m.flags[key]
	if exists {
		now := time.Now()
		flag.IsArchived = true
		flag.ModifiedAt = &now
	}
	return nil
}

// Tests for IsEnabled method

func TestIsEnabled_FlagExistsAndEnabled(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	result := service.IsEnabled("test_enabled_flag")

	if !result {
		t.Errorf("Expected IsEnabled to return true for enabled flag, got false")
	}
}

func TestIsEnabled_FlagExistsAndDisabled(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	result := service.IsEnabled("test_disabled_flag")

	if result {
		t.Errorf("Expected IsEnabled to return false for disabled flag, got true")
	}
}

func TestIsEnabled_FlagDoesNotExist(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	result := service.IsEnabled("non_existent_flag")

	if result {
		t.Errorf("Expected IsEnabled to return false for non-existent flag, got true")
	}
}

func TestIsEnabled_FlagArchived(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	result := service.IsEnabled("test_archived_flag")

	if result {
		t.Errorf("Expected IsEnabled to return false for archived flag, got true")
	}
}

func TestIsEnabled_DatabaseError(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	mockRepo.shouldFail = true
	service := NewFeatureFlagService(mockRepo)

	result := service.IsEnabled("test_enabled_flag")

	if result {
		t.Errorf("Expected IsEnabled to return false on database error, got true")
	}
}

// Tests for GetFlag method

func TestGetFlag_FlagExists(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	flag, err := service.GetFlag("test_enabled_flag")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if flag == nil {
		t.Fatal("Expected flag to be returned, got nil")
	}
	if flag.Key != "test_enabled_flag" {
		t.Errorf("Expected flag key 'test_enabled_flag', got '%s'", flag.Key)
	}
	if !flag.Enabled {
		t.Errorf("Expected flag to be enabled, got disabled")
	}
}

func TestGetFlag_FlagDoesNotExist(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	flag, err := service.GetFlag("non_existent_flag")

	if err != nil {
		t.Errorf("Expected no error for non-existent flag, got %v", err)
	}
	if flag != nil {
		t.Errorf("Expected nil flag for non-existent flag, got %v", flag)
	}
}

func TestGetFlag_DatabaseError(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	mockRepo.shouldFail = true
	service := NewFeatureFlagService(mockRepo)

	flag, err := service.GetFlag("test_enabled_flag")

	if err == nil {
		t.Error("Expected error on database failure, got nil")
	}
	if flag != nil {
		t.Errorf("Expected nil flag on error, got %v", flag)
	}
}

// Tests for GetAllFlags method

func TestGetAllFlags_ReturnsMap(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	flags, err := service.GetAllFlags()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if flags == nil {
		t.Fatal("Expected flags map to be returned, got nil")
	}

	// Should return 2 non-archived flags
	if len(flags) != 2 {
		t.Errorf("Expected 2 non-archived flags, got %d", len(flags))
	}

	// Check that archived flag is not included
	if _, exists := flags["test_archived_flag"]; exists {
		t.Error("Expected archived flag to be excluded from results")
	}

	// Check that active flags are included
	if _, exists := flags["test_enabled_flag"]; !exists {
		t.Error("Expected test_enabled_flag to be in results")
	}
	if _, exists := flags["test_disabled_flag"]; !exists {
		t.Error("Expected test_disabled_flag to be in results")
	}
}

func TestGetAllFlags_DatabaseError(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	mockRepo.shouldFail = true
	service := NewFeatureFlagService(mockRepo)

	flags, err := service.GetAllFlags()

	if err == nil {
		t.Error("Expected error on database failure, got nil")
	}
	if flags != nil {
		t.Errorf("Expected nil flags on error, got %v", flags)
	}
}

// Tests for GetAllFlagsList method

func TestGetAllFlagsList_ReturnsList(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	service := NewFeatureFlagService(mockRepo)

	flags, err := service.GetAllFlagsList()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if flags == nil {
		t.Fatal("Expected flags list to be returned, got nil")
	}

	// Should return 2 non-archived flags
	if len(flags) != 2 {
		t.Errorf("Expected 2 non-archived flags, got %d", len(flags))
	}

	// Verify no archived flags in list
	for _, flag := range flags {
		if flag.IsArchived {
			t.Errorf("Expected no archived flags in list, found %s", flag.Key)
		}
	}
}

func TestGetAllFlagsList_DatabaseError(t *testing.T) {
	mockRepo := newMockFeatureFlagRepository()
	mockRepo.shouldFail = true
	service := NewFeatureFlagService(mockRepo)

	flags, err := service.GetAllFlagsList()

	if err == nil {
		t.Error("Expected error on database failure, got nil")
	}
	if flags != nil {
		t.Errorf("Expected nil flags on error, got %v", flags)
	}
}
