package repository_test

import (
	"errors"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMocks() (*gorm.DB, sqlmock.Sqlmock, *redis.Client, redismock.ClientMock, *es.Client, error) {
	// Setup SQL mock
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Setup Redis mock
	redisCli, redisMock := redismock.NewClientMock()

	// Setup Elasticsearch mock (simplified)
	esCfg := es.Config{}
	esClient, _ := es.NewClient(esCfg)

	return db, sqlMock, redisCli, redisMock, esClient, nil
}

func TestCreateServer(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	// Test creating a server with "On" status
	t.Run("Create server with On status", func(t *testing.T) {
		server := &domain.Server{
			ServerID:   "srv-001",
			ServerName: "Test Server",
			Status:     "On",
			IPv4:       "192.168.1.1",
			Port:       8080,
		}

		// SQL mock expectations
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "servers"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Redis mock expectations
		redisMock.ExpectSetBit("server_status", 1, 1).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		id, err := repo.CreateServer(server)

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test creating a server with "Off" status
	t.Run("Create server with Off status", func(t *testing.T) {
		server := &domain.Server{
			ServerID:   "srv-002",
			ServerName: "Test Server 2",
			Status:     "Off",
			IPv4:       "192.168.1.2",
			Port:       8081,
		}

		// SQL mock expectations
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "servers"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectCommit()

		// Redis mock expectations
		redisMock.ExpectSetBit("server_status", 2, 0).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		id, err := repo.CreateServer(server)

		assert.NoError(t, err)
		assert.Equal(t, 2, id)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test with database error
	t.Run("Database error", func(t *testing.T) {
		server := &domain.Server{
			ServerID:   "srv-003",
			ServerName: "Test Server 3",
			Status:     "On",
			IPv4:       "192.168.1.3",
			Port:       8082,
		}

		// SQL mock expectations - simulate error
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "servers"`).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		repo := repository.NewServerRepository(db, redisCli, esClient)
		id, err := repo.CreateServer(server)

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestViewServers(t *testing.T) {
	db, mock, redisCli, _, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("View servers with filters", func(t *testing.T) {
		filter := &dto.ServerFilter{
			ServerID:   "srv-001",
			ServerName: "Test",
			Status:     "On",
			IPv4:       "192.168.1.1",
			Port:       8080,
		}

		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, "srv-001", "Test Server", "On", "192.168.1.1", 8080)

		mock.ExpectQuery(`SELECT .+ FROM "servers" WHERE .+ ORDER BY .+`).
			WillReturnRows(rows)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		servers, err := repo.ViewServers(filter, 1, 10, "id", "asc")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(servers))
		assert.Equal(t, "srv-001", servers[0].ServerID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No filters", func(t *testing.T) {
		filter := &dto.ServerFilter{}

		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, "srv-001", "Server 1", "On", "192.168.1.1", 8080).
			AddRow(2, "srv-002", "Server 2", "Off", "192.168.1.2", 8081)

		mock.ExpectQuery(`SELECT .+ FROM "servers" WHERE .+ ORDER BY .+`).
			WillReturnRows(rows)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		servers, err := repo.ViewServers(filter, 1, 10, "id", "asc")

		assert.NoError(t, err)
		assert.Equal(t, 2, len(servers))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		filter := &dto.ServerFilter{}

		mock.ExpectQuery(`SELECT .+ FROM "servers" .+`).
			WillReturnError(errors.New("database error"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		servers, err := repo.ViewServers(filter, 1, 10, "id", "asc")

		assert.Error(t, err)
		assert.Nil(t, servers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateServers(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	// Test successful insertion of multiple servers
	t.Run("Successfully insert multiple servers", func(t *testing.T) {
		servers := []domain.Server{
			{ServerID: "srv-001", ServerName: "Server 1", Status: "On", IPv4: "192.168.1.1", Port: 8080},
			{ServerID: "srv-002", ServerName: "Server 2", Status: "Off", IPv4: "192.168.1.2", Port: 8081},
		}

		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, "srv-001", "Server 1", "On", "192.168.1.1", 8080).
			AddRow(2, "srv-002", "Server 2", "Off", "192.168.1.2", 8081)

		mock.ExpectQuery(`INSERT INTO servers \(server_id, server_name, status, ipv4, port\) VALUES`).
			WillReturnRows(rows)

		redisMock.ExpectSetBit("server_status", 1, 1).SetVal(0)
		redisMock.ExpectSetBit("server_status", 2, 0).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		inserted, nonInserted, err := repo.CreateServers(servers)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(inserted))
		assert.Equal(t, 0, len(nonInserted))
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test partial insertion (some servers succeed, some fail)
	t.Run("Partial insertion of servers", func(t *testing.T) {
		servers := []domain.Server{
			{ServerID: "srv-003", ServerName: "Server 3", Status: "On", IPv4: "192.168.1.3", Port: 8083},
			{ServerID: "srv-004", ServerName: "Server 4", Status: "Off", IPv4: "192.168.1.4", Port: 8084},
		}

		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(3, "srv-003", "Server 3", "On", "192.168.1.3", 8083)

		mock.ExpectQuery(`INSERT INTO servers \(server_id, server_name, status, ipv4, port\) VALUES`).
			WillReturnRows(rows)

		redisMock.ExpectSetBit("server_status", 3, 1).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		inserted, nonInserted, err := repo.CreateServers(servers)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(inserted))
		assert.Equal(t, 1, len(nonInserted))
		assert.Equal(t, "srv-003", inserted[0].ServerID)
		assert.Equal(t, "srv-004", nonInserted[0].ServerID)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test database error
	t.Run("Database error during insertion", func(t *testing.T) {
		servers := []domain.Server{
			{ServerID: "srv-005", ServerName: "Server 5", Status: "On", IPv4: "192.168.1.5", Port: 8085},
		}

		mock.ExpectQuery(`INSERT INTO servers \(server_id, server_name, status, ipv4, port\) VALUES`).
			WillReturnError(errors.New("database error"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		inserted, nonInserted, err := repo.CreateServers(servers)

		assert.Error(t, err)
		assert.Nil(t, inserted)
		assert.Nil(t, nonInserted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateServer(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	// Test updating a server with "On" status
	t.Run("Update server with On status", func(t *testing.T) {
		serverID := "srv-001"
		updatedData := map[string]interface{}{
			"server_name": "Updated Server",
			"status":      "On",
			"ipv4":        "192.168.1.100",
		}

		// SQL mock expectations for update
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// SQL mock expectations for getting server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, serverID, "Updated Server", "On", "192.168.1.100", 8080)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)

		// Redis mock expectations
		redisMock.ExpectSetBit("server_status", 1, 1).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test updating a server with "Off" status
	t.Run("Update server with Off status", func(t *testing.T) {
		serverID := "srv-002"
		updatedData := map[string]interface{}{
			"server_name": "Updated Server 2",
			"status":      "Off",
		}

		// SQL mock expectations for update
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// SQL mock expectations for getting server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(2, serverID, "Updated Server 2", "Off", "192.168.1.2", 8081)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)

		// Redis mock expectations
		redisMock.ExpectSetBit("server_status", 2, 0).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test with database update error
	t.Run("Database update error", func(t *testing.T) {
		serverID := "srv-003"
		updatedData := map[string]interface{}{
			"status": "On",
		}

		// SQL mock expectations - simulate error
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnError(errors.New("database update error"))
		mock.ExpectRollback()

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.Error(t, err)
		assert.Equal(t, "database update error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test with error getting server ID
	t.Run("Error getting server ID", func(t *testing.T) {
		serverID := "srv-004"
		updatedData := map[string]interface{}{
			"status": "On",
		}

		// SQL mock expectations for update
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// SQL mock expectations for getting server ID - simulate error
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnError(errors.New("server not found"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.Error(t, err)
		assert.Equal(t, "server not found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test with Redis error
	t.Run("Redis error", func(t *testing.T) {
		serverID := "srv-005"
		updatedData := map[string]interface{}{
			"status": "On",
		}

		// SQL mock expectations for update
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// SQL mock expectations for getting server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(5, serverID, "Some Server", "On", "192.168.1.5", 8080)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)

		// Redis mock expectations - simulate error
		redisMock.ExpectSetBit("server_status", 5, 1).SetErr(errors.New("redis connection error"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.Error(t, err)
		assert.Equal(t, "redis connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	// Test updating without status field
	t.Run("Update without status field", func(t *testing.T) {
		serverID := "srv-006"
		updatedData := map[string]interface{}{
			"server_name": "Just Name Updated",
			"ipv4":        "192.168.1.200",
		}

		// SQL mock expectations for update
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// SQL mock expectations for getting server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(6, serverID, "Just Name Updated", "Off", "192.168.1.200", 8080)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)

		// Redis mock expectations - should still update to 0 since status is not "On"
		redisMock.ExpectSetBit("server_status", 6, 0).SetVal(0)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServer(serverID, updatedData)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestDeleteServer(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	// Test successful deletion
	t.Run("Successfully delete server", func(t *testing.T) {
		serverID := "srv-001"
		
		// SQL mock expectations for getting the server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, serverID, "Test Server", "On", "192.168.1.1", 8080)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)
		
		// SQL mock expectations for deleting
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "servers" WHERE`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		
		// Redis mock expectations
		redisMock.ExpectSetBit("server_status", 1, 0).SetVal(0)
		
		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.DeleteServer(serverID)
		
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
	
	// Test error when getting server ID
	t.Run("Server not found", func(t *testing.T) {
		serverID := "srv-not-exist"
		
		// SQL mock expectations - server not found
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnError(gorm.ErrRecordNotFound)
		
		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.DeleteServer(serverID)
		
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	// Test error when deleting server
	t.Run("Error when deleting server", func(t *testing.T) {
		serverID := "srv-002"
		
		// SQL mock expectations for getting the server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(2, serverID, "Test Server 2", "Off", "192.168.1.2", 8081)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)
		
		// SQL mock expectations for deleting - simulate error
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "servers" WHERE`).
			WillReturnError(errors.New("deletion error"))
		mock.ExpectRollback()
		
		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.DeleteServer(serverID)
		
		assert.Error(t, err)
		assert.Equal(t, "deletion error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	// // Test error when updating Redis
	t.Run("Redis error", func(t *testing.T) {
		serverID := "srv-003"
		
		// SQL mock expectations for getting the server ID
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(3, serverID, "Test Server 3", "On", "192.168.1.3", 8082)
		mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 ORDER BY "servers"."server_id" LIMIT \$2`).
			WithArgs(serverID, 1).
			WillReturnRows(rows)
		
		// SQL mock expectations for deleting
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "servers" WHERE`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		
		// Redis mock expectations - simulate error
		redisMock.ExpectSetBit("server_status", 3, 0).SetErr(errors.New("redis connection error"))
		
		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.DeleteServer(serverID)
		
		assert.Error(t, err)
		assert.Equal(t, "redis connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestUpdateServerStatus(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("Update status from Off to On", func(t *testing.T) {
		serverID := 1
		newStatus := "On"

		// Redis expectations
		redisMock.ExpectGetBit("server_status", int64(serverID)).SetVal(0) // Current status is Off (0)
		redisMock.ExpectSetBit("server_status", int64(serverID), 1).SetVal(0)

		// SQL expectations
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "servers" SET "status"=\$1,"last_updated"=\$2 WHERE id = \$3`).
			WithArgs(newStatus, sqlmock.AnyArg(), serverID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServerStatus(serverID, newStatus)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("No change in status", func(t *testing.T) {
		serverID := 1
		newStatus := "On"

		// Redis expectations - status is already On
		redisMock.ExpectGetBit("server_status", int64(serverID)).SetVal(1) // Current status is already On (1)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.UpdateServerStatus(serverID, newStatus)

		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestGetAllAddresses(t *testing.T) {
	db, mock, redisCli, _, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("Successfully get all addresses", func(t *testing.T) {
		// Define expected results
		expectedAddresses := []dto.ServerAddress{
			{ID: 1, IPv4: "192.168.1.1", Port: 8080},
			{ID: 2, IPv4: "192.168.1.2", Port: 8081},
			{ID: 3, IPv4: "192.168.1.3", Port: 8082},
		}

		// Create mock rows
		rows := sqlmock.NewRows([]string{"id", "ipv4", "port"})
		for _, addr := range expectedAddresses {
			rows.AddRow(addr.ID, addr.IPv4, addr.Port)
		}

		// Set up SQL mock expectations
		mock.ExpectQuery(`SELECT "id","ipv4","port" FROM "servers"`).
			WillReturnRows(rows)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		addresses, err := repo.GetAllAddresses()

		assert.NoError(t, err)
		assert.Equal(t, expectedAddresses, addresses)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		// Set up SQL mock to return an error
		mock.ExpectQuery(`SELECT "id","ipv4","port" FROM "servers"`).
			WillReturnError(errors.New("database error"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		addresses, err := repo.GetAllAddresses()

		assert.Error(t, err)
		assert.Nil(t, addresses)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result", func(t *testing.T) {
		// Set up SQL mock to return empty result
		rows := sqlmock.NewRows([]string{"id", "ipv4", "port"})

		mock.ExpectQuery(`SELECT "id","ipv4","port" FROM "servers"`).
			WillReturnRows(rows)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		addresses, err := repo.GetAllAddresses()

		assert.NoError(t, err)
		assert.Empty(t, addresses)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetNumOnServers(t *testing.T) {
	db, _, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("Get number of On servers", func(t *testing.T) {
		expectedCount := 5
		
		// Redis expectations
		redisMock.ExpectBitCount("server_status", nil).SetVal(int64(expectedCount))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		count, err := repo.GetNumOnServers()

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestGetNumServers(t *testing.T) {
	db, mock, redisCli, _, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("Get total number of servers", func(t *testing.T) {
		expectedCount := 10
		
		// SQL expectations - when db.Model(&domain.Server{}).Count(&count) is called
		mock.ExpectQuery(`SELECT count\(\*\) FROM "servers"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		count, err := repo.GetNumServers()

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		// SQL expectations - simulate error
		mock.ExpectQuery(`SELECT count\(\*\) FROM "servers"`).
			WillReturnError(errors.New("database error"))

		repo := repository.NewServerRepository(db, redisCli, esClient)
		count, err := repo.GetNumServers()

		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAddServerStatus(t *testing.T) {
	// This test is complex because it involves Elasticsearch
	// A simplified version is provided, but may need adjustments
	db, _, redisCli, _, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	// Create a mock implementation for testing AddServerStatus
	// In a real implementation, you might want to use a proper ES mock
	repo := repository.NewServerRepository(db, redisCli, esClient)
	
	// Note: This is more of an integration test since ES is hard to mock
	// In a real test environment, you might want to skip this test or use a test ES instance
	t.Run("Add server status - integration test", func(t *testing.T) {
		t.Skip("Skipping Elasticsearch integration test")
		err := repo.AddServerStatus(1, "On")
		assert.NoError(t, err)
	})
}

func TestSyncServerStatus(t *testing.T) {
	db, mock, redisCli, redisMock, esClient, err := setupMocks()
	if err != nil {
		t.Fatalf("Failed to setup mocks: %v", err)
	}

	t.Run("Sync server status", func(t *testing.T) {
		// Mock DB query to get servers
		rows := sqlmock.NewRows([]string{"id", "server_id", "server_name", "status", "ipv4", "port"}).
			AddRow(1, "srv-001", "Server 1", "On", "192.168.1.1", 8080).
			AddRow(2, "srv-002", "Server 2", "Off", "192.168.1.2", 8081).
			AddRow(3, "srv-003", "Server 3", "On", "192.168.1.3", 8082)

		mock.ExpectQuery(`SELECT (.+) FROM "servers"`).
			WillReturnRows(rows)

		// Redis expectations
		redisMock.ExpectDel("server_status").SetVal(1)
		redisMock.ExpectSetBit("server_status", 1, 1).SetVal(0)
		redisMock.ExpectSetBit("server_status", 3, 1).SetVal(0)
		redisMock.ExpectBitCount("server_status", nil).SetVal(2)

		repo := repository.NewServerRepository(db, redisCli, esClient)
		err := repo.SyncServerStatus()

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}
