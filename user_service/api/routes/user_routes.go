package api

import (
	"user_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler handler.UserHandler) {
	r.HandleFunc("/create", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/getUserByID", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/getAllUsers", userHandler.GetAllUsers).Methods("GET")
}