package services

import (
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = 13

type ICryptoService interface {
	HashPassword(password string) (*string, error)
	ValidatePassword(password string, hash string) (bool, error)
}

type CryptoService struct {
	pepper string
}

func NewCryptoService(pepper string) ICryptoService {
	return &CryptoService{
		pepper: pepper,
	}
}

// preparePasswordForBcrypt creates a SHA-256 hash of the password+pepper
// to ensure we stay under bcrypt's 72-byte limit
func (s *CryptoService) preparePasswordForBcrypt(password string) []byte {
	combined := []byte(password + s.pepper)
	hash := sha256.Sum256(combined)
	return hash[:]
}

func (s *CryptoService) HashPassword(password string) (*string, error) {

	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	} else if len(password) < 10 {
		return nil, errors.New("password must be at least 10 characters long")
	}

	var bytesToHash = s.preparePasswordForBcrypt(password)
	var hashBytes, err = bcrypt.GenerateFromPassword(bytesToHash, bcryptCost)
	if err != nil {
		return nil, err
	}
	var hashedString = string(hashBytes)

	return &hashedString, nil
}

func (s *CryptoService) ValidatePassword(password string, hash string) (bool, error) {

	var bytesToHash = s.preparePasswordForBcrypt(password)
	var err = bcrypt.CompareHashAndPassword([]byte(hash), bytesToHash)
	if err != nil {
		return false, err
	}
	return true, nil
}
