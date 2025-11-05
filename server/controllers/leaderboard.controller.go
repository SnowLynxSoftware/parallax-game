package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type LeaderboardController struct {
	service        services.ILeaderboardService
	authMiddleware middleware.IAuthMiddleware
}

func NewLeaderboardController(service services.ILeaderboardService, authMiddleware middleware.IAuthMiddleware) *LeaderboardController {
	return &LeaderboardController{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

func (c *LeaderboardController) MapController() *chi.Mux {
	r := chi.NewRouter()

	// API routes
	r.Get("/{type}", c.getLeaderboard)

	return r
}

func (c *LeaderboardController) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	leaderboardType := chi.URLParam(r, "type")

	// Validate type
	validTypes := map[string]bool{
		"legendary":   true,
		"power":       true,
		"expeditions": true,
	}
	if !validTypes[leaderboardType] {
		http.Error(w, "Invalid leaderboard type", http.StatusBadRequest)
		return
	}

	result, err := c.service.GetLeaderboard(leaderboardType, user.Id)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Leaderboard temporarily unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
