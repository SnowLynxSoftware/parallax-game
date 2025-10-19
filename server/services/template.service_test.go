package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: TemplateService is out of scope for unit testing due to embedded file system dependency
// The service would require complex mocking of the embedded template files which is beyond
// the scope of pure unit tests

func TestTemplateService_Creation(t *testing.T) {
	// Arrange & Act
	templateService := NewTemplateService()

	// Assert
	assert.NotNil(t, templateService)
}
