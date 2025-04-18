package postgres

import (
	"user_service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/flashhhhh/pkg/logging"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	logging.LogMessage("user_service", "Connecting to Postgres Database...", "INFO")
	logging.LogMessage("user_service", "DSN = " + dsn, "DEBUG")
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.LogMessage("user_service", "Failed to connect to Postgres Database: "+err.Error(), "ERROR")
		return nil, err
	}

	logging.LogMessage("user_service", "Connected to Postgres Database", "INFO")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	logging.LogMessage("user_service", "Running migrations...", "INFO")

	// Migrate the schema
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		logging.LogMessage("user_service", "Failed to run migrations: "+err.Error(), "ERROR")
		return err
	}

	logging.LogMessage("user_service", "Migrations completed successfully", "INFO")
	return nil
}