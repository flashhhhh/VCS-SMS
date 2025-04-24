package routes

import (
	"server_administration_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, serverHandler handler.ServerHandler) {
	r.HandleFunc("/server/create", serverHandler.CreateServer).Methods("POST")
	r.HandleFunc("/server/update", serverHandler.UpdateServer).Methods("PUT")
	r.HandleFunc("/server/delete", serverHandler.DeleteServer).Methods("DELETE")
	r.HandleFunc("/server/import", serverHandler.ImportServers).Methods("POST")
}