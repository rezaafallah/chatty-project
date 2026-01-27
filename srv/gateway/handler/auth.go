package handler

import (
	"my-project/internal/core"
	"my-project/pkg/api"
	"my-project/types"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Logic *core.AuthLogic
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Error(c, 400, err.Error())
		return
	}
	mnemonic, err := h.Logic.Register(req)
	if err != nil {
		api.Error(c, 500, "Registration failed")
		return
	}
	api.Success(c, gin.H{"mnemonic": mnemonic})
}