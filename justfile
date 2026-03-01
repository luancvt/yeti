# Run with live-reload (watches .go, .templ, .yang files)
dev:
    air

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

# Download YANG models from YangModels/yang repo
# Usage: just fetch-models xr-7112 vendor/cisco/xr/7112
fetch-models name repo_path:
    ./scripts/fetch-models.sh {{ name }} {{ repo_path }}
