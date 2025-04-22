package service

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/repository"
)

type ServerService interface {
	CreateServer(server_id, server_name, status, ipv4 string, port int) error
}

type serverService struct {
	serverRepository repository.ServerRepository
}

func NewServerService(serverRepository repository.ServerRepository) ServerService {
	return &serverService{
		serverRepository: serverRepository,
	}
}

func (s *serverService) CreateServer(server_id, server_name, status, ipv4 string, port int) error {
	server := &domain.Server{
		ServerID:   server_id,
		ServerName: server_name,
		Status:     status,
		IPAddress:  ipv4,
		Port: 	 port,
	}

	err := s.serverRepository.CreateServer(server)
	if err != nil {
		return err
	}
	return nil
}