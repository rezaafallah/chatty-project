package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"my-project/srv/gateway/ws"
	"my-project/internal/adapter/redis"
	"my-project/internal/service"
	"my-project/pkg/logger"
)

type WSHandler struct {
	Hub       *ws.Hub
	Redis     *redis.Client
	Sanitizer *service.Sanitizer
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WSHandler) HandleConnection(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Setup().Error("Failed to upgrade ws", err)
		return
	}

	client := &ws.Client{
		Hub:       h.Hub,
		Redis:     h.Redis,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		UserID:    userID,
		Log:       logger.Setup(), 
		Sanitizer: h.Sanitizer, 
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}