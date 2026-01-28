package main

import (
	"context"
	"log"
	"os"

	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/core"
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

	// 2. Logic
	// FIX: Pass JWT Secret to NewAuthLogic
	jwtSecret := os.Getenv("JWT_SECRET")
	authLogic := core.NewAuthLogic(db, jwtSecret)

	// 3. WebSocket Hub
	hub := ws.NewHub()
	go hub.Run() 

	// 4. Redis Subscriber
	sub := worker.NewSubscriber(rdb.RDB, hub) 
	go sub.Start(context.Background())

	// 5. Handlers
	authHandler := &handler.AuthHandler{Logic: authLogic}
	wsHandler := &handler.WSHandler{Hub: hub, Redis: rdb} 

	// 6. Router 
	r := gateway.SetupRouter(jwtSecret, authHandler, wsHandler)
	
	log.Println("Gateway running on :8080")
	r.Run(":8080")
}