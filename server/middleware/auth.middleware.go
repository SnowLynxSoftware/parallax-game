package middleware

import (
	"errors"
	"net/http"

	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type IAuthMiddleware interface {
	Authorize(r *http.Request) (*AuthorizedUserContext, error)
	ValidateSystemAPIKey(r *http.Request) error
}

type AuthMiddleware struct {
	userRepository repositories.IUserRepository
	tokenService   services.ITokenService
	systemAPIKey   string
}

func NewAuthMiddleware(userRepository repositories.IUserRepository, tokenService services.ITokenService, systemAPIKey string) IAuthMiddleware {
	return &AuthMiddleware{
		userRepository: userRepository,
		tokenService:   tokenService,
		systemAPIKey:   systemAPIKey,
	}
}

// If a request is authorized, it will return this context to the controller
// so that information from the user can be used as an immutable object.
type AuthorizedUserContext struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

func (m *AuthMiddleware) Authorize(r *http.Request) (*AuthorizedUserContext, error) {

	cookie, err := r.Cookie("access_token")
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, errors.New("access token not found in request")
	}

	userId, err := m.tokenService.ValidateToken(&cookie.Value)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}

	userEntity, err := m.userRepository.GetUserById(*userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return nil, err
	}

	if userEntity.IsArchived {
		return nil, errors.New("user is archived")
	}

	if !userEntity.IsVerified {
		return nil, errors.New("user is not verified")
	}

	return &AuthorizedUserContext{
		Id:        int(userEntity.ID),
		Email:     userEntity.Email,
		Username:  userEntity.DisplayName,
		CreatedAt: userEntity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil

}

func (m *AuthMiddleware) ValidateSystemAPIKey(r *http.Request) error {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		return errors.New("missing API key")
	}
	if apiKey != m.systemAPIKey {
		return errors.New("invalid API key")
	}
	return nil
}
