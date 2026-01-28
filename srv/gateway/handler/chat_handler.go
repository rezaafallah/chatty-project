package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"my-project/internal/core"
	"my-project/pkg/api"
)

type ChatHandler struct {
	Logic *core.ChatLogic
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	// user_id from jwt token
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		api.Error(c, 401, "Invalid User ID")
		return
	}

	msgs, err := h.Logic.GetHistory(userID)
	if err != nil {
		api.Error(c, 500, "Failed to fetch history")
		return
	}

	api.Success(c, msgs)
}