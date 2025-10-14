package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"egaldeutsch-be/internal/config"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}
	logrus.Info("Successfully connected to database")

	// Set optimized connection pool settings for production
	// Based on PostgreSQL best practices and Go performance guidelines
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)   // ~4x CPU cores, adjust based on load
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)   // Keep some idle connections ready
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Refresh connections every 5 minutes
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Close idle connections after 1 minute

	return &Database{db}, nil
}
