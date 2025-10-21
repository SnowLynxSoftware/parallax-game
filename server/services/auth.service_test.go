package services

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/snowlynxsoftware/parallax-game/config"
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Ensure config.IAppConfig interface is available for NewAuthService
var _ config.IAppConfig = (*MockConfigService)(nil)

// Ensure MockEmailTemplates satisfies IEmailTemplates
var _ IEmailTemplates = (*MockEmailTemplates)(nil)

// MockTokenService is a mock implementation of ITokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) ValidateToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) GenerateVerificationToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) GenerateLoginWithEmailToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockTokenService) ValidateVerificationToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) ValidateLoginWithEmailToken(token *string) (*int, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID int) (*string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

// MockCryptoService is a mock implementation of ICryptoService
type MockCryptoService struct {
	mock.Mock
}

func (m *MockCryptoService) HashPassword(password string) (*string, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockCryptoService) ValidatePassword(password string, hashedPassword string) (bool, error) {
	args := m.Called(password, hashedPassword)
	return args.Bool(0), args.Error(1)
}

// MockEmailService is a mock implementation of IEmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(options *EmailSendOptions) bool {
	args := m.Called(options)
	return args.Bool(0)
}

func (m *MockEmailService) GetTemplates() IEmailTemplates {
	args := m.Called()
	return args.Get(0).(IEmailTemplates)
}

// MockEmailTemplates is a mock implementation of IEmailTemplates
type MockEmailTemplates struct {
	mock.Mock
}

func (m *MockEmailTemplates) GetNewUserEmailTemplate(baseURL string, verificationToken string) string {
	args := m.Called(baseURL, verificationToken)
	return args.String(0)
}

func (m *MockEmailTemplates) GetLoginEmailTemplate(baseURL string, loginToken string) string {
	args := m.Called(baseURL, loginToken)
	return args.String(0)
}

func (m *MockEmailTemplates) GetPasswordResetEmailTemplate(baseURL string, verificationToken string) string {
	args := m.Called(baseURL, verificationToken)
	return args.String(0)
}

// MockConfigService is a mock implementation of IAppConfig
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) GetCloudEnv() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) IsDebugMode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetBaseURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetDBConnectionString() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetAuthHashPepper() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetJWTSecretKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetMJAPIKeyPublic() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetMJAPIKeyPrivate() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetCorsAllowedOrigin() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetCookieDomain() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetSystemAPIKey() string {
	args := m.Called()
	return args.String(0)
}

// Test RegisterNewUser - Success
func TestAuthService_RegisterNewUser_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"
	verificationToken := "verification_token_123"
	emailTemplate := "<html>Verification email</html>"

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       userDTO.Email,
		DisplayName: userDTO.DisplayName,
		CreatedAt:   time.Now(),
	}

	// Mock expectations
	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.MatchedBy(func(dto *models.UserCreateDTO) bool {
		return dto.Email == userDTO.Email && dto.Password == hashedPassword
	})).Return(expectedUser, nil)
	mockTokenService.On("GenerateVerificationToken", 123).Return(&verificationToken, nil)
	mockConfigService.On("GetBaseURL").Return("http://localhost:3000")
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetNewUserEmailTemplate", "http://localhost:3000", verificationToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.MatchedBy(func(options *EmailSendOptions) bool {
		return options.ToEmail == userDTO.Email && options.Subject == "Parallax - Verify Your Account"
	})).Return(true)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test RegisterNewUser - User Already Exists
func TestAuthService_RegisterNewUser_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "existing@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	existingUser := &repositories.UserEntity{
		ID:    123,
		Email: userDTO.Email,
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(existingUser, nil)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "a user already exists with the specified email")
	mockUserRepo.AssertExpectations(t)
}

