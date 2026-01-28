package main

import (
	"context"
	"log"
	"os"

	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/core"
	"my-project/pkg/consts"
)

func main() {
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.New(os.Getenv("REDIS_ADDR"))
	chatLogic := core.NewChatLogic(db, rdb)
	log.Println("Core Worker Started (Listening to Queue)...")

	ctx := context.Background()
	for {
		result, err := rdb.RDB.BLPop(ctx, 0, consts.QueueChatInbound).Result()
		if err != nil {
			log.Println("Queue Error:", err)
			continue
		}

		payload := result[1]
		
		log.Printf("Processing message: %s", payload)
		err = chatLogic.ProcessIncomingMessage(ctx, []byte(payload))
		if err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}
}