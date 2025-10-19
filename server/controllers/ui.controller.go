package controllers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type UIController struct {
	templateService    services.ITemplateService
	staticService      services.IStaticService
	authMiddleware     middleware.IAuthMiddleware
	featureFlagService services.IFeatureFlagService
}

func NewUIController(templateService services.ITemplateService, staticService services.IStaticService, authMiddleware middleware.IAuthMiddleware, featureFlagService services.IFeatureFlagService) IController {
	return &UIController{
		templateService:    templateService,
		staticService:      staticService,
		authMiddleware:     authMiddleware,
		featureFlagService: featureFlagService,
	}
}

func (c *UIController) MapController() *chi.Mux {
	router := chi.NewRouter()

	// Root redirect to welcome
	router.Get("/", c.redirectToWelcome)

	router.Get("/welcome", c.welcome)
	router.Get("/register", c.register)
	router.Get("/login", c.login)
	router.Get("/dashboard", c.dashboard)
	router.Get("/account", c.account)
	router.Get("/reset-password", c.resetPassword)
	router.Get("/terms", c.terms)
	router.Get("/privacy", c.privacy)

	// Static file serving
	router.Get("/static/*", c.serveStatic)

	return router
}

func (c *UIController) redirectToWelcome(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Redirecting root to welcome page")
	http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)
}

func (c *UIController) welcome(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving welcome page")

	// Check if we're in prelaunch mode
	isPrelaunchMode := c.featureFlagService.IsEnabled("prelaunch_mode")

	pageData := services.PageData{
		Title:       "Welcome",
		Description: "Welcome to Parallax - A journey through the rifts of space-time",
		Data: map[string]interface{}{
			"PrelaunchMode": isPrelaunchMode,
		},
	}

	err := c.templateService.RenderTemplate(w, "welcome", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) register(w http.ResponseWriter, r *http.Request) {
	// Check if we're in prelaunch mode--and redirect to welcome page if so.
	isPrelaunchMode := c.featureFlagService.IsEnabled("prelaunch_mode")
	if isPrelaunchMode {
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	util.LogDebug("Serving register page")

	pageData := services.PageData{
		Title:       "Register",
		Description: "Welcome to Parallax - Register for an account to get started",
		Data:        nil,
	}

	err := c.templateService.RenderTemplate(w, "register", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) login(w http.ResponseWriter, r *http.Request) {
	// Check if we're in prelaunch mode--and redirect to welcome page if so.
	isPrelaunchMode := c.featureFlagService.IsEnabled("prelaunch_mode")
	if isPrelaunchMode {
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	util.LogDebug("Serving login page")

	pageData := services.PageData{
		Title:       "Login",
		Description: "Login to Parallax",
		Data:        nil,
	}

	err := c.templateService.RenderTemplate(w, "login", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) dashboard(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving dashboard page")

	authUser, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Build dashboard data
	dashboardData := map[string]interface{}{
		"Username":   authUser.Username,
	}

	pageData := services.PageData{
		Title:       "Dashboard",
		Description: "Parallax Dashboard",
		Data:        dashboardData,
	}

	err = c.templateService.RenderTemplate(w, "dashboard", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) account(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving account page")

	// Require authentication - redirect to login if not authenticated
	authUser, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Prepare page data with authenticated user context
	pageData := services.PageData{
		Title:       "Account Settings",
		Description: "Manage your Parallax account settings",
		Data:        authUser,
	}

	// Render the account template
	err = c.templateService.RenderTemplate(w, "account", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) terms(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving terms of service page")

	pageData := services.PageData{
		Title:       "Terms of Service",
		Description: "Parallax Terms of Service - Read our terms and conditions",
		Data: map[string]interface{}{
			"LastUpdated": "October 19, 2025",
		},
	}

	err := c.templateService.RenderTemplate(w, "terms", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) resetPassword(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving reset password page")

	pageData := services.PageData{
		Title:       "Reset Password",
		Description: "Smarter Lynx - Reset your account password",
		Data:        map[string]interface{}{},
	}

	err := c.templateService.RenderTemplate(w, "reset-password", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) privacy(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving privacy policy page")

	pageData := services.PageData{
		Title:       "Privacy Policy",
		Description: "Parallax - Learn how we protect your data",
		Data: map[string]interface{}{
			"PrivacyLastUpdated": "October 19, 2025",
		},
	}

	err := c.templateService.RenderTemplate(w, "privacy", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) serveStatic(w http.ResponseWriter, r *http.Request) {
	// Extract the file path from the URL
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")

	err := c.staticService.ServeStaticFile(w, r, filePath)
	if err != nil {
		// Error is already handled and logged in the static service
		return
	}
}
