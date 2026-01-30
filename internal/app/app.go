package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/auth"
	"my-project/pkg/config"
	"my-project/pkg/logger"
	"my-project/pkg/logic"
	"my-project/pkg/repository"
	service "my-project/pkg/utils"
	"my-project/srv/gateway"
	"my-project/srv/gateway/handler"
	"my-project/srv/gateway/worker"
	"my-project/srv/gateway/ws"
)

// Application
type Application struct {
	DB     *postgres.DB
	Redis  *redis.Client
	Router *gin.Engine
	Config *config.Config
}

// NewApplication
func NewApplication(cfg *config.Config) (*Application, error) {
	// 1. Config & Infra
	db, err := postgres.New(cfg.DB_DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	rdb := redis.New(cfg.RedisAddr)

	// 2. Repositories & Tools
	userRepo := repository.NewUserRepository(db.Conn)
	msgRepo := repository.NewMessageRepository(db.Conn)

	tokenMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	sanitizer := service.NewSanitizer()

	// 3. Logic (Services)
	authLogic := logic.NewAuthLogic(userRepo, tokenMgr)
	chatLogic := logic.NewChatLogic(msgRepo, rdb, logger.Setup())

	// 4. WebSocket Hub & Workers
	hub := ws.NewHub(rdb)
	go hub.Run()

	subscriber := worker.NewSubscriber(rdb.RDB, hub)
	go subscriber.Start(context.Background())

	// 5. Handlers
	authHandler := &handler.AuthHandler{Logic: authLogic}
	
	wsHandler := &handler.WSHandler{
		Hub:       hub,
		Broker:    rdb,
		Sanitizer: sanitizer,
		Log:       logger.Setup(),
	}
	
	chatHandler := &handler.ChatHandler{Logic: chatLogic}

	// 6. Router
	router := gateway.SetupRouter(cfg.JWTSecret, authHandler, wsHandler, chatHandler)

	return &Application{
		DB:     db,
		Redis:  rdb,
		Router: router,
		Config: cfg,
	}, nil
}

// Start
func (app *Application) Start(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Router,
	}

	go func() {
		log.Printf("Gateway running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful Shutdown Logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := app.DB.Conn.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}
	if err := app.Redis.RDB.Close(); err != nil {
		log.Println("Error closing Redis:", err)
	}

	log.Println("Server exiting")
}