# The 'mocks' target generates mock implementations for interfaces using 'go generate'.
# It first prints a message indicating that mocks are being generated,
# then runs 'go generate' recursively in all subdirectories.

# The 'up' target depends on 'mocks', ensuring that mocks are generated before running the application.
# It prints a message indicating that the application is starting,
# then runs the application using 'go run .'.
mocks:
	@echo "Generating mocks..."
	go generate ./...

up: mocks test
	@echo "Running the application..."
	go run .

test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed."

lint:
	@echo "Running linter..."
	golangci-lint run ./...

install:
	@echo "Installing dependencies..."
	go mod tidy
	go mod vendor

install-tools:
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@latest

.PHONY: mocks up test
