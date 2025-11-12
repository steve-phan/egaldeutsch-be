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
	"egaldeutsch-be/internal/redis"
	authmodule "egaldeutsch-be/modules/auth"
	"egaldeutsch-be/modules/quiz"
	"egaldeutsch-be/modules/user"
	websocketmodule "egaldeutsch-be/modules/websocket"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	db     *database.Database
	redis  *redis.RedisClient
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

	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}
	fmt.Printf("Redis client initialized successfully : %+v\n", redisClient)

	// Initialize modules with dependency injection
	authRepo := authmodule.NewRepository(db.DB)
	authService := auth.NewService(cfg.Jwt, authRepo)
	userModule := user.NewModule(db.DB, cfg.Jwt)
	authModule := authmodule.NewModule(authService, userModule.Service, cfg.Jwt)
	quizModule := quiz.NewModule(db.DB)
	wsModule := websocketmodule.NewModule(redisClient)

	// Run database migrations
	if err := runMigrations(db, userModule, authModule, quizModule, wsModule); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Setup HTTP router
	router := createRouter(cfg.Jwt, userModule, authModule, quizModule, wsModule)

	return &Server{
		config: cfg,
		router: router,
		db:     db,
		redis:  redisClient,
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
func runMigrations(db *database.Database, userModule *user.Module, authModule *authmodule.Module, quizModule *quiz.Module, wsModule *websocketmodule.Module) error {
	// Start with the first module's models
	modelsToMigrate := userModule.GetModelsForMigration()
	// Append the other modules' models by unpacking them one at a time
	modelsToMigrate = append(modelsToMigrate, authModule.GetModelsForMigration()...)
	modelsToMigrate = append(modelsToMigrate, quizModule.GetModelsForMigration()...)
	modelsToMigrate = append(modelsToMigrate, wsModule.GetModelsForMigration()...)

	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	return nil
}

// createRouter sets up the HTTP router with all routes and middleware.
func createRouter(jwtCfg config.JwtConfig, userModule *user.Module, authModule *authmodule.Module, quizModule *quiz.Module, wsModule *websocketmodule.Module) *gin.Engine {
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
		wsModule.RegisterRoutes(api, jwtCfg)
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
