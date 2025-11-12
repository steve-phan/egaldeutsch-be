package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"

	"egaldeutsch-be/internal/config"
	"egaldeutsch-be/modules/websocket/internal/hub"
	"egaldeutsch-be/modules/websocket/internal/models"
)

type WSHandler struct {
	hub     *hub.Hub
	jwtCfg  config.JwtConfig
}

func NewWSHandler(hub *hub.Hub, jwtCfg config.JwtConfig) *WSHandler {
	return &WSHandler{
		hub:    hub,
		jwtCfg: jwtCfg,
	}
}

// HandleConnection upgrades HTTP connection to WebSocket and manages the client lifecycle
func (h *WSHandler) HandleConnection(c *gin.Context) {
	// Get room ID from URL parameter
	var params models.RoomIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get username from context (you might want to add this to your auth middleware)
	username := c.GetString("username")
	if username == "" {
		username = "Anonymous" // Fallback, you should get this from user service
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		Subprotocols: []string{"chat"},
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade to WebSocket")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Create client
	client := hub.NewClient(h.hub, conn, userID, username, params.RoomID)

	// Register client with hub
	h.hub.Register <- client

	// Create context for this connection
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	logrus.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": username,
		"room_id":  params.RoomID,
	}).Info("WebSocket connection established")

	// Start read and write pumps
	go client.WritePump(ctx)
	client.ReadPump(ctx)
}

// GetRoomHistory returns chat history for a room
func (h *WSHandler) GetRoomHistory(c *gin.Context) {
	// Get room ID from URL parameter
	var params models.RoomIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	var query models.GetHistoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default limit
	if query.Limit <= 0 {
		query.Limit = 50 // Default to last 50 messages
	}
	if query.Limit > 100 {
		query.Limit = 100 // Max 100 messages
	}

	// Get history from hub (which gets it from Redis)
	history, err := h.hub.GetRoomHistory(c.Request.Context(), params.RoomID, int64(query.Limit))
	if err != nil {
		logrus.WithError(err).WithField("room_id", params.RoomID).Error("Failed to get room history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": history,
		"count":    len(history),
	})
}

// GetRoomInfo returns information about a room
func (h *WSHandler) GetRoomInfo(c *gin.Context) {
	// Get room ID from URL parameter
	var params models.RoomIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get room info from hub
	userCount := h.hub.GetRoomUserCount(params.RoomID)

	c.JSON(http.StatusOK, gin.H{
		"room_id":    params.RoomID,
		"user_count": userCount,
	})
}

// CreateRoom creates a new chat room (optional feature)
func (h *WSHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Generate room ID (you could use a more sophisticated approach)
	roomID := uuid.New().String()

	room := &models.Room{
		ID:          roomID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
		IsActive:    true,
	}

	// TODO: Save room to database if you want persistent rooms
	// For now, rooms are ephemeral (exist only while users are connected)

	logrus.WithFields(logrus.Fields{
		"room_id":   roomID,
		"room_name": req.Name,
		"created_by": userID,
	}).Info("Chat room created")

	c.JSON(http.StatusCreated, room)
}

// ListRooms returns a list of active rooms (optional feature)
func (h *WSHandler) ListRooms(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// For now, return empty list since rooms are ephemeral
	// TODO: Implement room persistence if needed
	rooms := []models.Room{}

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
		"count": len(rooms),
	})
}