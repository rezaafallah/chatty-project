package main

import (
	"log"
	"my-project/internal/app"
	"my-project/pkg/config"
)

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// 2. App
	worker, err := app.NewWorkerApp(cfg)
	if err != nil {
		log.Fatal("App init error:", err)
	}

	// 3. Start (Blocking)
	worker.Start()
}