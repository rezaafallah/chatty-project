package gateway

import (
	"my-project/srv/gateway/handler"
	"my-project/srv/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(secret string, authH *handler.AuthHandler, wsH *handler.WSHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", authH.Register)
	}

	protected := v1.Group("/")
	protected.Use(middleware.Auth(secret))
	{
		
		protected.GET("/ws", wsH.HandleConnection)
	}
	return r
}