run: build
	@./bin/auth-service

build:
	@echo "Building Auth service..."
	@go build -o ./bin/auth-service ./cmd/main.go
	@echo "Done"