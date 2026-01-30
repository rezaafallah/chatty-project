package main

import (
	"log"
	"my-project/internal/app"
)

func main() {
	// 1. Initialize Application (Wiring)
	application, err := app.NewApplication()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// 2. Start Server (Running)
	application.Start(":8080")
}