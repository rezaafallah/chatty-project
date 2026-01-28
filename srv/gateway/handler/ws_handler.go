package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"my-project/srv/gateway/ws"
	"my-project/internal/adapter/redis"
	"my-project/pkg/logger"
)

type WSHandler struct {
	Hub   *ws.Hub
	Redis *redis.Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// در پروداکشن باید Origin چک شود
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WSHandler) HandleConnection(c *gin.Context) {
	// 1. دریافت UserID از Context (که میدل‌ور Auth ست کرده)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// 2. ارتقا به WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Setup().Error("Failed to upgrade ws", err)
		return
	}

	// 3. ایجاد Client و ثبت در Hub
	client := &ws.Client{
		Hub:    h.Hub,
		Redis:  h.Redis,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
		Log:    logger.Setup(),
	}

	client.Hub.Register <- client

	// 4. اجرای Pumpها در Goroutine
	go client.WritePump()
	go client.ReadPump()
}