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
	// 1. Bind to DTO
	var req dto.RegisterReq 
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, 400, err.Error())
		return
	}

	// 2. Convert DTO to Domain Model 
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, 400, "Invalid request body")
		return
	}

	token, err := h.Logic.Login(req.Username, req.Password)
	if err != nil {
		api.Error(c, 401, "Login failed: "+err.Error())
		return
	}

	api.Success(c, dto.LoginRes{
		Token:     token,
		ExpiresIn: 72 * 3600,
	})
}