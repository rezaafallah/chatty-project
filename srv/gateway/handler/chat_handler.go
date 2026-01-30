package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"my-project/pkg/logic"
	"my-project/srv/gateway/response"
)

type ChatHandler struct {
	Logic *logic.ChatLogic
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	// user_id from jwt token
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, 401, "Invalid User ID")
		return
	}

	msgs, err := h.Logic.GetHistory(userID)
	if err != nil {
		response.Error(c, 500, "Failed to fetch history")
		return
	}

	response.Success(c, msgs)
}