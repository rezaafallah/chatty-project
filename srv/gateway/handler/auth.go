package handler

import (
	"my-project/pkg/logic"
	"my-project/srv/gateway/response"
	"my-project/srv/gateway/dto"
	"my-project/types"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Logic *logic.AuthLogic
}

func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Bind to DTO
	var req dto.RegisterReq 
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// 2. Convert DTO to Domain Model 
	domainReq := types.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	}

	mnemonic, err := h.Logic.Register(c.Request.Context(), domainReq)
	if err != nil {
		response.Error(c, 500, "Registration failed")
		return
	}
	response.Success(c, gin.H{"mnemonic": mnemonic})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request body")
		return
	}

	token, err := h.Logic.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, 401, "Login failed: "+err.Error())
		return
	}

	response.Success(c, dto.LoginRes{
		Token:     token,
		ExpiresIn: 72 * 3600,
	})
}

func (h *AuthHandler) RecoverAccount(c *gin.Context) {
	var req dto.RecoverReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid mnemonic format")
		return
	}

	token, err := h.Logic.RecoverAccount(c.Request.Context(), req.Mnemonic)
	if err != nil {
		response.Error(c, 401, "Recovery failed: invalid mnemonic")
		return
	}

	response.Success(c, dto.LoginRes{
		Token:     token,
		ExpiresIn: 72 * 3600,
	})
}