package routes

import (
	"server_administration_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, serverHandler handler.ServerHandler) {
	r.HandleFunc("/server/create", serverHandler.CreateServer).Methods("POST")
}