package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/flashhhhh/pkg/logging"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServerRepository interface {
	CreateServer(server *domain.Server) (int, error)
	CreateServers(servers []domain.Server) error
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
	DeleteServer(serverID string) error
	
	UpdateServerStatus(id int, status string) error
	GetAllAddresses() ([]dto.ServerAddress, error)

	AddServerStatus(id int, status string) error
}

type serverRepository struct {
	db *gorm.DB
	redis *redis.Client
	es *elasticsearch.Client
}

func NewServerRepository(db *gorm.DB, redis *redis.Client, es *elasticsearch.Client) ServerRepository {
	return &serverRepository{
		db: db,
		redis: redis,
		es: es,
	}
}

func (r *serverRepository) CreateServer(server *domain.Server) (int, error) {
	err := r.db.Create(server).Error
	if err != nil {
		return 0, err
	}

	// Turn the bit id in Redis bitmap 
	status := 0
	if server.Status == "On" {
		status = 1
	}

	/*
		WARNING: We are thinking how to make sure Redis must be updated after creating a server
		To make sure Redis is synced with the database
	*/
	r.redis.SetBit(context.Background(), "server_status", int64(server.ID), status).Err()

	return server.ID, nil
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

/*
	UPDATE: Using Raw query to receive just the inserted records
*/
func (r *serverRepository) CreateServers(servers []domain.Server) error {
	// Use raw SQL to insert multiple rows and get the inserted IDs

	query := `
		INSERT INTO servers (server_id, server_name, status, ipv4, port) VALUES 
	`

	for i, server := range servers {
		query += fmt.Sprintf("('%s', '%s', '%s', '%s', %d)",
			server.ServerID, server.ServerName, server.Status, server.IPv4, server.Port)
		
		if i < len(servers)-1 {
			query += ", "
		}
	}

	query += " ON CONFLICT DO NOTHING RETURNING *";

	var result []domain.Server
	err := r.db.Raw(query).Scan(&result).Error
	if err != nil {
		logging.LogMessage("server_administration_service", "Error inserting servers: "+err.Error(), "ERROR")
		return err
	}

	for _, result := range result {
		logging.LogMessage("server_administration_service", "Server " + strconv.Itoa(result.ID) + " inserted successfully", "INFO")

		logInformation := strconv.Itoa(result.ID) + 
						" " + result.ServerID +
						" " + result.ServerName +
						" " + result.Status +
						" " + result.IPv4 +
						" " + strconv.Itoa(result.Port)
		logging.LogMessage("server_administration_service", logInformation, "DEBUG") 

		status := 0
		if result.Status == "On" {
			status = 1
		}

		logging.LogMessage("server_administration_service", "Redis bitmap updated for server ID: "+strconv.Itoa(result.ID), "INFO")
		logging.LogMessage("server_administration_service", "Server is " + result.Status, "DEBUG")

		r.redis.SetBit(context.Background(), "server_status", int64(result.ID), status).Err()
	}

	return nil
}

func (r *serverRepository) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	if err := r.db.Model(&domain.Server{}).
		Where("server_id = ?", serverID).
		Updates(updatedData).Error; err != nil {
			return err
	}

	// Get the server's id
	var server domain.Server
	if err := r.db.Where("server_id = ?", serverID).First(&server).Error; err != nil {
		return err
	}

	// Update Redis bitmap
	statusValue := 0
	if updatedData["status"] == "On" {
		statusValue = 1
	}

	if err := r.redis.SetBit(context.Background(), "server_status", int64(server.ID), statusValue).Err(); err != nil {
		return err
	}

	return nil
}

func (r *serverRepository) DeleteServer(serverID string) error {
	// Get the server's ID before deleting
	var server domain.Server
	if err := r.db.Where("server_id = ?", serverID).First(&server).Error; err != nil {
		return err
	}

	// Delete the server
	if err := r.db.Delete(&domain.Server{}, serverID).Error; err != nil {
		return err
	}

	// Update Redis bitmap
	if err := r.redis.SetBit(context.Background(), "server_status", int64(server.ID), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r *serverRepository) UpdateServerStatus(id int, status string) error {
	/*
		WARNING: Not handling the case when Redis is crashed
	*/
	currentStatus := int(r.redis.GetBit(context.Background(), "server_status", int64(id)).Val())
	nextStatus := 0
	if status == "On" {
		nextStatus = 1
	}

	if currentStatus == nextStatus {
		logging.LogMessage("server_administration_service", "Server ID "+strconv.Itoa(id)+" status is already "+status, "INFO")
		return nil
	}

	// Update Redis bitmap
	r.redis.SetBit(context.Background(), "server_status", int64(id), nextStatus)

	if err := r.db.Model(&domain.Server{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) GetAllAddresses() ([]dto.ServerAddress, error) {
	var addresses []dto.ServerAddress
	if err := r.db.Model(&domain.Server{}).
		Select("id", "ipv4", "port").
		Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *serverRepository) AddServerStatus(id int, status string) error {
	// Add server status to Elasticsearch
	ctx := context.Background()

	doc := map[string]interface{}{
		"id":     id,
		"status": status,
		"timestamp": time.Now(),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return err
	}

	res, err := r.es.Index(
		"server_status",
		&buf,
		r.es.Index.WithContext(ctx),
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error indexing document: %s", res.String())
	}
	return nil
}