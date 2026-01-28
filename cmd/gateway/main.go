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
	"my-project/srv/gateway/ws" // Import جدید
)

func main() {
	// 1. Infra
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.New(os.Getenv("REDIS_ADDR"))

	// 2. Logic
	authLogic := core.NewAuthLogic(db)

	// 3. WebSocket Hub (Start Logic)
	hub := ws.NewHub()
	go hub.Run() // اجرای هاب در بک‌گراند

	// 4. Redis Subscriber (اتصال Core به Hub)
	// حالا Hub ما متد BroadcastToUser را دارد و با اینترفیس Subscriber سازگار است
	sub := worker.NewSubscriber(rdb.RDB, hub) 
	go sub.Start(context.Background())

	// 5. Handlers
	authHandler := &handler.AuthHandler{Logic: authLogic}
	wsHandler := &handler.WSHandler{Hub: hub, Redis: rdb} // هندلر جدید

	// 6. Router (باید روت WS را اضافه کنید)
	// نکته: باید تابع SetupRouter را آپدیت کنید که wsHandler را هم بگیرد
	r := gateway.SetupRouter(os.Getenv("JWT_SECRET"), authHandler, wsHandler)
	
	log.Println("Gateway running on :8080")
	r.Run(":8080")
}