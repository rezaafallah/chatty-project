package srv

import (
	"github.com/gin-gonic/gin"
	"my-project/internal/adapter/redis"
	"my-project/types"
	// ... imports
)

type GatewayServer struct {
	Engine *gin.Engine
	Redis  *redis.Client
}

func NewGateway(r *redis.Client) *GatewayServer {
	g := &GatewayServer{
		Engine: gin.Default(),
		Redis:  r,
	}
	
	// Routes
	g.Engine.GET("/ws", g.wsHandler)
	return g
}

func (s *GatewayServer) wsHandler(c *gin.Context) {
	// 1. Upgrade to WS
	// 2. Read Message
	// 3. Push to Redis Queue (Do NOT process logic here)
	// s.Redis.RPush("chat.queue", msgBytes)
}