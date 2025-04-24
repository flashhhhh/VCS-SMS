package repository

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"

	"gorm.io/gorm"
)

type ServerRepository interface {
	CreateServer(server *domain.Server) error
	UpdateServer(server *domain.Server) error
	DeleteServer(serverID string) error
	
	UpdateServerStatus(serverID string, status string) error
	GetAllAddresses() ([]dto.ServerAddress, error)
}

type serverRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{
		db: db,
	}
}

func (r *serverRepository) CreateServer(server *domain.Server) error {
	if err := r.db.Create(server).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) UpdateServer(server *domain.Server) error {
	if err := r.db.Save(server).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) DeleteServer(serverID string) error {
	if err := r.db.Delete(&domain.Server{}, serverID).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) UpdateServerStatus(serverID string, status string) error {
	if err := r.db.Model(&domain.Server{}).
		Where("server_id = ?", serverID).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) GetAllAddresses() ([]dto.ServerAddress, error) {
	var addresses []dto.ServerAddress
	if err := r.db.Model(&domain.Server{}).
		Select("server_id", "ipv4", "port").
		Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}