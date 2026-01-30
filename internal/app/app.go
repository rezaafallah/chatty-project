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
	"my-project/pkg/logic"
	"my-project/pkg/repository"
	"my-project/pkg/utils"
	"my-project/pkg/auth"
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
}

// NewApplication
func NewApplication() (*Application, error) {
	// 1. Config & Infra
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN is not set")
	}
	db, err := postgres.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is not set")
	}
	rdb := redis.New(redisAddr)

	// 2. Repositories & Tools
	userRepo := repository.NewUserRepository(db.Conn)
	msgRepo := repository.NewMessageRepository(db.Conn)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}
	tokenMgr := auth.NewJWTManager(jwtSecret, 72*time.Hour)
	sanitizer := service.NewSanitizer()

	// 3. Logic (Services)
	authLogic := logic.NewAuthLogic(userRepo, tokenMgr)
	chatLogic := logic.NewChatLogic(msgRepo, rdb)

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
	}
	chatHandler := &handler.ChatHandler{Logic: chatLogic}

	// 6. Router
	router := gateway.SetupRouter(jwtSecret, authHandler, wsHandler, chatHandler)

	return &Application{
		DB:     db,
		Redis:  rdb,
		Router: router,
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