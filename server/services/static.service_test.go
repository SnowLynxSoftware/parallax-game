package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: StaticService is out of scope for unit testing due to embedded file system dependency
// The service would require complex mocking of the embedded file system which is beyond
// the scope of pure unit tests

func TestStaticService_Creation(t *testing.T) {
	// Arrange & Act
	staticService := NewStaticService()

	// Assert
	assert.NotNil(t, staticService)
}
