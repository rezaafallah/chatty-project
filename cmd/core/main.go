package main

import (
	"context"
	"log"
	"os"

	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/pkg/consts"
	"my-project/pkg/logic"
	"my-project/pkg/repository"
)

func main() {
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.New(os.Getenv("REDIS_ADDR"))

	msgRepo := repository.NewMessageRepository(db.Conn)
	chatLogic := logic.NewChatLogic(msgRepo, rdb)
	
	log.Println("Core Worker Started (Listening to Queue)...")

	runWorker(rdb, chatLogic)
}

func runWorker(rdb *redis.Client, chatLogic *logic.ChatLogic) {
	ctx := context.Background()
	for {
		result, err := rdb.RDB.BLPop(ctx, 0, consts.QueueChatInbound).Result()
		if err != nil {
			log.Println("Queue Error:", err)
			continue
		}

		payload := result[1]
		
		err = chatLogic.ProcessIncomingMessage(ctx, []byte(payload))
		if err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}
}