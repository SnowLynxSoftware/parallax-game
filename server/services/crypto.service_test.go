package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	password := "securePassword123"
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashedPassword == nil || *hashedPassword == "" {
		t.Fatalf("expected a valid hashed password, got nil or empty string")
	}

	// Ensure the hash is not the same as the input password
	if *hashedPassword == password+pepper {
		t.Fatalf("hashed password should not match the raw password + pepper")
	}
}

func TestValidatePassword(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	password := "securePassword123"
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test valid password
	isValid, err := service.ValidatePassword(password, *hashedPassword)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !isValid {
		t.Fatalf("expected password to be valid, got invalid")
	}

	// Test invalid password
	isValid, err = service.ValidatePassword("wrongPassword", *hashedPassword)
	if err == nil || err != bcrypt.ErrMismatchedHashAndPassword {
		t.Fatalf("expected bcrypt.ErrMismatchedHashAndPassword, got %v", err)
	}
	if isValid {
		t.Fatalf("expected password to be invalid, got valid")
	}
}

func TestEmptyPassword(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	// Test empty password
	_, err := service.HashPassword("")
	if err == nil {
		t.Fatalf("expected an error for empty password, got nil")
	}

	expectedError := "password cannot be empty"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestShortPassword(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	// Test password shorter than 10 characters
	shortPassword := "short"
	_, err := service.HashPassword(shortPassword)
	if err == nil {
		t.Fatalf("expected an error for short password, got nil")
	}

	expectedError := "password must be at least 10 characters long"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestHashPassword_ExactlyTenCharacters(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	// Test password with exactly 10 characters (minimum valid length)
	password := "1234567890"
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error for 10 character password, got %v", err)
	}

	if hashedPassword == nil || *hashedPassword == "" {
		t.Fatalf("expected a valid hashed password, got nil or empty string")
	}
}

func TestHashPassword_NineCharacters(t *testing.T) {
	pepper := "testPepper"
	service := NewCryptoService(pepper)

	// Test password with 9 characters (one below minimum)
	password := "123456789"
	_, err := service.HashPassword(password)
	if err == nil {
		t.Fatalf("expected an error for 9 character password, got nil")
	}

	expectedError := "password must be at least 10 characters long"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}
