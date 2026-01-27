package main

import (
	"log"
	"os"
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
)

func main() {
	// Init infrastructure
	_, err := postgres.New(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	_ = redis.New(os.Getenv("REDIS_ADDR"))

	log.Println("Core Worker Started (Listening to Queue)...")
	select {} // Block forever
}