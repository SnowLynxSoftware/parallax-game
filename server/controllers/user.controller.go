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

type UserController struct {
	userService    services.IUserService
	authMiddleware middleware.IAuthMiddleware
}

func NewUserController(userService services.IUserService, authMiddleware middleware.IAuthMiddleware) *UserController {
	return &UserController{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
}

func (c *UserController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", c.getUsers)
	r.Get("/{id}", c.getUserById)
	r.Put("/{id}", c.updateUser)
	r.Patch("/{id}/archived", c.toggleUserArchived)
	return r
}

func (c *UserController) toggleUserArchived(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to update this user", http.StatusForbidden)
		return
	}
	userIdStr := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	err = c.userService.ToggleUserArchived(&userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to toggle user archived status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (c *UserController) updateUser(w http.ResponseWriter, r *http.Request) {
	userContext, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you must be logged in to perform this request", http.StatusUnauthorized)
		return
	}

	userIdStr := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	if userContext.Id != userId {
		http.Error(w, "you are not authorized to update this user", http.StatusForbidden)
		return
	}

	var userUpdateDTO models.UserUpdateDTO
	err = json.NewDecoder(r.Body).Decode(&userUpdateDTO)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	if userUpdateDTO.Email == "" || userUpdateDTO.DisplayName == "" {
		http.Error(w, "email and display name are required", http.StatusBadRequest)
		return
	}

	updatedUser, err := c.userService.UpdateUser(&userUpdateDTO, &userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}
	if updatedUser == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	returnStr, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *UserController) getUserById(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	userIdStr := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUserById(userId)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
		return
	}

	returnStr, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}

func (c *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	_, err := c.authMiddleware.Authorize(r)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "you are not authorized to perform this request", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	pageSize := 25       // Default page size
	page := 1            // Default page number
	searchString := ""   // Default search string
	statusFilter := ""   // Default status filter
	userTypeFilter := "" // Default user type filter

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 {
			pageSize = psInt
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		searchString = search
	}
	if status := r.URL.Query().Get("status"); status != "" {
		statusFilter = status
	}
	if userType := r.URL.Query().Get("user_type"); userType != "" {
		userTypeFilter = userType
	}

	offset := (page - 1) * pageSize

	results, err := c.userService.GetUsers(pageSize, offset, searchString, statusFilter, userTypeFilter)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		http.Error(w, "failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// TODO: This is a temporary fix to ensure the page number is right, but I need to fix this in the service layer
	results.Page = page

	returnStr, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnStr)
}
