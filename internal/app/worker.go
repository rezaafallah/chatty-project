package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/pkg/config"
	"my-project/pkg/consts"
	"my-project/pkg/logic"
	"my-project/pkg/repository"
	"my-project/pkg/logger"
)

type WorkerApp struct {
	DB        *postgres.DB
	Redis     *redis.Client
	ChatLogic *logic.ChatLogic
}

func NewWorkerApp(cfg *config.Config) (*WorkerApp, error) {
	db, err := postgres.New(cfg.DB_DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	rdb := redis.New(cfg.RedisAddr)
	msgRepo := repository.NewMessageRepository(db.Conn)
	chatLogic := logic.NewChatLogic(msgRepo, rdb, logger.Setup())

	return &WorkerApp{
		DB:        db,
		Redis:     rdb,
		ChatLogic: chatLogic,
	}, nil
}

func (w *WorkerApp) Start() {
	log.Println("Core Worker Started (Waiting for messages)...")

	// Context 
	ctx, cancel := context.WithCancel(context.Background())
	
	// listen to shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// run Goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := w.Redis.RDB.BLPop(ctx, 0, consts.QueueChatInbound).Result()
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					log.Println("Queue Error (Retry in 1s):", err)
					continue
				}

				payload := result[1]
				if err := w.ChatLogic.ProcessIncomingMessage(ctx, []byte(payload)); err != nil {
					log.Printf("Processing Failed: %v", err)
				}
			}
		}
	}()

	<-sigChan
	log.Println("Shutting down worker...")
	
	cancel()

	if err := w.DB.Conn.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}
	if err := w.Redis.RDB.Close(); err != nil {
		log.Println("Error closing Redis:", err)
	}

	log.Println("Worker exited gracefully")
}