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

type TeamController struct {
	teamService    services.ITeamService
	authMiddleware middleware.IAuthMiddleware
}

func NewTeamController(teamService services.ITeamService, authMiddleware middleware.IAuthMiddleware) *TeamController {
	return &TeamController{
		teamService:    teamService,
		authMiddleware: authMiddleware,
	}
}

func (c *TeamController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.getUserTeams)
	r.Get("/{teamId}", c.getTeamById)
	r.Post("/equip", c.equipItem)
	r.Post("/unequip", c.unequipItem)
	r.Post("/consume", c.consumeItem)
	r.Post("/{teamId}/unlock", c.unlockTeam)
	return r
}

func (c *TeamController) getUserTeams(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teams, err := c.teamService.GetUserTeams(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (c *TeamController) getTeamById(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teamIdStr := chi.URLParam(r, "teamId")
	teamId, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := c.teamService.GetTeamById(teamId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (c *TeamController) equipItem(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto models.EquipItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	team, err := c.teamService.EquipItemToTeam(int64(user.Id), dto.TeamID, dto.Slot, dto.InventoryID)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (c *TeamController) unequipItem(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto models.UnequipItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	team, err := c.teamService.UnequipItemFromTeam(int64(user.Id), dto.TeamID, dto.Slot)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (c *TeamController) consumeItem(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var dto models.ConsumeItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	team, err := c.teamService.ConsumeItemOnTeam(int64(user.Id), dto.TeamID, dto.InventoryID)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (c *TeamController) unlockTeam(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teamIdStr := chi.URLParam(r, "teamId")
	teamId, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	err = c.teamService.UnlockTeam(int64(user.Id), teamId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("team unlocked"))
}
