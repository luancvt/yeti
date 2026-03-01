# Generate templ files
generate:
    templ generate

# Run the app
run: generate
    go run ./cmd/yeti/

# Build the app
build: generate
    go build -o yeti ./cmd/yeti/
    just clean

# Clean build artifacts
clean:
    rm -f yeti
