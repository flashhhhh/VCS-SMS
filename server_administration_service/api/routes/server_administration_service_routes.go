package routes

import (
	"net/http"
	"server_administration_service/api/middlewares"
	"server_administration_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, serverHandler handler.ServerHandler) {
	r.Handle("/create", middlewares.UserMiddleware(http.HandlerFunc(serverHandler.CreateServer))).Methods("POST")
	r.Handle("/view", middlewares.GuestMiddleware(http.HandlerFunc(serverHandler.ViewServers))).Methods("GET")
	r.Handle("/update", middlewares.AdminMiddleware(http.HandlerFunc(serverHandler.UpdateServer))).Methods("PUT")
	r.Handle("/delete", middlewares.AdminMiddleware(http.HandlerFunc(serverHandler.DeleteServer))).Methods("DELETE")
	r.Handle("/import", middlewares.UserMiddleware(http.HandlerFunc(serverHandler.ImportServers))).Methods("POST")
	r.Handle("/export", middlewares.AdminMiddleware(http.HandlerFunc(serverHandler.ExportServers))).Methods("GET")
}