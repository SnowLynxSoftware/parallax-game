package controllers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type UIController struct {
	templateService    services.ITemplateService
	staticService      services.IStaticService
	authMiddleware     middleware.IAuthMiddleware
	featureFlagService services.IFeatureFlagService
	teamService        services.ITeamService
	riftService        services.IRiftService
	inventoryService   services.IInventoryService
}

func NewUIController(templateService services.ITemplateService, staticService services.IStaticService, authMiddleware middleware.IAuthMiddleware, featureFlagService services.IFeatureFlagService, teamService services.ITeamService, riftService services.IRiftService, inventoryService services.IInventoryService) IController {
	return &UIController{
		templateService:    templateService,
		staticService:      staticService,
		authMiddleware:     authMiddleware,
		featureFlagService: featureFlagService,
		teamService:        teamService,
		riftService:        riftService,
		inventoryService:   inventoryService,
	}
}

func (c *UIController) MapController() *chi.Mux {
	router := chi.NewRouter()

	// Root redirect to welcome
	router.Get("/", c.redirectToWelcome)

	router.Get("/welcome", c.welcome)
	router.Get("/register", c.register)
	router.Get("/login", c.login)
	router.Get("/teams", c.teams)
	router.Get("/expeditions", c.expeditions)
	router.Get("/fishing", c.fishing)
	router.Get("/inventory", c.inventory)
	router.Get("/dungeons", c.dungeons)
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

func (c *UIController) account(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving account page")

	// Require authentication - redirect to login if not authenticated
	authUser, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(authUser.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool) // Fallback to all unlocked false
	}

	// Prepare page data with authenticated user context
	pageData := services.PageData{
		Title:       "Account Settings",
		Description: "Manage your Parallax account settings",
		Data: map[string]interface{}{
			"Username":    authUser.Username,
			"Email":       authUser.Email,     // Added Email field
			"CreatedAt":   authUser.CreatedAt, // Added CreatedAt field
			"NavbarState": navbarState,
		},
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
		Description: "Parallax - Reset your account password",
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

func (c *UIController) teams(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving teams page")

	// Get authenticated user
	user, err := c.authMiddleware.Authorize(r)
	if err != nil || user == nil {
		util.LogDebug("User not authenticated, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(user.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool)
	}

	// Get teams with expedition status
	teams, err := c.teamService.GetUserTeams(int64(user.Id))
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := services.PageData{
		Title:       "Teams",
		Description: "Manage your expedition teams",
		Data: map[string]interface{}{
			"Username":    user.Username,
			"Teams":       teams,
			"NavbarState": navbarState,
		},
	}

	err = c.templateService.RenderTemplate(w, "teams", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) expeditions(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving expeditions page")

	// Get authenticated user
	user, err := c.authMiddleware.Authorize(r)
	if err != nil || user == nil {
		util.LogDebug("User not authenticated, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(user.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool)
	}

	// Get all teams for this user
	allTeams, err := c.teamService.GetUserTeams(int64(user.Id))
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Filter to only available teams (not on expedition)
	availableTeams := make([]*models.TeamResponseDTO, 0)
	for _, team := range allTeams {
		if team.IsUnlocked && !team.OnExpedition {
			availableTeams = append(availableTeams, team)
		}
	}

	// Get all rifts with unlock status
	rifts, err := c.riftService.GetAllRifts(int64(user.Id))
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := services.PageData{
		Title:       "Expeditions",
		Description: "Launch expeditions to parallel worlds",
		Data: map[string]interface{}{
			"Username":    user.Username,
			"Teams":       availableTeams,
			"Rifts":       rifts,
			"NavbarState": navbarState,
		},
	}

	err = c.templateService.RenderTemplate(w, "expeditions", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) inventory(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving inventory page")

	// Get authenticated user
	user, err := c.authMiddleware.Authorize(r)
	if err != nil || user == nil {
		util.LogDebug("User not authenticated, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(user.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool)
	}

	pageData := services.PageData{
		Title:       "Inventory",
		Description: "View and manage your collected loot",
		Data: map[string]interface{}{
			"Username":    user.Username,
			"NavbarState": navbarState,
		},
	}

	err = c.templateService.RenderTemplate(w, "inventory", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) fishing(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving fishing page")

	// Get authenticated user
	user, err := c.authMiddleware.Authorize(r)
	if err != nil || user == nil {
		util.LogDebug("User not authenticated, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(user.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool)
	}

	pageData := services.PageData{
		Title:       "Fishing",
		Description: "View and manage your fishing activities",
		Data: map[string]interface{}{
			"Username":    user.Username,
			"NavbarState": navbarState,
		},
	}

	err = c.templateService.RenderTemplate(w, "fishing", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) dungeons(w http.ResponseWriter, r *http.Request) {
	util.LogDebug("Serving dungeons page")

	// Get authenticated user
	user, err := c.authMiddleware.Authorize(r)
	if err != nil || user == nil {
		util.LogDebug("User not authenticated, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get navbar unlock state
	navbarState, err := c.getNavbarUnlockState(int64(user.Id))
	if err != nil {
		util.LogError(err)
		navbarState = make(map[string]bool)
	}

	// Fetch external content from weshould.run
	externalContent, scripts, styles, fetchErr := c.fetchExternalContent("https://weshould.run/")
	if fetchErr != nil {
		util.LogError(fetchErr)
		util.LogDebug("Failed to fetch external content, redirecting to /teams")
		http.Redirect(w, r, "/teams", http.StatusSeeOther)
		return
	}

	_ = externalContent // Keep for potential future use

	pageData := services.PageData{
		Title:       "Dungeons",
		Description: "Explore the dungeons",
		Data: map[string]interface{}{
			"Username":      user.Username,
			"Scripts":       scripts,
			"InlineScripts": []string{}, // Currently extracting only external scripts
			"Styles":        styles,
			"NavbarState":   navbarState,
		},
	}

	fmt.Println(pageData.Data)

	err = c.templateService.RenderTemplate(w, "dungeons", pageData)
	if err != nil {
		util.LogError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// fetchExternalContent fetches HTML from a URL and extracts script and style tags
func (c *UIController) fetchExternalContent(url string) (string, []string, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		util.LogDebug("Unexpected status code: " + string(rune(resp.StatusCode)))
		return "", nil, nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}

	htmlContent := string(body)
	scripts := extractScriptSrcs(htmlContent)
	styles := extractStylesAndLinks(htmlContent)

	return htmlContent, scripts, styles, nil
}

// extractScriptSrcs extracts script src URLs from HTML using regex
func extractScriptSrcs(htmlContent string) []string {
	var scripts []string

	// Match <script src="..."> tags
	scriptSrcRegex := regexp.MustCompile(`<script\s+[^>]*src\s*=\s*["']([^"']+)["'][^>]*>`)
	matches := scriptSrcRegex.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			scripts = append(scripts, match[1])
		}
	}

	return scripts
}

// extractStylesAndLinks extracts style content and stylesheet links from HTML using regex
func extractStylesAndLinks(htmlContent string) []string {
	var styles []string

	// Match <link rel="stylesheet" href="..."> tags
	linkRegex := regexp.MustCompile(`<link\s+[^>]*rel\s*=\s*["']stylesheet["'][^>]*href\s*=\s*["']([^"']+)["'][^>]*>`)
	linkMatches := linkRegex.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range linkMatches {
		if len(match) > 1 && match[1] != "" {
			styles = append(styles, match[1])
		}
	}

	// Also match href first format: <link ... href="..." ... rel="stylesheet">
	linkRegex2 := regexp.MustCompile(`<link\s+[^>]*href\s*=\s*["']([^"']+)["'][^>]*rel\s*=\s*["']stylesheet["'][^>]*>`)
	linkMatches2 := linkRegex2.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range linkMatches2 {
		if len(match) > 1 && match[1] != "" {
			styles = append(styles, match[1])
		}
	}

	// Match <style>...</style> tags
	styleRegex := regexp.MustCompile(`<style[^>]*>([\s\S]*?)</style>`)
	styleMatches := styleRegex.FindAllStringSubmatch(htmlContent, -1)

	for _, match := range styleMatches {
		if len(match) > 1 && strings.TrimSpace(match[1]) != "" {
			styles = append(styles, match[1])
		}
	}

	return styles
}

// getNavbarUnlockState fetches the unlock status for navbar features
func (c *UIController) getNavbarUnlockState(userId int64) (map[string]bool, error) {
	unlockState := make(map[string]bool)

	hasFishingRod, err := c.inventoryService.HasItemByName(userId, "Golden Fishing Rod")
	if err != nil {
		util.LogError(err)
		return nil, err
	}

	hasCompass, err := c.inventoryService.HasItemByName(userId, "Explorers Compass")
	if err != nil {
		util.LogError(err)
		return nil, err
	}

	unlockState["HasFishingRod"] = hasFishingRod
	unlockState["HasCompass"] = hasCompass

	return unlockState, nil
}
