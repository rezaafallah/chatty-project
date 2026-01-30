package handler

import (
	"github.com/gin-gonic/gin"
	"my-project/pkg/broker"
	"my-project/pkg/logger"
	"my-project/pkg/network" // استفاده از Wrapper شبکه
	service "my-project/pkg/utils"
	"my-project/srv/gateway/response" // استفاده از پاسخ استاندارد
	"my-project/srv/gateway/ws"
)

type WSHandler struct {
	Hub       *ws.Hub
	Broker    broker.MessageBroker
	Sanitizer *service.Sanitizer
	Log       logger.Logger // اضافه کردن Logger برای جلوگیری از Setup تکراری
}

// استفاده از Wrapper برای تنظیمات Upgrader
var upgrader = network.NewUpgrader(1024, 1024)

func (h *WSHandler) HandleConnection(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	// استفاده از متد Upgrade داخل Wrapper
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.Log.Error("Failed to upgrade ws", err)
		return
	}

	client := &ws.Client{
		Hub:       h.Hub,
		Broker:    h.Broker,
		Conn:      conn, // تایپ این الان network.SocketConnection است
		Send:      make(chan []byte, 256),
		UserID:    userID,
		Log:       h.Log, // استفاده از لاگر تزریق شده
		Sanitizer: h.Sanitizer,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}