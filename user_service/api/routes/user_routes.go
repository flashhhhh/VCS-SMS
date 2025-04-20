package api

import (
	"user_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler handler.UserHandler) {
	r.HandleFunc("/user/create", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/user/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/user/getUserByID", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/user/getAllUsers", userHandler.GetAllUsers).Methods("GET")
}