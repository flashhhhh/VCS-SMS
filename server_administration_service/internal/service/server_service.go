package service

import (
	"bytes"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/flashhhhh/pkg/logging"
	"github.com/xuri/excelize/v2"
)

type ServerService interface {
	CreateServer(server_id, server_name, status, ipv4 string, port int) (int, error)
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
	DeleteServer(server_id string) error
	ImportServers(buf []byte) ([]domain.Server, []domain.Server, error)
	ExportServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]byte, error)
	
	UpdateServerStatus(id int, status string) error
	GetAllAddresses() ([]dto.ServerAddress, error)

	AddServerStatus(id int, status string) error
	GetNumOnServers() (int, error)
	GetNumServers() (int, error)
	GetServerUptimeRatio(startTime, endTime time.Time) (float64, error)
}

type serverService struct {
	serverRepository repository.ServerRepository
}

func NewServerService(serverRepository repository.ServerRepository) ServerService {
	return &serverService{
		serverRepository: serverRepository,
	}
}

func (s *serverService) CreateServer(server_id, server_name, status, ipv4 string, port int) (int, error) {
	server := &domain.Server{
		ServerID:   server_id,
		ServerName: server_name,
		Status:     status,
		IPv4:  ipv4,
		Port: 	 port,
	}

	id, err := s.serverRepository.CreateServer(server)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *serverService) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error) {
	servers, err := s.serverRepository.ViewServers(serverFilter, from, to, sortedColumn, order)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (s *serverService) UpdateServer(server_id string, updatedData map[string]interface{}) error {
	err := s.serverRepository.UpdateServer(server_id, updatedData)
	return err
}

func (s *serverService) DeleteServer(server_id string) error {
	err := s.serverRepository.DeleteServer(server_id)
	return err
}

func (s *serverService) UpdateServerStatus(id int, status string) error {
	err := s.serverRepository.UpdateServerStatus(id, status)
	return err
}

func (s *serverService) GetAllAddresses() ([]dto.ServerAddress, error) {
	addresses, err := s.serverRepository.GetAllAddresses()
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *serverService) ImportServers(buf []byte) ([]domain.Server, []domain.Server, error) {
	f, err := excelize.OpenReader(strings.NewReader(string(buf)))
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to open Excel file: "+err.Error(), "ERROR")
		return nil, nil, err
	}

	rows, err := f.GetRows("Servers")
	if err != nil || len(rows) < 2 {
		logging.LogMessage("server_administration_service", "Failed to get rows from Excel file: "+err.Error(), "ERROR")
		return nil, nil, err
	}

	servers := make([]domain.Server, 0)

	for _, row := range rows[1:] {
		serverID := row[0]
		serverName := row[1]
		status := row[2]
		ipv4 := row[3]
		
		port := 80
		if len(row) > 4 {
			port, err = strconv.Atoi(row[4])
			if err != nil {
				logging.LogMessage("server_administration_service", "Invalid port value: "+row[4], "ERROR")
				return nil, nil, err
			}
		}

		server := domain.Server{
			ServerID:   serverID,
			ServerName: serverName,
			Status:     status,
			IPv4:       ipv4,
			Port:       port,
		}

		servers = append(servers, server)
	}

	insertedServers, nonInsertedServers, err := s.serverRepository.CreateServers(servers)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to insert servers: "+err.Error(), "ERROR")
		return nil, nil, err
	}
	
	logging.LogMessage("server_administration_service", "Servers imported successfully", "INFO")
	return insertedServers, nonInsertedServers, nil	
}

func (s *serverService) ExportServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]byte, error) {
	servers, err := s.serverRepository.ViewServers(serverFilter, from, to, sortedColumn, order)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to export servers: "+err.Error(), "ERROR")
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Servers"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Server ID", "Server Name", "Status", "IPv4", "Port"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for i, server := range servers {
		f.SetCellValue(sheet, "A"+strconv.Itoa(i+2), server.ServerID)
		f.SetCellValue(sheet, "B"+strconv.Itoa(i+2), server.ServerName)
		f.SetCellValue(sheet, "C"+strconv.Itoa(i+2), server.Status)
		f.SetCellValue(sheet, "D"+strconv.Itoa(i+2), server.IPv4)
		f.SetCellValue(sheet, "E"+strconv.Itoa(i+2), server.Port)
	}

	var buf bytes.Buffer
	err = f.Write(&buf)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to write Excel file to buffer: "+err.Error(), "ERROR")
		return nil, err
	}

	logging.LogMessage("server_administration_service", "Servers exported successfully", "INFO")
	return buf.Bytes(), nil
}

func (s *serverService) AddServerStatus(id int, status string) error {
	err := s.serverRepository.AddServerStatus(id, status)
	return err
}

func (s *serverService) GetNumOnServers() (int, error) {
	return s.serverRepository.GetNumOnServers()
}

func (s *serverService) GetNumServers() (int, error) {
	return s.serverRepository.GetNumServers()
}

func (s *serverService) GetServerUptimeRatio(startTime, endTime time.Time) (float64, error) {
	return s.serverRepository.GetServerUptimeRatio(startTime, endTime)
}