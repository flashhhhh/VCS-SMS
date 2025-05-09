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
	CreateServers(servers []domain.Server) ([]domain.Server, []domain.Server, error)
	ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error)
	UpdateServer(server_id string, updatedData map[string]interface{}) error
	DeleteServer(serverID string) error
	
	UpdateServerStatus(id int, status string) error
	GetAllAddresses() ([]dto.ServerAddress, error)

	AddServerStatus(id int, status string) error
	GetNumOnServers() (int, error)
	GetNumServers() (int, error)
	GetServerUptimeRatio(startTime, endTime time.Time) (float64, error)

	SyncServerStatus() error
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
func (r *serverRepository) CreateServers(servers []domain.Server) (inserted []domain.Server, nonInserted []domain.Server, err error) {
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

	query += " ON CONFLICT DO NOTHING RETURNING *"

	var result []domain.Server
	err = r.db.Raw(query).Scan(&result).Error
	if err != nil {
		logging.LogMessage("server_administration_service", "Error inserting servers: "+err.Error(), "ERROR")
		return nil, nil, err
	}

	// Determine non-inserted records
	insertedMap := make(map[string]bool)
	for _, server := range result {
		insertedMap[server.ServerID] = true
		inserted = append(inserted, server)
	}

	for _, server := range servers {
		if !insertedMap[server.ServerID] {
			nonInserted = append(nonInserted, server)
		}
	}

	// Update Redis bitmap for inserted records
	for _, server := range result {
		logging.LogMessage("server_administration_service", "Server "+strconv.Itoa(server.ID)+" inserted successfully", "INFO")

		status := 0
		if server.Status == "On" {
			status = 1
		}

		if err := r.redis.SetBit(context.Background(), "server_status", int64(server.ID), status).Err(); err != nil {
			logging.LogMessage("server_administration_service", "Error updating Redis bitmap for server ID: "+strconv.Itoa(server.ID)+", error: "+err.Error(), "ERROR")
		}
	}

	return inserted, nonInserted, nil
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

func (r *serverRepository) GetNumOnServers() (int, error) {
	numOnServers, _ := r.redis.BitCount(context.Background(), "server_status", nil).Result()
	return int(numOnServers), nil
}

func (r *serverRepository) GetNumServers() (int, error) {
	var count int64
	if err := r.db.Model(&domain.Server{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *serverRepository) GetServerUptimeRatio(startTime, endTime time.Time) (float64, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": startTime.Format(time.RFC3339),
					"lte": endTime.Format(time.RFC3339),
				},
			},
		},
		"aggs": map[string]interface{}{
			"per_server": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "id",
					"size":  10000,
				},
				"aggs": map[string]interface{}{
					"total_docs": map[string]interface{}{
						"value_count": map[string]interface{}{
							"field": "id",
						},
					},
					"on_count": map[string]interface{}{
						"filter": map[string]interface{}{
							"term": map[string]interface{}{
								"status.keyword": "On",
							},
						},
					},
					"on_ratio": map[string]interface{}{
						"bucket_script": map[string]interface{}{
							"buckets_path": map[string]interface{}{
								"on":  "on_count._count",
								"all": "total_docs.value",
							},
							"script": "params.all > 0 ? params.on / params.all : 0",
						},
					},
				},
			},
			"avg_ratio": map[string]interface{}{
				"avg_bucket": map[string]interface{}{
					"buckets_path": "per_server>on_ratio.value",
				},
			},
		},
	}

	// Encode the query
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("server_status"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
		r.es.Search.WithPretty(),
	)

	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("Error getting server uptime ratio: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	avgRatio := result["aggregations"].(map[string]interface{})["avg_ratio"].(map[string]interface{})["value"]
	if avgRatio == nil {
		return 0, nil
	}
	return avgRatio.(float64), nil
}

func (r *serverRepository) SyncServerStatus() error {
	// Get all server statuses from the database
	var servers []domain.Server
	if err := r.db.Find(&servers).Error; err != nil {
		return err
	}

	ctx := context.Background()

	// Delete old data
	r.redis.Del(ctx, "server_status")

	// Set all bits first
	for _, server := range servers {
		status := 0
		if server.Status == "On" {
			status = 1
		}

		if status == 1 {
			println("Server ID: " + strconv.Itoa(server.ID) + " is On")
			_, err := r.redis.SetBit(ctx, "server_status", int64(server.ID), status).Result()
			if err != nil {
				println("SetBit failed for server ID: " + strconv.Itoa(server.ID) + " error: " + err.Error())
			}
		}
	}

	numOnServers, _ := r.redis.BitCount(ctx, "server_status", nil).Result()
	println("Num On server: " + strconv.Itoa(int(numOnServers)))

	return nil
}