.PHONY: run-gateway run-core up down build clean tidy

run-gateway:
	go run cmd/gateway/main.go

run-core:
	go run cmd/core/main.go

up:
	docker-compose up --build -d
	@echo "Services started! Gateway: http://localhost:8080"

down:
	docker-compose down

logs:
	docker-compose logs -f

tidy:
	go mod tidy
	go mod vendor

test:
	go test ./...