package handler

import (
	"encoding/json"
	"net/http"
	"server_administration_service/internal/service"

	"github.com/flashhhhh/pkg/logging"
)

type ServerHandler interface {
	CreateServer(w http.ResponseWriter, r *http.Request)
}

type serverHandler struct {
	service service.ServerService
}

func NewServerHandler(service service.ServerService) ServerHandler {
	return &serverHandler{
		service: service,
	}
}

func (h *serverHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to decode request body for request CreateServer: "+err.Error(), "ERROR")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	serverID, _ := requestBody["server_id"].(string)
	serverName, _ := requestBody["server_name"].(string)
	status, _ := requestBody["status"].(string)
	ipAddress, _ := requestBody["ipv4"].(string)
	port, ok := requestBody["port"].(int)
	if !ok {
		port = 80 // Default value if port is not provided
	}
	
	err = h.service.CreateServer(serverID, serverName, status, ipAddress, port)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to create server: "+err.Error(), "ERROR")
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Server created successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Server created successfully"))
}