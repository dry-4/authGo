package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hells/models"
	"hells/services"
	"hells/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get page and limit from query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	page := 1
	limit := 10

	// Parse page
	if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
		page = pageNum
	}

	// Parse limit
	if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
		limit = limitNum
	}

	// Call service to list users
	users, total, err := services.ListUsers(page, limit)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Call service to find user
	user, err := services.FindUserByID(uint(userID))
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	// Clear sensitive data before sending
	user.PasswordHash = ""

	utils.SendJSONResponse(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Parse request body
	var updateData models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateData); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	// Get current user from context (set by AuthMiddleware)
	currentUserID := context.Get(r, "user_id").(uint)
	currentUserRole := context.Get(r, "role").(string)

	// Validate update permissions
	if currentUserID != uint(userID) && currentUserRole != "Admin" {
		utils.SendErrorResponse(w, http.StatusForbidden, "Unauthorized to update this user")
		return
	}

	// Fetch existing user
	existingUser, err := services.FindUserByID(uint(userID))
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	// Update allowed fields
	if updateData.Name != "" {
		existingUser.Name = updateData.Name
	}
	if updateData.Email != "" {
		existingUser.Email = updateData.Email
	}

	// Prevent role change for non-admins
	if currentUserRole == "Admin" && updateData.RoleID != 0 {
		existingUser.RoleID = updateData.RoleID
	}

	// Save updates
	if err := services.UpdateUser(existingUser); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Clear sensitive data before sending
	existingUser.PasswordHash = ""

	utils.SendJSONResponse(w, http.StatusOK, existingUser)
}
