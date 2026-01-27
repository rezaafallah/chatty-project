package main

import (
	"my-project/internal/adapter/postgres"
	"my-project/internal/adapter/redis"
	"my-project/internal/core"
	"my-project/pkg/logger"
	"my-project/srv"
	"os"
)

func main() {
	log := logger.Setup()
	
	// Init Infra
	db, _ := postgres.New(os.Getenv("DB_DSN"))
	rdb := redis.New(os.Getenv("REDIS_ADDR"))

	// Init Logic (The Brain)
	appLogic := core.NewLogic(db, rdb, log)

	// Init Worker Server
	worker := &srv.WorkerServer{Logic: appLogic}
	
	log.Info("Starting Core Worker...")
	worker.Start()
}