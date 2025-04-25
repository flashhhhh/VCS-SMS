package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/service"
	"strconv"
	"time"

	"github.com/flashhhhh/pkg/logging"
)

type ServerHandler interface {
	CreateServer(w http.ResponseWriter, r *http.Request)
	ViewServers(w http.ResponseWriter, r *http.Request)
	UpdateServer(w http.ResponseWriter, r *http.Request)
	DeleteServer(w http.ResponseWriter, r *http.Request)
	ImportServers(w http.ResponseWriter, r *http.Request)
	ExportServers(w http.ResponseWriter, r *http.Request)
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
	portFloat, ok := requestBody["port"].(float64)

	port := 80 // default port
	if ok {
		port = int(portFloat)
	}
	
	id, err := h.service.CreateServer(serverID, serverName, status, ipAddress, port)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to create server: "+err.Error(), "ERROR")
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Server created successfully with ID: "+strconv.Itoa(id), "INFO")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Server created successfully"))
}

func (h *serverHandler) ViewServers(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid 'from' query parameter: "+err.Error(), "ERROR")
		logging.LogMessage("server_administration_service", "from: "+fromStr, "DEBUG")
		http.Error(w, "Invalid 'from' query parameter", http.StatusBadRequest)
		return
	}

	toStr := r.URL.Query().Get("to")
	to, err := strconv.Atoi(toStr)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid 'to' query parameter: "+err.Error(), "ERROR")
		logging.LogMessage("server_administration_service", "to: "+toStr, "DEBUG")
		http.Error(w, "Invalid 'to' query parameter", http.StatusBadRequest)
		return
	}

	sortedColumn := r.URL.Query().Get("sort_column")
	order := r.URL.Query().Get("sort_order")

	serverID := r.URL.Query().Get("server_id")
	serverName := r.URL.Query().Get("server_name")
	status := r.URL.Query().Get("status")
	ipv4 := r.URL.Query().Get("ipv4")
	portStr := r.URL.Query().Get("port")

	serverFilter := dto.ServerFilter{}

	if serverID != "" {
		serverFilter.ServerID = serverID
	}
	if serverName != "" {
		serverFilter.ServerName = serverName
	}
	if status != "" {
		serverFilter.Status = status
	}
	if ipv4 != "" {
		serverFilter.IPv4 = ipv4
	}
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			logging.LogMessage("server_administration_service", "Invalid port in query params", "ERROR")
			http.Error(w, "Invalid port", http.StatusBadRequest)
			return
		}
		serverFilter.Port = port
	} else {
		serverFilter.Port = -1
	}

	servers, err := h.service.ViewServers(&serverFilter, from, to, sortedColumn, order)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to view servers: "+err.Error(), "ERROR")
		http.Error(w, "Failed to view servers", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Servers retrieved successfully", "INFO")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(servers)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to marshal servers response: "+err.Error(), "ERROR")
		http.Error(w, "Failed to process servers data", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (h *serverHandler) UpdateServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		logging.LogMessage("server_administration_service", "Server ID is required for update", "ERROR")
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to decode request body for request UpdateServer: "+err.Error(), "ERROR")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedData := make(map[string]interface{})
	serverName, existed := requestBody["server_name"].(string)
	if existed {
		updatedData["server_name"] = serverName
	}

	status, existed := requestBody["status"].(string)
	if existed {
		updatedData["status"] = status
	}

	ipAddress, existed := requestBody["ipv4"].(string)
	if existed {
		updatedData["ipv4"] = ipAddress
	}

	portFloat, existed := requestBody["port"].(float64)
	if existed {
		updatedData["port"] = int(portFloat)
	}

	err = h.service.UpdateServer(serverID, updatedData)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to update server: "+err.Error(), "ERROR")
		http.Error(w, "Failed to update server", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Server updated successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server updated successfully"))
}

func (h *serverHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		logging.LogMessage("server_administration_service", "Server ID is required for deletion", "ERROR")
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteServer(serverID)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid server ID: "+serverID+" - "+err.Error(), "ERROR")
		http.Error(w, "Invalid server ID", http.StatusNotFound)
		return
	}

	logging.LogMessage("server_administration_service", "Server deleted successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server deleted successfully"))
}

func (h *serverHandler) ImportServers(w http.ResponseWriter, r *http.Request) {
	serversFile, _, err := r.FormFile("servers_file")

	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to get file from request: "+err.Error(), "ERROR")
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer serversFile.Close()

	var buf []byte
	buf, err = io.ReadAll(serversFile)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to read file: "+err.Error(), "ERROR")
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	err = h.service.ImportServers(buf)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to import servers: "+err.Error(), "ERROR")
		http.Error(w, "Failed to import servers", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Servers imported successfully", "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Servers imported successfully"))
}

func (h *serverHandler) ExportServers(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid 'from' query parameter: "+err.Error(), "ERROR")
		logging.LogMessage("server_administration_service", "from: "+fromStr, "DEBUG")
		http.Error(w, "Invalid 'from' query parameter", http.StatusBadRequest)
		return
	}

	toStr := r.URL.Query().Get("to")
	to, err := strconv.Atoi(toStr)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid 'to' query parameter: "+err.Error(), "ERROR")
		logging.LogMessage("server_administration_service", "to: "+toStr, "DEBUG")
		http.Error(w, "Invalid 'to' query parameter", http.StatusBadRequest)
		return
	}

	sortedColumn := r.URL.Query().Get("sort_column")
	order := r.URL.Query().Get("sort_order")

	serverID := r.URL.Query().Get("server_id")
	serverName := r.URL.Query().Get("server_name")
	status := r.URL.Query().Get("status")
	ipv4 := r.URL.Query().Get("ipv4")
	portStr := r.URL.Query().Get("port")

	serverFilter := dto.ServerFilter{}

	if serverID != "" {
		serverFilter.ServerID = serverID
	}
	if serverName != "" {
		serverFilter.ServerName = serverName
	}
	if status != "" {
		serverFilter.Status = status
	}
	if ipv4 != "" {
		serverFilter.IPv4 = ipv4
	}
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			logging.LogMessage("server_administration_service", "Invalid port in query params", "ERROR")
			http.Error(w, "Invalid port", http.StatusBadRequest)
			return
		}
		serverFilter.Port = port
	} else {
		serverFilter.Port = -1
	}

	serverBuf, err := h.service.ExportServers(&serverFilter, from, to, sortedColumn, order)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to export servers: "+err.Error(), "ERROR")
		http.Error(w, "Failed to export servers", http.StatusInternalServerError)
		return
	}

	filename := "servers_" + time.Now().Format("2006-01-02_15-04-05") + ".xlsx"

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("File-Name", filename)
	w.WriteHeader(http.StatusOK)
	w.Write(serverBuf)
}