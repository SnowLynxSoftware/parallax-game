package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetNewUserEmailTemplate - Success
func TestEmailTemplates_GetNewUserEmailTemplate_Success(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	verificationToken := "test-verification-token-123"

	// Act
	result := emailTemplates.GetNewUserEmailTemplate(baseURL, verificationToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, verificationToken)
	assert.Contains(t, result, "/api/auth/verify?token=")
	assert.Contains(t, result, "Welcome to Parallax")
	assert.Contains(t, result, "Verify Account")
	assert.Contains(t, result, "Parallax")
}

// Test GetNewUserEmailTemplate - With Empty BaseURL
func TestEmailTemplates_GetNewUserEmailTemplate_EmptyBaseURL(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := ""
	verificationToken := "test-verification-token-123"

	// Act
	result := emailTemplates.GetNewUserEmailTemplate(baseURL, verificationToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, verificationToken)
	assert.Contains(t, result, "/api/auth/verify?token=")
	// Should still contain HTML structure even with empty baseURL
	assert.Contains(t, result, "<a href=")
}

// Test GetNewUserEmailTemplate - With Empty Token
func TestEmailTemplates_GetNewUserEmailTemplate_EmptyToken(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	verificationToken := ""

	// Act
	result := emailTemplates.GetNewUserEmailTemplate(baseURL, verificationToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, "/api/auth/verify?token=")
	// Should contain the button even with empty token
	assert.Contains(t, result, "Verify Account")
}

// Test GetNewUserEmailTemplate - With Special Characters
func TestEmailTemplates_GetNewUserEmailTemplate_SpecialCharacters(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com/path?param=value&other=test"
	verificationToken := "token-with-special-chars-!@#$%"

	// Act
	result := emailTemplates.GetNewUserEmailTemplate(baseURL, verificationToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, verificationToken)
	assert.Contains(t, result, "/api/auth/verify?token=")
}

// Test GetLoginEmailTemplate - Success
func TestEmailTemplates_GetLoginEmailTemplate_Success(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	loginToken := "test-login-token-456"

	// Act
	result := emailTemplates.GetLoginEmailTemplate(baseURL, loginToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, loginToken)
	assert.Contains(t, result, "/api/auth/login-with-email?token=")
	assert.Contains(t, result, "Log in to your account")
	assert.Contains(t, result, "Log In Instantly")
	assert.Contains(t, result, "If you did not request this email")
}

// Test GetLoginEmailTemplate - With Empty BaseURL
func TestEmailTemplates_GetLoginEmailTemplate_EmptyBaseURL(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := ""
	loginToken := "test-login-token-456"

	// Act
	result := emailTemplates.GetLoginEmailTemplate(baseURL, loginToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, loginToken)
	assert.Contains(t, result, "/api/auth/login-with-email?token=")
	// Should still contain HTML structure even with empty baseURL
	assert.Contains(t, result, "<a href=")
}

// Test GetLoginEmailTemplate - With Empty Token
func TestEmailTemplates_GetLoginEmailTemplate_EmptyToken(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	loginToken := ""

	// Act
	result := emailTemplates.GetLoginEmailTemplate(baseURL, loginToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, "/api/auth/login-with-email?token=")
	// Should contain the button even with empty token
	assert.Contains(t, result, "Log In Instantly")
}

// Test GetLoginEmailTemplate - With Special Characters
func TestEmailTemplates_GetLoginEmailTemplate_SpecialCharacters(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com/login?redirect=/dashboard"
	loginToken := "login-token-with-symbols-&%*"

	// Act
	result := emailTemplates.GetLoginEmailTemplate(baseURL, loginToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, loginToken)
	assert.Contains(t, result, "/api/auth/login-with-email?token=")
}

// Test GetPasswordResetEmailTemplate - Success
func TestEmailTemplates_GetPasswordResetEmailTemplate_Success(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	resetToken := "test-reset-token-789"

	// Act
	result := emailTemplates.GetPasswordResetEmailTemplate(baseURL, resetToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, resetToken)
	assert.Contains(t, result, "/reset-password?token=")
	assert.Contains(t, result, "Reset your password")
	assert.Contains(t, result, "Reset Password")
}

// Test GetPasswordResetEmailTemplate - With Empty BaseURL
func TestEmailTemplates_GetPasswordResetEmailTemplate_EmptyBaseURL(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := ""
	resetToken := "test-reset-token-789"

	// Act
	result := emailTemplates.GetPasswordResetEmailTemplate(baseURL, resetToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, resetToken)
	assert.Contains(t, result, "/reset-password?token=")
	assert.Contains(t, result, "<a href=")
}

// Test GetPasswordResetEmailTemplate - With Empty Token
func TestEmailTemplates_GetPasswordResetEmailTemplate_EmptyToken(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	resetToken := ""

	// Act
	result := emailTemplates.GetPasswordResetEmailTemplate(baseURL, resetToken)

	// Assert
	assert.NotEmpty(t, result)
	assert.Contains(t, result, baseURL)
	assert.Contains(t, result, "/reset-password?token=")
	assert.Contains(t, result, "Reset Password")
}

// Test Template Structure - Verify HTML Structure
func TestEmailTemplates_VerifyHTMLStructure(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	token := "test-token"

	// Act
	newUserTemplate := emailTemplates.GetNewUserEmailTemplate(baseURL, token)
	loginTemplate := emailTemplates.GetLoginEmailTemplate(baseURL, token)
	resetTemplate := emailTemplates.GetPasswordResetEmailTemplate(baseURL, token)

	// Assert - Check for enhanced HTML structure
	for _, template := range []string{newUserTemplate, loginTemplate, resetTemplate} {
		assert.Contains(t, template, "<!DOCTYPE html>")
		assert.Contains(t, template, "<body")
		assert.Contains(t, template, "<table")
		assert.Contains(t, template, "Parallax")
		assert.Contains(t, template, "Facebook")
		assert.Contains(t, template, "Website")
	}
}

// Test Template Differences
func TestEmailTemplates_TemplateDifferences(t *testing.T) {
	// Arrange
	emailTemplates := NewEmailTemplates()
	baseURL := "https://example.com"
	token := "test-token"

	// Act
	newUserTemplate := emailTemplates.GetNewUserEmailTemplate(baseURL, token)
	loginTemplate := emailTemplates.GetLoginEmailTemplate(baseURL, token)
	resetTemplate := emailTemplates.GetPasswordResetEmailTemplate(baseURL, token)

	// Assert - Templates should be different
	assert.NotEqual(t, newUserTemplate, loginTemplate)
	assert.NotEqual(t, newUserTemplate, resetTemplate)
	assert.NotEqual(t, loginTemplate, resetTemplate)

	// Verify unique content in each template
	assert.Contains(t, newUserTemplate, "Welcome to Parallax")
	assert.Contains(t, newUserTemplate, "Verify Account")
	assert.Contains(t, newUserTemplate, "/api/auth/verify")

	assert.Contains(t, loginTemplate, "Log in to your account")
	assert.Contains(t, loginTemplate, "Log In Instantly")
	assert.Contains(t, loginTemplate, "/api/auth/login-with-email")

	assert.Contains(t, resetTemplate, "Reset your password")
	assert.Contains(t, resetTemplate, "Reset Password")
	assert.Contains(t, resetTemplate, "/reset-password")
}
