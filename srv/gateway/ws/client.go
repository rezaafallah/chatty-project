package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"my-project/internal/adapter/redis"
	"my-project/pkg/consts"
	"my-project/types"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	Hub    *Hub
	Redis  *redis.Client
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	Log    *logrus.Logger
}

// مدل ورودی: چیزی که فرانت (Postman) می‌فرستد
type IncomingReq struct {
	ReceiverID string `json:"to"`      // ID گیرنده
	Content    string `json:"content"` // متن پیام
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

		// 1. خواندن JSON ورودی
		var req IncomingReq
		if err := json.Unmarshal(message, &req); err != nil {
			c.Log.Error("Invalid JSON format from client")
			continue
		}

		// 2. اعتبارسنجی UUIDها
		senderUUID, err := uuid.Parse(c.UserID)
		if err != nil {
			c.Log.Error("Invalid Sender UUID")
			break
		}
		
		receiverUUID, err := uuid.Parse(req.ReceiverID)
		if err != nil {
			c.Log.Errorf("Invalid Receiver UUID: %s", req.ReceiverID)
			continue
		}

		// 3. ساخت پیام استاندارد
		domainMsg := types.Message{
			ID:         uuid.New(),
			SenderID:   senderUUID,
			ReceiverID: receiverUUID, // اینجا گیرنده درست ست می‌شود
			Content:    req.Content,
			CreatedAt:  time.Now().Unix(),
		}

		// 4. تبدیل به JSON و ارسال به صف
		bytes, _ := json.Marshal(domainMsg)

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