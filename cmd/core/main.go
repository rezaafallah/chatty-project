package main

import (
	"context"
	"log"
	"os"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/core"
)

func main() {
	// 1. Init Infra
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.New(os.Getenv("REDIS_ADDR"))

	// 2. Init Logic
	chatLogic := core.NewChatLogic(db, rdb)

	log.Println("Core Worker Started (Listening to Queue)...")

	// 3. Start Consumer Loop
	ctx := context.Background()
	for {
		// خواندن از صف (QueueChatInbound = "chat.inbound")
		// فرض بر این است که متد PopQueue یا مشابه در redis client داری
		// اگر نداری باید از rdb.RDB.BLPop استفاده کنی
		result, err := rdb.RDB.BLPop(ctx, 0, "chat.inbound").Result()
		if err != nil {
			log.Println("Queue Error:", err)
			continue
		}

		// result[0] نام صف است، result[1] محتوای پیام (JSON)
		payload := result[1]
		
		log.Printf("Processing message: %s", payload)
		
		err = chatLogic.ProcessIncomingMessage(ctx, []byte(payload))
		if err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}
}