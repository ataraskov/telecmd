# Variables
BINARY_NAME=telecmd

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(BINARY_NAME) cmd/telecmd/main.go

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Test the application
.PHONY: test
test:
	go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Build Docker image
.PHONY: build-docker
build-docker:
	docker build -t $(BINARY_NAME) .

# Run Docker container
.PHONY: run-docker
run-docker: build-docker
	docker run --rm -it --env-file .env --name $(BINARY_NAME) $(BINARY_NAME)
