package handler

import (
	"encoding/json"
	"net/http"
	"user_service/internal/service"

	"github.com/flashhhhh/pkg/logging"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	logging.LogMessage("user_service", "Initializing UserHandler", "INFO")
	
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logging.LogMessage("user_service", "Failed to decode request body for creating user: "+err.Error(), "ERROR")

		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	username, _ := requestBody["username"].(string)
	password, _ := requestBody["password"].(string)
	name, _ := requestBody["name"].(string)
	email, _ := requestBody["email"].(string)
	role, _ := requestBody["role"].(string)
	ctx := r.Context()

	logging.LogMessage("user_service", "Creating user with username: "+username + ", password: "+password + ", name: "+name + ", email: "+email + ", role: "+role, "DEBUG")
	
	user, err := h.userService.CreateUser(ctx, username, password, name, email, role)
	if err != nil {
		logging.LogMessage("user_service", "Failed to create user: "+err.Error(), "ERROR")

		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"userID": user,
	}

	logging.LogMessage("user_service", "User created successfully: "+user, "INFO")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logging.LogMessage("user_service", "Failed to decode request body for logining: "+err.Error(), "ERROR")

		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	username, _ := requestBody["username"].(string)
	password, _ := requestBody["password"].(string)
	ctx := r.Context()

	logging.LogMessage("user_service", "Logging in user with username: "+username + ", password: "+password, "DEBUG")
	token, err := h.userService.Login(ctx, username, password)
	if err != nil {
		logging.LogMessage("user_service", "Failed to login user: "+err.Error(), "ERROR")

		if err.Error() == "Invalid password" {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to login user", http.StatusInternalServerError)
		}

		return
	}

	response := map[string]interface{}{
		"message": "User logged in successfully",
		"token":   token,
	}

	logging.LogMessage("user_service", "User logged in successfully", "INFO")
	logging.LogMessage("user_service", "Token: "+token, "DEBUG")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	ctx := r.Context()

	logging.LogMessage("user_service", "Getting user by ID: "+userID, "DEBUG")
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		logging.LogMessage("user_service", "Failed to get user by ID: "+err.Error(), "ERROR")

		if err.Error() == "User "+userID+" not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user by ID", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message": "User retrieved successfully",
		"user":    user,
	}

	logging.LogMessage("user_service", "User " + userID + " retrieved successfully", "INFO")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logging.LogMessage("user_service", "Getting all users", "DEBUG")
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		logging.LogMessage("user_service", "Failed to get all users: "+err.Error(), "ERROR")

		http.Error(w, "Failed to get all users", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "All users retrieved successfully",
		"users":   users,
	}

	logging.LogMessage("user_service", "All users retrieved successfully", "INFO")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}