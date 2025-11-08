package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/models"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type ExpeditionController struct {
	expeditionService services.IExpeditionService
	authMiddleware    middleware.IAuthMiddleware
}

func NewExpeditionController(expeditionService services.IExpeditionService, authMiddleware middleware.IAuthMiddleware) *ExpeditionController {
	return &ExpeditionController{
		expeditionService: expeditionService,
		authMiddleware:    authMiddleware,
	}
}

func (c *ExpeditionController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/start", c.startExpedition)
	r.Get("/active", c.getActiveExpeditions)
	r.Get("/history", c.getExpeditionHistory)
	r.Post("/{expeditionId}/claim", c.claimRewards)
	return r
}

func (c *ExpeditionController) startExpedition(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto models.StartExpeditionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expedition, err := c.expeditionService.StartExpedition(int64(user.Id), dto.TeamID, dto.RiftID)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expedition)
}

func (c *ExpeditionController) getActiveExpeditions(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expeditions, err := c.expeditionService.GetActiveExpeditions(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expeditions)
}

func (c *ExpeditionController) getExpeditionHistory(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get limit from query param (default 10)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	expeditions, err := c.expeditionService.GetExpeditionHistory(int64(user.Id), limit)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expeditions)
}

func (c *ExpeditionController) claimRewards(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expeditionIdStr := chi.URLParam(r, "expeditionId")
	expeditionId, err := strconv.ParseInt(expeditionIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expedition ID", http.StatusBadRequest)
		return
	}

	rewards, err := c.expeditionService.ClaimExpeditionRewards(int64(user.Id), expeditionId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rewards)
}