// Test RegisterNewUser - Password Hashing Error
func TestAuthService_RegisterNewUser_PasswordHashingError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(nil, errors.New("hashing failed"))

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "hashing failed")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test RegisterNewUser - Email Send Failure
func TestAuthService_RegisterNewUser_EmailSendFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"
	verificationToken := "verification_token_123"
	emailTemplate := "<html>Verification email</html>"

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       userDTO.Email,
		DisplayName: userDTO.DisplayName,
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.AnythingOfType("*models.UserCreateDTO")).Return(expectedUser, nil)
	mockTokenService.On("GenerateVerificationToken", 123).Return(&verificationToken, nil)
	mockConfigService.On("GetBaseURL").Return("http://localhost:3000")
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetNewUserEmailTemplate", "http://localhost:3000", verificationToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.AnythingOfType("*services.EmailSendOptions")).Return(false)

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "the user was created but the verification email failed to send")
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test SendLoginEmail - Success
func TestAuthService_SendLoginEmail_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	email := "test@example.com"
	loginToken := "login_token_123"
	emailTemplate := "<html>Login email</html>"

	user := &repositories.UserEntity{
		ID:    123,
		Email: email,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)
	mockTokenService.On("GenerateLoginWithEmailToken", 123).Return(&loginToken, nil)
	mockConfigService.On("GetBaseURL").Return("http://localhost:3000")
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetLoginEmailTemplate", "http://localhost:3000", loginToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.MatchedBy(func(options *EmailSendOptions) bool {
		return options.ToEmail == email && options.Subject == "Parallax - Login Email"
	})).Return(true)

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}

// Test SendLoginEmail - User Banned
// Test SendLoginEmail - User Not Found
func TestAuthService_SendLoginEmail_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	email := "nonexistent@example.com"

	mockUserRepo.On("GetUserByEmail", email).Return(nil, errors.New("user not found"))

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockUserRepo.AssertExpectations(t)
}

// Test LoginWithEmailLink - Success
func TestAuthService_LoginWithEmailLink_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123
	accessToken := "access_token_123"

	mockTokenService.On("GenerateAccessToken", userId).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", &userId).Return(true, nil)

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, accessToken, result.AccessToken)
	assert.Equal(t, "", result.RefreshToken)
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test LoginWithEmailLink - Token Generation Error
func TestAuthService_LoginWithEmailLink_TokenGenerationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123

	mockTokenService.On("GenerateAccessToken", userId).Return(nil, errors.New("token generation failed"))

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockTokenService.AssertExpectations(t)
}

// Test LoginWithEmailLink - Last Login Update Error
func TestAuthService_LoginWithEmailLink_LastLoginUpdateError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123
	accessToken := "access_token_123"

	mockTokenService.On("GenerateAccessToken", userId).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", &userId).Return(false, errors.New("database error"))

	// Act
	result, err := authService.LoginWithEmailLink(&userId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test VerifyNewUser - Success
func TestAuthService_VerifyNewUser_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	verificationToken := "verification_token_123"
	userId := 123

	mockTokenService.On("ValidateToken", &verificationToken).Return(&userId, nil)
	mockUserRepo.On("MarkUserVerified", &userId).Return(true, nil)

	// Act
	result, err := authService.VerifyNewUser(&verificationToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userId, *result)
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test VerifyNewUser - Invalid Token
func TestAuthService_VerifyNewUser_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	verificationToken := "invalid_token"

	mockTokenService.On("ValidateToken", &verificationToken).Return(nil, errors.New("token validation failed"))

	// Act
	result, err := authService.VerifyNewUser(&verificationToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "the token could not be verified")
	mockTokenService.AssertExpectations(t)
}

// Test VerifyNewUser - Database Error
func TestAuthService_VerifyNewUser_DatabaseError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	verificationToken := "verification_token_123"
	userId := 123

	mockTokenService.On("ValidateToken", &verificationToken).Return(&userId, nil)
	mockUserRepo.On("MarkUserVerified", &userId).Return(false, errors.New("database error"))

	// Act
	result, err := authService.VerifyNewUser(&verificationToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockTokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test UpdateUserPassword - Success
func TestAuthService_UpdateUserPassword_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123
	password := "newPassword123"
	hashedPassword := "hashed_new_password"

	mockCryptoService.On("HashPassword", password).Return(&hashedPassword, nil)
	mockUserRepo.On("UpdateUserPassword", &userId, hashedPassword).Return(true, nil)

	// Act
	result, err := authService.UpdateUserPassword(&userId, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userId, *result)
	mockCryptoService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test UpdateUserPassword - Password Hashing Error
func TestAuthService_UpdateUserPassword_HashingError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123
	password := "newPassword123"

	mockCryptoService.On("HashPassword", password).Return(nil, errors.New("hashing failed"))

	// Act
	result, err := authService.UpdateUserPassword(&userId, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "hashing failed")
	mockCryptoService.AssertExpectations(t)
}

