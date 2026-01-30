package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"my-project/pkg/logger"
	"my-project/internal/port"
	"my-project/pkg/utils"
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
	Broker port.MessageBroker
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	Log       logger.Logger
	Sanitizer *service.Sanitizer
}

type IncomingReq struct {
	ReceiverID string `json:"to"`
	Content    string `json:"content"`
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

		var req IncomingReq
		if err := json.Unmarshal(message, &req); err != nil {
			c.Log.Error("Invalid JSON format from client")
			continue
		}
			req.Content = c.Sanitizer.Clean(req.Content)

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

		domainMsg := types.Message{
			// ID:         uuid.New(), core
			SenderID:   senderUUID,
			ReceiverID: receiverUUID,
			Content:    req.Content,
			CreatedAt:  time.Now().Unix(),
		}

		bytes, _ := json.Marshal(domainMsg)

		err = c.Broker.PushQueue(context.Background(), consts.QueueChatInbound, bytes)
		if err != nil {
			c.Log.Error("Failed to push to Broker queue", err)
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