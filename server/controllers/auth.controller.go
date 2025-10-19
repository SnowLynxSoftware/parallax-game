package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type AuthController struct {
	authMiddleware    middleware.IAuthMiddleware
	authService       services.IAuthService
	shouldEnableHTTPS bool
	cookieDomain      string
}

func NewAuthController(authMiddleware middleware.IAuthMiddleware, authService services.IAuthService, shouldEnableHTTPS bool, cookieDomain string) IController {
	return &AuthController{
		authMiddleware:    authMiddleware,
		authService:       authService,
		shouldEnableHTTPS: shouldEnableHTTPS,
		cookieDomain:      cookieDomain,
	}
}

func (c *AuthController) MapController() *chi.Mux {
	router := chi.NewRouter()
	// Public Routes
	router.Post("/login", c.login)
	router.Post("/logout", c.logout)
	router.Post("/register", c.register)
	router.Get("/verify", c.verify)
	router.Post("/send-login-email", c.sendLoginEmail)
	router.Post("/send-reset-password-email", c.sendResetPasswordEmail)
	router.Post("/reset-password", c.resetPassword)
	router.Get("/login-with-email", c.loginWithEmail)

	// Protected Routes
	router.Get("/token", c.tokenInfo)
	return router
}

func (c *AuthController) logout(w http.ResponseWriter, r *http.Request) {
	c.clearCookie(w, "access_token")
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	response, err := c.authService.Login(&authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	returnStr, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	// log.Info().Str("Access Token: ", response.AccessToken).Msg("")

	c.setCookie(w, "access_token", response.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *AuthController) resetPassword(w http.ResponseWriter, r *http.Request) {

	verificationToken := r.URL.Query().Get("token")

	userId, err := c.authService.VerifyNewUser(&verificationToken)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var userUpdatePasswordDto models.UserUpdatePasswordDTO
	err = json.NewDecoder(r.Body).Decode(&userUpdatePasswordDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.authService.UpdateUserPassword(userId, userUpdatePasswordDto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("password updated successfully"))

}

func (c *AuthController) sendLoginEmail(w http.ResponseWriter, r *http.Request) {
	var userCreateDTO models.UserCreateDTO

	err := json.NewDecoder(r.Body).Decode(&userCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userCreateDTO.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	userEntity, err := c.authService.SendLoginEmail(strings.ToLower(userCreateDTO.Email))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "an error occurred when attempting to send the login email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("user (%v) email was sent.", userEntity.Email)))
}

func (c *AuthController) sendResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
	var userCreateDTO models.UserCreateDTO

	err := json.NewDecoder(r.Body).Decode(&userCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userCreateDTO.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	userEntity, err := c.authService.SendResetPasswordEmail(strings.ToLower(userCreateDTO.Email))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "an error occurred when attempting to send the password reset email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("user (%v) email was sent.", userEntity.Email)))
}

func (c *AuthController) register(w http.ResponseWriter, r *http.Request) {
	var userCreateDTO models.UserCreateDTO

	// Parse JSON body instead of form data
	err := json.NewDecoder(r.Body).Decode(&userCreateDTO)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if userCreateDTO.DisplayName == "" || userCreateDTO.Password == "" || userCreateDTO.Email == "" {
		http.Error(w, "Email, Username, and Password are required", http.StatusBadRequest)
		return
	}

	userEntity, err := c.authService.RegisterNewUser(&userCreateDTO)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "An error occurred when attempting to register your user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Account created successfully! Please check your email (%v) for verification.", userEntity.Email)))
}

func (c *AuthController) loginWithEmail(w http.ResponseWriter, r *http.Request) {

	verificationToken := r.URL.Query().Get("token")

	userId, err := c.authService.VerifyNewUser(&verificationToken)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response, err := c.authService.LoginWithEmailLink(userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	c.setCookie(w, "access_token", response.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully logged in with email link"))
}

func (c *AuthController) verify(w http.ResponseWriter, r *http.Request) {

	verificationToken := r.URL.Query().Get("token")

	_, err := c.authService.VerifyNewUser(&verificationToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Account Verified - Smarter Lynx</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		body {
			font-family: 'Helvetica Neue', Arial, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			display: flex;
			align-items: center;
			justify-content: center;
			padding: 20px;
		}
		.container {
			background: white;
			border-radius: 16px;
			box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
			padding: 48px 32px;
			max-width: 500px;
			width: 100%;
			text-align: center;
		}
		.success-icon {
			width: 80px;
			height: 80px;
			background: #10b981;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			margin: 0 auto 24px;
		}
		.success-icon svg {
			width: 48px;
			height: 48px;
			fill: white;
		}
		h1 {
			font-size: 28px;
			font-weight: 700;
			color: #0f172a;
			margin-bottom: 16px;
		}
		p {
			font-size: 16px;
			line-height: 1.6;
			color: #475569;
			margin-bottom: 32px;
		}
		.redirect-info {
			font-size: 14px;
			color: #64748b;
			margin-bottom: 24px;
		}
		.button {
			display: inline-block;
			background: #2563eb;
			color: white;
			padding: 14px 32px;
			border-radius: 999px;
			text-decoration: none;
			font-weight: 600;
			font-size: 16px;
			transition: background 0.3s ease;
		}
		.button:hover {
			background: #1d4ed8;
		}
	</style>
	<script>
		let countdown = 5;
		function updateCountdown() {
			const element = document.getElementById('countdown');
			if (element) {
				element.textContent = countdown;
			}
			if (countdown === 0) {
				window.location.href = '/login';
			} else {
				countdown--;
				setTimeout(updateCountdown, 1000);
			}
		}
		window.onload = function() {
			updateCountdown();
		};
	</script>
</head>
<body>
	<div class="container">
		<div class="success-icon">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
				<path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z"/>
			</svg>
		</div>
		<h1>Account Verified Successfully!</h1>
		<p>Your account has been verified and is now active. You can now log in and start using Smarter Lynx.</p>
		<p class="redirect-info">Redirecting to login page in <span id="countdown">5</span> seconds...</p>
		<a href="/login" class="button">Go to Login Now</a>
	</div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))

}

func (c *AuthController) tokenInfo(w http.ResponseWriter, r *http.Request) {

	userContext, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "an error occurred when attempting to get token info", http.StatusUnauthorized)
		return
	}

	returnStr, err := json.Marshal(userContext)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)

}

func (c *AuthController) setCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Domain:   c.cookieDomain,
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.shouldEnableHTTPS, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   59 * 60,
	})
}

func (c *AuthController) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Domain:   c.cookieDomain,
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.shouldEnableHTTPS, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
