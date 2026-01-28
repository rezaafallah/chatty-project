package ws

import (
	"sync"
	"github.com/sirupsen/logrus"
	"my-project/pkg/logger"
)

type Hub struct {
	// نقشه آنلاین‌ها: UserID -> Client
	// نکته: اگر کاربر چند دیوایس داشته باشد باید map[string][]*Client باشد. فعلاً ساده در نظر می‌گیریم.
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
	Log        *logrus.Logger
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Log:        logger.Setup(),
	}
}

func (h *Hub) Run() {
	h.Log.Info("WebSocket Hub Started")
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			// اگر کانکشن قبلی بود، قطعش کن (Single Device Policy)
			if old, ok := h.Clients[client.UserID]; ok {
				close(old.Send)
				delete(h.Clients, client.UserID)
			}
			h.Clients[client.UserID] = client
			h.Mutex.Unlock()
			h.Log.Infof("User connected: %s", client.UserID)

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.Mutex.Unlock()
			h.Log.Infof("User disconnected: %s", client.UserID)
		}
	}
}

// BroadcastToUser ارسال پیام به کاربر خاص (توسط Worker صدا زده می‌شود)
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	h.Mutex.RLock()
	defer h.Mutex.RUnlock()

	if client, ok := h.Clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			// اگر بافر پر بود، کلاینت را قطع کن (Bloated connection)
			close(client.Send)
			delete(h.Clients, userID)
		}
	} else {
		// کاربر آفلاین است. (می‌توانید اینجا لاگ کنید یا Push Notification بفرستید)
		// h.Log.Debugf("User %s is offline, message dropped (or push notify)", userID)
	}
}