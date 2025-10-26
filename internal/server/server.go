package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"egaldeutsch-be/internal/auth"
	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/database"
	"egaldeutsch-be/internal/middleware"
	authmodule "egaldeutsch-be/modules/auth"
	"egaldeutsch-be/modules/quiz"
	"egaldeutsch-be/modules/user"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	db     *database.Database
}

// NewServer creates a new server instance with all dependencies properly initialized.
// It follows Go philosophy by returning errors instead of calling fatal, allowing
// callers to decide how to handle failures.
func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Configure Gin mode based on environment
	configureGinMode(cfg.Server.Host)

	// Initialize database connection
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize modules with dependency injection
	authRepo := authmodule.NewRepository(db.DB)
	authService := auth.NewService(cfg.Jwt, authRepo)
	userModule := user.NewModule(db.DB, cfg.Jwt)
	authModule := authmodule.NewModule(authService, userModule.Service, cfg.Jwt)
	quizModule := quiz.NewModule(db.DB)

	// Run database migrations
	if err := runMigrations(db, userModule, authModule); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Setup HTTP router
	router := createRouter(cfg.Jwt, userModule, authModule, quizModule)

	return &Server{
		config: cfg,
		router: router,
		db:     db,
	}, nil
}

// configureGinMode sets Gin mode based on the host configuration.
func configureGinMode(host string) {
	gin.SetMode(gin.ReleaseMode)
	if host == "localhost" {
		gin.SetMode(gin.DebugMode)
	}
}

// runMigrations performs database migrations for all modules.
func runMigrations(db *database.Database, userModule *user.Module, authModule *authmodule.Module) error {
	modelsToMigrate := append(
		userModule.GetModelsForMigration(),
		authmodule.GetModelsForMigration()...,
	)

	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	return nil
}

// createRouter sets up the HTTP router with all routes and middleware.
func createRouter(jwtCfg config.JwtConfig, userModule *user.Module, authModule *authmodule.Module, quizModule *quiz.Module) *gin.Engine {
	router := gin.New()

	// Add middleware in correct order
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", healthCheckHandler)

	// API routes
	api := router.Group("/api/v1")
	{
		userModule.RegisterRoutes(api, jwtCfg)
		authModule.RegisterRoutes(api, jwtCfg)
		quizModule.RegisterRoutes(api)

	}

	return router
}

// healthCheckHandler provides a simple health check endpoint.
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "egaldeutsch-be",
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	logrus.Infof("Starting server on %s", addr)

	return s.router.Run(addr)
}
