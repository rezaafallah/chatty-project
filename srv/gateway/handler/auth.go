package handler

import (
	"my-project/internal/core"
	"my-project/pkg/api"
	"my-project/srv/gateway/dto" // Import DTO package
	"my-project/types"           // Import Types for conversion
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Logic *core.AuthLogic
}

func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Bind to DTO (نه مدل دامین)
	var req dto.RegisterReq 
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, 400, err.Error())
		return
	}

	// 2. Convert DTO to Domain Model (اگه متد Logic ورودی DTO نمیگیره)
	// یا متد Logic رو تغییر بده که DTO نگیره، بلکه پارامترهای ساده بگیره
	domainReq := types.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	}

	mnemonic, err := h.Logic.Register(domainReq)
	if err != nil {
		api.Error(c, 500, "Registration failed")
		return
	}
	api.Success(c, gin.H{"mnemonic": mnemonic})
}