package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExpirationInMinutes         = 60
	verificationTokenExpirationInHours     = 3
	loginWithEmailTokenExpirationInMinutes = 10
	refreshTokenExpirationInHours          = 160
	claimIssuer                            = "https://parallax.com"
)

type ITokenService interface {
	GenerateAccessToken(id int) (*string, error)
	GenerateLoginWithEmailToken(id int) (*string, error)
	GenerateVerificationToken(id int) (*string, error)
	GenerateRefreshToken(id int) (*string, error)
	ValidateToken(tokenToVerify *string) (*int, error)
}

type TokenService struct {
	jwtSecretKey string
}

func NewTokenService(jwtSecretKey string) ITokenService {
	return &TokenService{
		jwtSecretKey: jwtSecretKey,
	}
}

func (s *TokenService) GenerateAccessToken(id int) (*string, error) {

	expirationTime := time.Now().Add(accessTokenExpirationInMinutes * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":  claimIssuer,
			"sub":  "api_access_token",
			"exp":  expirationTime,
			"user": id,
		})
	signedToken, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (s *TokenService) GenerateLoginWithEmailToken(id int) (*string, error) {

	expirationTime := time.Now().Add(loginWithEmailTokenExpirationInMinutes * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":  claimIssuer,
			"sub":  "loginwithemail_token",
			"exp":  expirationTime,
			"user": id,
		})
	signedToken, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (s *TokenService) GenerateVerificationToken(id int) (*string, error) {

	expirationTime := time.Now().Add(verificationTokenExpirationInHours * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":  claimIssuer,
			"sub":  "api_verification_token",
			"exp":  expirationTime,
			"user": id,
		})
	signedToken, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (s *TokenService) GenerateRefreshToken(id int) (*string, error) {

	expirationTime := time.Now().Add(refreshTokenExpirationInHours * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":  claimIssuer,
			"sub":  "api_refresh_token",
			"exp":  expirationTime,
			"user": id,
		})
	signedToken, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (s *TokenService) ValidateToken(tokenToVerify *string) (*int, error) {

	parsedToken, err := jwt.Parse(*tokenToVerify, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.jwtSecretKey), nil
	})
	if err != nil {
		return nil, errors.New("JWT could not be validated")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		userClaim, exists := claims["user"]
		if !exists {
			return nil, errors.New("JWT claims could not be validated")
		}

		userFloat, ok := userClaim.(float64)
		if !ok {
			return nil, errors.New("JWT claims could not be validated")
		}

		userId := int(userFloat)
		return &userId, nil
	} else {
		return nil, errors.New("JWT claims could not be validated")
	}
}
