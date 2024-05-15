run: build
	@./bin/authentication-service

build:
	@echo "Building Authentication service..."
	@go build -o ./bin/authentication-service ./cmd/main.go
	@echo "Done"