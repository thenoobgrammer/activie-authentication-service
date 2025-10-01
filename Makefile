include .env
export

DSN := $(DSN)
MIGRATIONS_PATH := internal/infra/database/migrations/
SERVICE_NAME := Authorization service

run: build
	@./bin/authentication-service

build:
	@echo "Building Authentication service..."
	@go build -o ./bin/authentication-service ./cmd/main.go
	@echo "Done"

# Database migrations commands
#------------------------------------------
up:
	@echo "Running migrate up with DSN: $(DSN)"
	@migrate -path $(MIGRATIONS_PATH) -database '$(DSN)' -verbose up

down:
	@echo "Running $(SERVICE_NAME) with DSN: $(DSN)"
	@migrate -path $(MIGRATIONS_PATH) -database '$(DSN)' -verbose down $(N)

forcev:
	@echo "Running $(SERVICE_NAME) with DSN: $(DSN)"
	@migrate -path $(MIGRATIONS_PATH) -database '$(DSN)' force $(version)

migration:
	@echo "Creating migration with name: $(name)"
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
