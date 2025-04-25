package repository

import (
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServerRepository interface {
	CreateServer(server *domain.Server) error
	CreateServers(servers []domain.Server) error
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
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

func (r *serverRepository) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error) {
	var servers []domain.Server
	query := r.db.Model(&domain.Server{})

	if serverFilter.ServerID != "" {
		query = query.Where("server_id = ?", serverFilter.ServerID)
	}

	if serverFilter.ServerName != "" {	
		query = query.Where("server_name LIKE ?", "%"+serverFilter.ServerName+"%")
	}

	if serverFilter.Status != "" {
		query = query.Where("status = ?", serverFilter.Status)
	}

	if serverFilter.IPv4 != "" {
		query = query.Where("ipv4 = ?", serverFilter.IPv4)
	}

	if serverFilter.Port >= 0 {
		query = query.Where("port = ?", serverFilter.Port)
	}

	query = query.Where("id BETWEEN ? AND ?", from, to)

	// sortedColumn is mandatory
	err := query.Order(sortedColumn + " " + order).Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *serverRepository) CreateServers(servers []domain.Server) error {
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&servers).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	if err := r.db.Model(&domain.Server{}).
		Where("server_id = ?", serverID).
		Updates(updatedData).Error; err != nil {
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