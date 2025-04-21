package service

import (
	"errors"
	"fake_server_service/internal/repository"
	"strconv"

	"github.com/flashhhhh/pkg/logging"
)

type FakeServerService interface {
	TriggerNewServer(serverID int) (error)
	CountRunningServers() (int, error)
}

type fakeServerService struct {
	repository repository.FakeServerRepository
}

func NewFakeServerService(repository *repository.FakeServerRepository) FakeServerService {
	return &fakeServerService{
		repository: *repository,
	}
}

func (s *fakeServerService) TriggerNewServer(serverID int) (error) {
	logging.LogMessage("fake_server_service", "Triggering new server with ID: " + strconv.Itoa(serverID), "INFO")

	enabled, err := s.repository.GetServerStatus(serverID)
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to get server status: "+err.Error(), "ERROR")
		return err
	}

	if enabled {
		logging.LogMessage("fake_server_service", "Server is already enabled", "ERROR")
		return errors.New("Server is already enabled")
	}

	err = s.repository.EnableServer(serverID)
	if err != nil {
		logging.LogMessage("fake_server_service", "Failed to enable server: "+err.Error(), "ERROR")
		return err
	}

	logging.LogMessage("fake_server_service", "Server enabled successfully", "INFO")
	return nil
}

func (s *fakeServerService) CountRunningServers() (int, error) {
	count, err := s.repository.CountEnabledServers()
	if err != nil {
		return 0, err
	}
	return count, nil
}