package main

import (
	"context"
	"log"
	"os"
	"time"

	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/core"
	"my-project/internal/repository"
	"my-project/internal/service"    
	"my-project/pkg/auth"            
	"my-project/srv/gateway"
	"my-project/srv/gateway/handler"
	"my-project/srv/gateway/worker"
	"my-project/srv/gateway/ws"
)

func main() {
	// 1. Infra
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.New(os.Getenv("REDIS_ADDR"))

	// 2. Dependencies (SRP Refactoring)
	userRepo := repository.NewUserRepository(db.Conn)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenMgr := auth.NewJWTManager(jwtSecret, 72*time.Hour)
	sanitizer := service.NewSanitizer()

	// 3. Logic
	msgRepo := repository.NewMessageRepository(db.Conn)
	authLogic := core.NewAuthLogic(userRepo, tokenMgr)
	chatLogic := core.NewChatLogic(msgRepo, rdb)

	// 4. WebSocket Hub
	hub := ws.NewHub()
	go hub.Run() 

	// 5. Redis Subscriber
	sub := worker.NewSubscriber(rdb.RDB, hub) 
	go sub.Start(context.Background())

	// 6. Handlers
	authHandler := &handler.AuthHandler{Logic: authLogic}
	// Sanitizer
	wsHandler := &handler.WSHandler{Hub: hub, Redis: rdb, Sanitizer: sanitizer}
	chatHandler := &handler.ChatHandler{Logic: chatLogic}

	// 7. Router 
	r := gateway.SetupRouter(jwtSecret, authHandler, wsHandler, chatHandler)
	
	log.Println("Gateway running on :8080")
	r.Run(":8080")
}