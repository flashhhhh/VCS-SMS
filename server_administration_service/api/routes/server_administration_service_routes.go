package routes

import (
	"net/http"
	"server_administration_service/api/middlewares"
	"server_administration_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, serverHandler handler.ServerHandler) {
	r.HandleFunc("/server/create", serverHandler.CreateServer).Methods("POST")
	r.HandleFunc("/server/view", serverHandler.ViewServers).Methods("GET")
	r.HandleFunc("/server/update", serverHandler.UpdateServer).Methods("PUT")
	r.HandleFunc("/server/delete", serverHandler.DeleteServer).Methods("DELETE")
	r.HandleFunc("/server/import", serverHandler.ImportServers).Methods("POST")
	r.Handle("/server/export", middlewares.AdminMiddleware(http.HandlerFunc(serverHandler.ExportServers))).Methods("GET")
}