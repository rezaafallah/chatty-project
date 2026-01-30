package ws

import (
	"context"
	"sync"
	"my-project/pkg/logger"
	"my-project/pkg/broker"
)

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
	Log        logger.Logger
	Broker     broker.MessageBroker
}

func NewHub(broker broker.MessageBroker) *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Log:        logger.Setup(),
		Broker:     broker,
	}
}

func (h *Hub) Run() {
	h.Log.Info("WebSocket Hub Started")
	ctx := context.Background()

	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			if old, ok := h.Clients[client.UserID]; ok {
				close(old.Send)
				delete(h.Clients, client.UserID)
			}
			h.Clients[client.UserID] = client
			h.Mutex.Unlock()
			
			h.Broker.SetUserOnline(ctx, client.UserID)
			h.Log.Infof("User connected: %s", client.UserID)

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.Mutex.Unlock()
			
			h.Broker.SetUserOffline(ctx, client.UserID)
			h.Log.Infof("User disconnected: %s", client.UserID)
		}
	}
}

// BroadcastToUser
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	h.Mutex.RLock()
	client, ok := h.Clients[userID]
	h.Mutex.RUnlock()

	if ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, userID)
		}
	}
}