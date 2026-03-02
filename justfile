# Run with live-reload (watches .go, .templ, .yang files)
dev:
    air

# Generate templ files
generate:
    templ generate

# Run the app
run: generate
    go run ./cmd/yeti/

# Verify the app compiles
build: generate
    go build -o /dev/null ./cmd/yeti/

# Run tests
test:
    go test ./...

# Check generated templ files are up to date
check:
    templ generate --check

# Download YANG models from YangModels/yang repo
# Usage: just fetch-models xr-7112 vendor/cisco/xr/7112
fetch-models name repo_path:
    ./scripts/fetch-models.sh {{ name }} {{ repo_path }}