// Test UpdateUserPassword - Database Update Error
func TestAuthService_UpdateUserPassword_DatabaseError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userId := 123
	password := "newPassword123"
	hashedPassword := "hashed_new_password"

	mockCryptoService.On("HashPassword", password).Return(&hashedPassword, nil)
	mockUserRepo.On("UpdateUserPassword", &userId, hashedPassword).Return(false, errors.New("database update failed"))

	// Act
	result, err := authService.UpdateUserPassword(&userId, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database update failed")
	mockCryptoService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// Test Login - Success
func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	email := "test@example.com"
	password := "password123"
	hashedPassword := "hashed_password_123"
	accessToken := "access_token_123"

	// Create valid Basic Auth header
	credentials := email + ":" + password
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	authHeader := "Basic " + encodedCredentials

	user := &repositories.UserEntity{
		ID:           123,
		Email:        email,
		PasswordHash: &hashedPassword,
		IsVerified:   true,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)
	mockCryptoService.On("ValidatePassword", password, hashedPassword).Return(true, nil)
	mockTokenService.On("GenerateAccessToken", 123).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", mock.MatchedBy(func(userId *int) bool { return *userId == 123 })).Return(true, nil)

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, accessToken, result.AccessToken)
	assert.Equal(t, "", result.RefreshToken)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test Login - Malformed Authorization Header (No Basic Prefix)
func TestAuthService_Login_MalformedHeaderNoBasic(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	authHeader := "Bearer some-token"

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode authorization header")
}

// Test Login - Invalid Base64 Encoding
func TestAuthService_Login_InvalidBase64(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	authHeader := "Basic invalid-base64!"

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode authorization header")
}

// Test Login - Missing Colon Separator
func TestAuthService_Login_MissingColonSeparator(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "usernamepassword" without colon, base64 encoded
	authHeader := "Basic dXNlcm5hbWVwYXNzd29yZA=="

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid authorization header format")
}

// Test Login - User Not Found
func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "user@example.com:password123" base64 encoded
	authHeader := "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZDEyMw=="

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(nil, errors.New("user not found"))

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockUserRepo.AssertExpectations(t)
}

// Test Login - Invalid Password
func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "user@example.com:wrongpassword" base64 encoded
	authHeader := "Basic dXNlckBleGFtcGxlLmNvbTp3cm9uZ3Bhc3N3b3Jk"
	hashedPassword := "correct_hashed_password"

	user := &repositories.UserEntity{
		ID:           123,
		Email:        "user@example.com",
		PasswordHash: &hashedPassword,
	}

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)
	mockCryptoService.On("ValidatePassword", "wrongpassword", hashedPassword).Return(false, errors.New("password mismatch"))

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test Login - Password Validation Returns False
func TestAuthService_Login_PasswordValidationFalse(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "user@example.com:wrongpassword" base64 encoded
	authHeader := "Basic dXNlckBleGFtcGxlLmNvbTp3cm9uZ3Bhc3N3b3Jk"
	hashedPassword := "correct_hashed_password"

	user := &repositories.UserEntity{
		ID:           123,
		Email:        "user@example.com",
		PasswordHash: &hashedPassword,
	}

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)
	mockCryptoService.On("ValidatePassword", "wrongpassword", hashedPassword).Return(false, nil)

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test Login - Token Generation Error
func TestAuthService_Login_TokenGenerationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "user@example.com:password123" base64 encoded
	authHeader := "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZDEyMw=="
	hashedPassword := "correct_hashed_password"

	user := &repositories.UserEntity{
		ID:           123,
		Email:        "user@example.com",
		PasswordHash: &hashedPassword,
	}

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)
	mockCryptoService.On("ValidatePassword", "password123", hashedPassword).Return(true, nil)
	mockTokenService.On("GenerateAccessToken", 123).Return(nil, errors.New("token generation failed"))

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

