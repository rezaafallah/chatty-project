package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"my-project/internal/adapter/redis"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	Hub      *Hub
	Redis    *redis.Client
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Log      *logrus.Logger
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Log.WithError(err).Error("WS Read Error")
			}
			break
		}

		// ساخت پیام برای صف Core
		inboundMsg := map[string]interface{}{
			"sender_id":   c.UserID,
			"receiver_id": "", // در MVP کلاینت باید بگوید به کی می‌فرستد، اینجا باید Parse شود
			"payload":     string(message), 
			"timestamp":   time.Now().Unix(),
		}
		
		// نکته: در یک سناریوی واقعی، کلاینت یک JSON می‌فرستد که فیلد `to` دارد.
		// اینجا فرض می‌کنیم کلاینت JSON معتبر می‌فرستد.
		
		bytes, _ := json.Marshal(inboundMsg)

		err = c.Redis.PushQueue(context.Background(), consts.QueueChatInbound, bytes)
		if err != nil {
			c.Log.Error("Failed to push to redis queue", err)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}