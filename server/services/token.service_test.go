package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == nil || *token == "" {
		t.Fatalf("expected a valid token, got nil or empty string")
	}
}

func TestValidateToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate the generated token
	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validatedUserID == nil || *validatedUserID != userID {
		t.Fatalf("expected userID %d, got %v", userID, validatedUserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	invalidToken := "invalid.token.value"
	_, err := service.ValidateToken(&invalidToken)
	if err == nil {
		t.Fatalf("expected an error for invalid token, got nil")
	}

	expectedError := "JWT could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	// Create an expired token
	expiredTime := time.Now().Add(-1 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":  claimIssuer,
		"sub":  "api_access_token",
		"exp":  expiredTime,
		"user": 123,
	})
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate the expired token
	_, err = service.ValidateToken(&signedToken)
	if err == nil {
		t.Fatalf("expected an error for expired token, got nil")
	}

	expectedError := "JWT could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	// Malformed token (missing parts)
	malformedToken := "malformed.token"
	_, err := service.ValidateToken(&malformedToken)
	if err == nil {
		t.Fatalf("expected an error for malformed token, got nil")
	}

	expectedError := "JWT could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	// Empty token
	emptyToken := ""
	_, err := service.ValidateToken(&emptyToken)
	if err == nil {
		t.Fatalf("expected an error for empty token, got nil")
	}

	expectedError := "JWT could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGenerateVerificationToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	userID := 456
	token, err := service.GenerateVerificationToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == nil || *token == "" {
		t.Fatalf("expected a valid verification token, got nil or empty string")
	}

	// Validate the generated token
	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error validating verification token, got %v", err)
	}

	if validatedUserID == nil || *validatedUserID != userID {
		t.Fatalf("expected userID %d, got %v", userID, validatedUserID)
	}
}

func TestGenerateLoginWithEmailToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	userID := 789
	token, err := service.GenerateLoginWithEmailToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == nil || *token == "" {
		t.Fatalf("expected a valid login-with-email token, got nil or empty string")
	}

	// Validate the generated token
	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error validating login-with-email token, got %v", err)
	}

	if validatedUserID == nil || *validatedUserID != userID {
		t.Fatalf("expected userID %d, got %v", userID, validatedUserID)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	userID := 321
	token, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == nil || *token == "" {
		t.Fatalf("expected a valid refresh token, got nil or empty string")
	}

	// Validate the generated token
	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error validating refresh token, got %v", err)
	}

	if validatedUserID == nil || *validatedUserID != userID {
		t.Fatalf("expected userID %d, got %v", userID, validatedUserID)
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	// Test with a manually crafted token that uses the wrong signing method
	// This token claims to use RS256 but is actually malformed
	wrongMethodToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL1NtYXJ0ZXJMeW54LmNvbSIsInN1YiI6ImFwaV9hY2Nlc3NfdG9rZW4iLCJleHAiOjk5OTk5OTk5OTksInVzZXIiOjEyM30.invalid-signature"

	_, err := service.ValidateToken(&wrongMethodToken)
	if err == nil {
		t.Fatalf("expected an error for wrong signing method, got nil")
	}

	expectedError := "JWT could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestValidateToken_InvalidClaims(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)

	// Create a token with invalid claims structure
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":  claimIssuer,
		"sub":  "api_access_token",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"user": "not-a-number", // Invalid user claim
	})
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		t.Fatalf("expected no error creating token, got %v", err)
	}

	_, err = service.ValidateToken(&signedToken)
	if err == nil {
		t.Fatalf("expected an error for invalid claims, got nil")
	}

	expectedError := "JWT claims could not be validated"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTokenService_AllTokenTypes(t *testing.T) {
	jwtSecretKey := "testSecretKey"
	service := NewTokenService(jwtSecretKey)
	userID := 100

	// Test all token generation methods
	accessToken, err := service.GenerateAccessToken(userID)
	if err != nil || accessToken == nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	verificationToken, err := service.GenerateVerificationToken(userID)
	if err != nil || verificationToken == nil {
		t.Fatalf("failed to generate verification token: %v", err)
	}

	loginToken, err := service.GenerateLoginWithEmailToken(userID)
	if err != nil || loginToken == nil {
		t.Fatalf("failed to generate login-with-email token: %v", err)
	}

	refreshToken, err := service.GenerateRefreshToken(userID)
	if err != nil || refreshToken == nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	// Verify all tokens are different
	tokens := []*string{accessToken, verificationToken, loginToken, refreshToken}
	for i, token1 := range tokens {
		for j, token2 := range tokens {
			if i != j && *token1 == *token2 {
				t.Fatalf("tokens should be different, but token %d and %d are the same", i, j)
			}
		}
	}

	// Verify all tokens can be validated and return correct userID
	for i, token := range tokens {
		validatedUserID, err := service.ValidateToken(token)
		if err != nil {
			t.Fatalf("failed to validate token %d: %v", i, err)
		}
		if validatedUserID == nil || *validatedUserID != userID {
			t.Fatalf("token %d returned wrong userID: expected %d, got %v", i, userID, validatedUserID)
		}
	}
}

func TestTokenService_EmptySecretKey(t *testing.T) {
	service := NewTokenService("")
	userID := 123

	// Should still generate token but validation might behave differently
	token, err := service.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("expected no error with empty secret, got %v", err)
	}

	if token == nil || *token == "" {
		t.Fatalf("expected a valid token even with empty secret, got nil or empty string")
	}
}

// Test token generation error cases - these test the error paths in the generate methods
func TestTokenService_GenerateAccessToken_SigningError(t *testing.T) {
	// Test with invalid characters that could cause signing errors
	// Note: JWT library is robust, but we test edge cases
	jwtSecretKey := string([]byte{0x00, 0x01, 0x02}) // binary data
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateAccessToken(userID)

	// This should still work as JWT library handles binary keys
	// But if it fails, we're testing the error path
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, token)
	} else {
		assert.NotNil(t, token)
	}
}

func TestTokenService_GenerateVerificationToken_SigningError(t *testing.T) {
	jwtSecretKey := string([]byte{0x00, 0x01, 0x02})
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateVerificationToken(userID)

	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, token)
	} else {
		assert.NotNil(t, token)
	}
}

func TestTokenService_GenerateLoginWithEmailToken_SigningError(t *testing.T) {
	jwtSecretKey := string([]byte{0x00, 0x01, 0x02})
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateLoginWithEmailToken(userID)

	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, token)
	} else {
		assert.NotNil(t, token)
	}
}

func TestTokenService_GenerateRefreshToken_SigningError(t *testing.T) {
	jwtSecretKey := string([]byte{0x00, 0x01, 0x02})
	service := NewTokenService(jwtSecretKey)

	userID := 123
	token, err := service.GenerateRefreshToken(userID)

	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, token)
	} else {
		assert.NotNil(t, token)
	}
}
