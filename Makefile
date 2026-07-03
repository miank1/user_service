APP_NAME=user-service

.PHONY: run dev build start test tidy clean docker

# Development (Live Reload)
dev:
	air

# Run without Air
run:
	go run ./cmd/main.go

# Build binary
build:
	go build -o bin/$(APP_NAME) ./cmd/main.go

# Run compiled binary
start: build
	./bin/$(APP_NAME)

# Run tests
test:
	go test ./...

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Clean generated files
clean:
	rm -rf bin
	rm -rf tmp

# Build Docker image
docker:
	docker build -t $(APP_NAME) .
