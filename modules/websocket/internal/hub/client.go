package hub

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uuid.UUID
	username string
	roomID   string
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID, username, roomID string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
		roomID:   roomID,
	}
}

// ReadPump reads messages from the WebSocket connection
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, message, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				logrus.WithField("user_id", c.userID).Info("WebSocket closed normally")
			} else {
				logrus.WithError(err).WithField("user_id", c.userID).Error("Error reading from WebSocket")
			}
			break
		}

		// Forward message to hub for broadcasting
		c.hub.Broadcast <- &BroadcastMessage{
			roomID:  c.roomID,
			message: message,
			sender:  c,
		}
	}
}

// WritePump sends messages to the WebSocket connection
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// Hub closed the channel
				c.conn.Close(websocket.StatusGoingAway, "")
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()

			if err != nil {
				logrus.WithError(err).WithField("user_id", c.userID).Error("Error writing to WebSocket")
				return
			}

		case <-ticker.C:
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(writeCtx)
			cancel()

			if err != nil {
				logrus.WithError(err).WithField("user_id", c.userID).Error("Error sending ping")
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