// Test Login - Last Login Update Error
func TestAuthService_Login_LastLoginUpdateError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	// "user@example.com:password123" base64 encoded
	authHeader := "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZDEyMw=="
	hashedPassword := "correct_hashed_password"
	accessToken := "access_token_123"

	user := &repositories.UserEntity{
		ID:           123,
		Email:        "user@example.com",
		PasswordHash: &hashedPassword,
	}

	userId := int(user.ID)

	mockUserRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)
	mockCryptoService.On("ValidatePassword", "password123", hashedPassword).Return(true, nil)
	mockTokenService.On("GenerateAccessToken", 123).Return(&accessToken, nil)
	mockUserRepo.On("UpdateUserLastLogin", &userId).Return(false, errors.New("last login update failed"))

	// Act
	result, err := authService.Login(&authHeader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "there was an issue trying to log this user in")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

// Test RegisterNewUser - User Creation Error
func TestAuthService_RegisterNewUser_UserCreationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.MatchedBy(func(dto *models.UserCreateDTO) bool {
		return dto.Email == userDTO.Email && dto.Password == hashedPassword
	})).Return(nil, errors.New("user creation failed"))

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user creation failed")
	mockUserRepo.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test RegisterNewUser - Verification Token Generation Error
func TestAuthService_RegisterNewUser_VerificationTokenError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	userDTO := &models.UserCreateDTO{
		Email:       "test@example.com",
		DisplayName: "Test User",
		Password:    "password123",
	}

	hashedPassword := "hashed_password_123"

	expectedUser := &repositories.UserEntity{
		ID:          123,
		Email:       userDTO.Email,
		DisplayName: userDTO.DisplayName,
		CreatedAt:   time.Now(),
	}

	mockUserRepo.On("GetUserByEmail", userDTO.Email).Return(nil, errors.New("user not found"))
	mockCryptoService.On("HashPassword", userDTO.Password).Return(&hashedPassword, nil)
	mockUserRepo.On("CreateNewUser", mock.MatchedBy(func(dto *models.UserCreateDTO) bool {
		return dto.Email == userDTO.Email && dto.Password == hashedPassword
	})).Return(expectedUser, nil)
	mockTokenService.On("GenerateVerificationToken", 123).Return(nil, errors.New("token generation failed"))

	// Act
	result, err := authService.RegisterNewUser(userDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token generation failed")
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockCryptoService.AssertExpectations(t)
}

// Test SendLoginEmail - Token Generation Error
func TestAuthService_SendLoginEmail_TokenGenerationError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	email := "test@example.com"

	user := &repositories.UserEntity{
		ID:    123,
		Email: email,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)
	mockTokenService.On("GenerateLoginWithEmailToken", 123).Return(nil, errors.New("token generation failed"))

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token generation failed")
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

// Test SendLoginEmail - Email Send Failure
func TestAuthService_SendLoginEmail_EmailSendFailure(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	mockCryptoService := new(MockCryptoService)
	mockEmailService := new(MockEmailService)
	mockEmailTemplates := new(MockEmailTemplates)

	mockConfigService := new(MockConfigService)
	authService := NewAuthService(mockUserRepo, mockTokenService, mockCryptoService, mockEmailService, mockConfigService)

	email := "test@example.com"
	loginToken := "login_token_123"
	emailTemplate := "<html>Login email</html>"

	user := &repositories.UserEntity{
		ID:    123,
		Email: email,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(user, nil)
	mockTokenService.On("GenerateLoginWithEmailToken", 123).Return(&loginToken, nil)
	mockConfigService.On("GetBaseURL").Return("http://localhost:3000")
	mockEmailService.On("GetTemplates").Return(mockEmailTemplates)
	mockEmailTemplates.On("GetLoginEmailTemplate", "http://localhost:3000", loginToken).Return(emailTemplate)
	mockEmailService.On("SendEmail", mock.MatchedBy(func(options *EmailSendOptions) bool {
		return options.ToEmail == email && options.Subject == "Parallax - Login Email"
	})).Return(false)

	// Act
	result, err := authService.SendLoginEmail(email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "the login by email failed to send")
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
	mockEmailTemplates.AssertExpectations(t)
}
