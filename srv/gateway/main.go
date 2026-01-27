package main

import (
	"log"
	"os"
	"my-project/internal/adapter/postgres"
	"my-project/internal/core"
	"my-project/srv/gateway"
	"my-project/srv/gateway/handler"
)

func main() {
	// Infra
	db, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	// Logic
	authLogic := core.NewAuthLogic(db)

	// Handler
	authHandler := &handler.AuthHandler{Logic: authLogic}

	// Server
	r := gateway.SetupRouter(os.Getenv("JWT_SECRET"), authHandler)
	r.Run(":8080")
}