package service

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/repository"
	"strconv"
)

type ServerService interface {
	CreateServer(server_id, server_name, status, ipv4 string, port int) error
	DeleteServer(server_id string) error
	
	UpdateServerStatus(server_id, status string) error
	GetAllAddresses() ([][2]string, error)
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
		IPv4:  ipv4,
		Port: 	 port,
	}

	err := s.serverRepository.CreateServer(server)
	if err != nil {
		return err
	}
	return nil
}

func (s *serverService) DeleteServer(server_id string) error {
	err := s.serverRepository.DeleteServer(server_id)
	return err
}

func (s *serverService) UpdateServerStatus(server_id, status string) error {
	err := s.serverRepository.UpdateServerStatus(server_id, status)
	return err
}

func (s *serverService) GetAllAddresses() ([][2]string, error) {
	addresses, err := s.serverRepository.GetAllAddresses()
	if err != nil {
		return nil, err
	}

	var addressList [][2]string
	for _, addressInfo := range addresses {
		address := addressInfo.IPv4 + ":" + strconv.Itoa(addressInfo.Port)
		addressList = append(addressList, [2]string{addressInfo.ServerID, address})
	}

	return addressList, nil
}