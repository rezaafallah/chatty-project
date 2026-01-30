package ws

import (
	"context"
	"encoding/json"
	"time"

	"my-project/pkg/broker"
	"my-project/pkg/consts"
	"my-project/pkg/logger"
	"my-project/pkg/network"
	"my-project/pkg/uid"  
	"my-project/pkg/utils"
	"my-project/types"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	Hub       *Hub
	Broker    broker.MessageBroker
	Conn      network.SocketConnection
	Send      chan []byte
	UserID    string
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
			if network.IsUnexpectedCloseError(err, network.CloseGoingAway, network.CloseAbnormalClosure) {
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
		senderUUID, err := uid.Parse(c.UserID)
		if err != nil {
			c.Log.Error("Invalid Sender UUID")
			break
		}

		receiverUUID, err := uid.Parse(req.ReceiverID)
		if err != nil {
			c.Log.Errorf("Invalid Receiver UUID: %s", req.ReceiverID)
			continue
		}

		domainMsg := types.Message{
			SenderID:   senderUUID,
			ReceiverID: receiverUUID,
			Content:    req.Content,
			CreatedAt:  time.Now().Unix(),
		}

		bytes, err := json.Marshal(domainMsg)
		if err != nil {
			c.Log.WithError(err).Error("Failed to marshal domain message")
			continue
		}

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
				c.Conn.WriteMessage(network.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(network.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(network.PingMessage, nil); err != nil {
				return
			}
		}
	}
}