.PHONY: build test clean run

# Build the devesis binary
build:
	go build -o bin/devesis ./cmd/devesis

# Run tests (essential for TDD)
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the application
run: build
	./bin/devesis