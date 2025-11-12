package websocket

import (
	"context"

	"github.com/gin-gonic/gin"

	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/internal/redis"
	"egaldeutsch-be/modules/websocket/internal/handlers"
	"egaldeutsch-be/modules/websocket/internal/hub"
)

type Module struct {
	hub     *hub.Hub
	handler *handlers.WSHandler
}

func NewModule(redisClient *redis.RedisClient) *Module {
	// Create hub with Redis client
	h := hub.NewHub(redisClient)

	// Start hub in background
	go h.Run(context.Background())

	// Create handler with hub
	handler := handlers.NewWSHandler(h, config.JwtConfig{}) // TODO: Pass JWT config

	return &Module{
		hub:     h,
		handler: handler,
	}
}

// RegisterRoutes registers WebSocket routes
func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwtCfg config.JwtConfig) {
	// Update handler with JWT config
	m.handler = handlers.NewWSHandler(m.hub, jwtCfg)

	ws := rg.Group("/ws")
	{
		// Protected routes - require authentication
		ws.GET("/chat/:room_id", m.handler.HandleConnection)
		ws.GET("/chat/:room_id/history", m.handler.GetRoomHistory)
		ws.GET("/chat/:room_id/info", m.handler.GetRoomInfo)

		// Room management (optional)
		ws.POST("/rooms", m.handler.CreateRoom)
		ws.GET("/rooms", m.handler.ListRooms)
	}
}

// GetModelsForMigration returns models that need to be migrated (for future persistence)
func (m *Module) GetModelsForMigration() []interface{} {
	// TODO: Return chat room and message models when implementing persistence
	return []interface{}{}
}
