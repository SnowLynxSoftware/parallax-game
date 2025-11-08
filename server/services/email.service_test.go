package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test EmailService Creation and GetTemplates
func TestEmailService_NewEmailService_And_GetTemplates(t *testing.T) {
	// Arrange
	apiKeyPublic := "test-api-key-12345"
	apiKeyPrivate := "test-api-key-12345"
	mockTemplates := NewEmailTemplates()

	// Act
	emailService := NewEmailService(apiKeyPublic, apiKeyPrivate, mockTemplates)

	// Assert
	assert.NotNil(t, emailService)

	// Test GetTemplates
	result := emailService.GetTemplates()
	assert.Equal(t, mockTemplates, result)
}
