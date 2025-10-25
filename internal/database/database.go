package database

import (
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"egaldeutsch-be/internal/config"
)

type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection using the provided configuration.
// It follows Go philosophy by being explicit about dependencies and failure modes.
func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	if err := validateDatabaseConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	db, err := openDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := configureConnectionPool(db, cfg); err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	logrus.Info("Successfully connected to database")
	return &Database{db}, nil
}

// validateDatabaseConfig ensures required database configuration is present.
func validateDatabaseConfig(cfg config.DatabaseConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if cfg.User == "" {
		return fmt.Errorf("database user is required")
	}
	return nil
}

// openDatabase establishes the database connection with proper logging configuration.
func openDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	connStr := buildConnectionString(cfg)

	gormConfig := &gorm.Config{
		Logger: createGormLogger(),
	}

	db, err := gorm.Open(postgres.Open(connStr), gormConfig)
	if err != nil {
		return nil, err
	}

	// Verify the connection is actually working
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// buildConnectionString creates PostgreSQL connection string from config.
func buildConnectionString(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

// createGormLogger creates a GORM logger that integrates with logrus.
func createGormLogger() gormlogger.Interface {
	return gormlogger.New(
		newLogrusWriter(logrus.StandardLogger().Out),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// configureConnectionPool sets up database connection pool with the provided settings.
func configureConnectionPool(db *gorm.DB, cfg config.DatabaseConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Apply connection pool settings with validation
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25 // sensible default
	}

	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 10 // sensible default
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	return nil
}

// newLogrusWriter adapts an io.Writer (logrus.Out) to gorm's logger writer expectations.
type logrusWriter struct {
	w io.Writer
}

func newLogrusWriter(w io.Writer) *logrusWriter {
	return &logrusWriter{w: w}
}

// Printf satisfies the logger writer interface used by gorm's logger.
func (lw *logrusWriter) Printf(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}
