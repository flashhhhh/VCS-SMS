package repository

import (
	"server_administration_service/internal/domain"

	"gorm.io/gorm"
)

type ServerRepository interface {
	CreateServer(server *domain.Server) error
	UpdateServer(server *domain.Server) error
	DeleteServer(serverID string) error
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