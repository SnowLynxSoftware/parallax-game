package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type InventoryController struct {
	inventoryService services.IInventoryService
	authMiddleware   middleware.IAuthMiddleware
}

func NewInventoryController(inventoryService services.IInventoryService, authMiddleware middleware.IAuthMiddleware) *InventoryController {
	return &InventoryController{
		inventoryService: inventoryService,
		authMiddleware:   authMiddleware,
	}
}

func (c *InventoryController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.getUserInventory)
	r.Get("/equipment", c.getEquipment)
	r.Get("/consumables", c.getConsumables)
	return r
}

func (c *InventoryController) getUserInventory(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	inventory, err := c.inventoryService.GetUserInventory(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func (c *InventoryController) getEquipment(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	equipment, err := c.inventoryService.GetEquipment(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

func (c *InventoryController) getConsumables(w http.ResponseWriter, r *http.Request) {
	user, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	consumables, err := c.inventoryService.GetConsumables(int64(user.Id))
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consumables)
}
