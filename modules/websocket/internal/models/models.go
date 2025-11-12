package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	MessageTypeChat     MessageType = "chat"
	MessageTypeJoin     MessageType = "join"
	MessageTypeLeave    MessageType = "leave"
	MessageTypeTyping   MessageType = "typing"
	MessageTypeError    MessageType = "error"
	MessageTypeRoomInfo MessageType = "room_info"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      MessageType `json:"type"`
	RoomID    string      `json:"room_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

// ChatMessage represents a stored chat message for persistence
type ChatMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TableName specifies the table name for the ChatMessage model
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// RoomInfo represents room information
type RoomInfo struct {
	RoomID    string    `json:"room_id"`
	UserCount int       `json:"user_count"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

// Room represents a chat room
type Room struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// TableName specifies the table name for the Room model
func (Room) TableName() string {
	return "chat_rooms"
}

// CreateRoomRequest represents the request to create a room
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// RoomIDParam represents the URI parameter for room ID
type RoomIDParam struct {
	RoomID string `uri:"room_id" binding:"required,min=1,max=100"`
}

// GetHistoryQuery represents query parameters for getting room history
type GetHistoryQuery struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}
