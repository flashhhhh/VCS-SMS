package service

import (
	"errors"
	"fake_server_service/internal/repository"
	"net"
	"strconv"
	"time"

	"github.com/flashhhhh/pkg/logging"
)

type FakeServerService interface {
	CheckServer(serverID int) (bool, error)
	HostServer(port, duration int)
	CountRunningServers() (int, error)
	DeleteServers() error
}

type fakeServerService struct {
	repository repository.FakeServerRepository
}

func NewFakeServerService(repository *repository.FakeServerRepository) FakeServerService {
	return &fakeServerService{
		repository: *repository,
	}
}

func (s *fakeServerService) CheckServer(serverID int) (bool, error) {
	logging.LogMessage("fake_server_service", "Triggering new server with ID: " + strconv.Itoa(serverID), "INFO")

	enabled, err := s.repository.GetServerStatus(serverID)
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to get server status: "+err.Error(), "ERROR")
		return false, err
	}

	if enabled {
		logging.LogMessage("fake_server_service", "Server is already enabled", "ERROR")
		return false, errors.New("Server is already enabled")
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(serverID))
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to start server: "+err.Error(), "ERROR")
		return false, err
	}
	defer listener.Close()

	return true, nil
}

func (s *fakeServerService) HostServer(port, duration int) {
	logging.LogMessage("fake_server_service", "Hosting new server on port: "+strconv.Itoa(port), "INFO")

	s.repository.EnableServer(port)
	defer func() {
		s.repository.DisableServer(port)
		logging.LogMessage("fake_server_service", "Server on port "+strconv.Itoa(port)+" has been disabled", "INFO")
	}()
	
	listener, _ := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	
	logging.LogMessage("fake_server_service", "Server will run for: "+strconv.Itoa(duration)+" seconds", "INFO")
	time.Sleep(time.Duration(duration) * time.Second)

	logging.LogMessage("fake_server_service", "Server is shutting down", "INFO")
	listener.Close()
}

func (s *fakeServerService) CountRunningServers() (int, error) {
	count, err := s.repository.CountEnabledServers()
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to count running servers: "+err.Error(), "ERROR")
		return 1000000000, err
	}
	return count, nil
}

func (s *fakeServerService) DeleteServers() error {
	err := s.repository.DeleteServers()
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to delete servers: "+err.Error(), "ERROR")
		return err
	}

	logging.LogMessage("fake_server_service", "Servers deleted successfully", "INFO")
	return nil
}