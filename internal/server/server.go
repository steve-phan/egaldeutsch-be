package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/database"
	"egaldeutsch-be/internal/middleware"
	"egaldeutsch-be/modules/user"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	db     *database.Database
}

func NewServer(cfg *config.Config) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	if cfg.Server.Host == "localhost" {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize user module
	userModule := user.NewModule(db.DB)

	// Auto-migrate models from all modules
	modelsToMigrate := userModule.GetModelsForMigration()
	// Add other modules here as they are implemented

	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "egaldeutsch-be",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Register module routes
		userModule.RegisterRoutes(api)

		// Article routes (TODO: Update to use services)
		// api.GET("/articles", articleHandler.GetArticles)
	}

	return &Server{
		config: cfg,
		router: router,
		db:     db,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	logrus.Infof("Starting server on %s", addr)

	return s.router.Run(addr)
}
