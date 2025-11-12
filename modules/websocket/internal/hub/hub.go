package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"egaldeutsch-be/internal/redis"
	"egaldeutsch-be/modules/websocket/internal/models"
)

type BroadcastMessage struct {
	roomID  string
	message []byte
	sender  *Client
}

type Hub struct {
	// Registered clients per room
	rooms map[string]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Broadcast messages to room
	Broadcast chan *BroadcastMessage

	// Redis client for message persistence
	redis *redis.RedisClient

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

func NewHub(redisClient *redis.RedisClient) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
		redis:      redisClient,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastToRoom(ctx, message)

		case <-ctx.Done():
			logrus.Info("Hub shutting down")
			return
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.roomID] == nil {
		h.rooms[client.roomID] = make(map[*Client]bool)
	}
	h.rooms[client.roomID][client] = true

	logrus.WithFields(logrus.Fields{
		"user_id":  client.userID,
		"username": client.username,
		"room_id":  client.roomID,
	}).Info("Client registered")

	// Send join message to room
	h.sendJoinMessage(client)

	// Send room info to new client
	h.sendRoomInfo(client)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.roomID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			// Remove room if empty
			if len(clients) == 0 {
				delete(h.rooms, client.roomID)
			}

			logrus.WithFields(logrus.Fields{
				"user_id":  client.userID,
				"username": client.username,
				"room_id":  client.roomID,
			}).Info("Client unregistered")

			// Send leave message to room
			h.sendLeaveMessage(client)
		}
	}
}

func (h *Hub) broadcastToRoom(ctx context.Context, msg *BroadcastMessage) {
	h.mu.RLock()
	clients := h.rooms[msg.roomID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	// Parse message
	var wsMsg models.WSMessage
	if err := json.Unmarshal(msg.message, &wsMsg); err != nil {
		logrus.WithError(err).Error("Failed to parse WebSocket message")
		return
	}

	// Set sender information
	wsMsg.UserID = msg.sender.userID
	wsMsg.Username = msg.sender.username
	wsMsg.RoomID = msg.roomID
	wsMsg.Timestamp = time.Now()

	// Store message in Redis if it's a chat message
	if wsMsg.Type == models.MessageTypeChat {
		h.storeMessageInRedis(ctx, &wsMsg)
	}

	// Marshal updated message
	messageBytes, err := json.Marshal(wsMsg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal WebSocket message")
		return
	}

	// Broadcast to all clients in room
	for client := range clients {
		select {
		case client.send <- messageBytes:
		default:
			// Client's send channel is full, close it
			close(client.send)
			delete(clients, client)
		}
	}
}

func (h *Hub) sendJoinMessage(client *Client) {
	msg := models.WSMessage{
		Type:      models.MessageTypeJoin,
		RoomID:    client.roomID,
		UserID:    client.userID,
		Username:  client.username,
		Content:   fmt.Sprintf("%s joined the room", client.username),
		Timestamp: time.Now(),
	}

	messageBytes, _ := json.Marshal(msg)
	h.Broadcast <- &BroadcastMessage{
		roomID:  client.roomID,
		message: messageBytes,
		sender:  client,
	}
}

func (h *Hub) sendLeaveMessage(client *Client) {
	msg := models.WSMessage{
		Type:      models.MessageTypeLeave,
		RoomID:    client.roomID,
		UserID:    client.userID,
		Username:  client.username,
		Content:   fmt.Sprintf("%s left the room", client.username),
		Timestamp: time.Now(),
	}

	messageBytes, _ := json.Marshal(msg)

	// Send directly to remaining clients
	h.mu.RLock()
	clients := h.rooms[client.roomID]
	h.mu.RUnlock()

	if clients != nil {
		for c := range clients {
			select {
			case c.send <- messageBytes:
			default:
			}
		}
	}
}

func (h *Hub) sendRoomInfo(client *Client) {
	h.mu.RLock()
	clients := h.rooms[client.roomID]
	h.mu.RUnlock()

	users := make([]string, 0, len(clients))
	for c := range clients {
		users = append(users, c.username)
	}

	roomInfo := models.RoomInfo{
		RoomID:    client.roomID,
		UserCount: len(clients),
		Users:     users,
	}

	msg := models.WSMessage{
		Type:      models.MessageTypeRoomInfo,
		RoomID:    client.roomID,
		Timestamp: time.Now(),
	}

	// Embed room info in content as JSON
	roomInfoBytes, _ := json.Marshal(roomInfo)
	msg.Content = string(roomInfoBytes)

	messageBytes, _ := json.Marshal(msg)

	select {
	case client.send <- messageBytes:
	default:
	}
}

func (h *Hub) storeMessageInRedis(ctx context.Context, msg *models.WSMessage) {
	// Store in Redis list with expiry (e.g., 24 hours)
	key := fmt.Sprintf("chat:room:%s", msg.RoomID)

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message for Redis")
		return
	}

	if err := h.redis.RPush(ctx, key, messageBytes).Err(); err != nil {
		logrus.WithError(err).Error("Failed to store message in Redis")
		return
	}

	// Set expiry of 24 hours
	h.redis.Expire(ctx, key, 24*time.Hour)
}

// GetRoomHistory retrieves chat history from Redis
func (h *Hub) GetRoomHistory(ctx context.Context, roomID string, limit int64) ([]models.WSMessage, error) {
	key := fmt.Sprintf("chat:room:%s", roomID)

	// Get last N messages (0 = oldest, -1 = newest)
	start := int64(0)
	stop := int64(-1)
	if limit > 0 {
		// Get last N messages
		count, err := h.redis.LLen(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if count > limit {
			start = count - limit
		}
	}

	messages, err := h.redis.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	result := make([]models.WSMessage, 0, len(messages))
	for _, msgStr := range messages {
		var msg models.WSMessage
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			logrus.WithError(err).Error("Failed to unmarshal message from Redis")
			continue
		}
		result = append(result, msg)
	}

	return result, nil
}

// GetRoomUserCount returns the number of users in a room
func (h *Hub) GetRoomUserCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}
