package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type RiftController struct {
	riftService    services.IRiftService
	authMiddleware middleware.IAuthMiddleware
}

func NewRiftController(riftService services.IRiftService, authMiddleware middleware.IAuthMiddleware) *RiftController {
	return &RiftController{
		riftService:    riftService,
		authMiddleware: authMiddleware,
	}
}

func (c *RiftController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.getAllRifts)
	r.Get("/{riftId}", c.getRiftById)
	return r
}

func (c *RiftController) getAllRifts(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rifts, err := c.riftService.GetAllRifts(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rifts)
}

func (c *RiftController) getRiftById(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	riftIdStr := chi.URLParam(r, "riftId")
	if riftIdStr == "" {
		http.Error(w, "Rift ID is required", http.StatusBadRequest)
		return
	}

	riftId, err := strconv.ParseInt(riftIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rift ID", http.StatusBadRequest)
		return
	}

	rift, err := c.riftService.GetRiftById(riftId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rift)
}
