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

.PHONY: mocks up test
