package main

import (
	"log"
	"my-project/internal/app"
	"my-project/pkg/config"
)

func main() {
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. Initialize App
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// 3. Start
	application.Start(":8080")
}