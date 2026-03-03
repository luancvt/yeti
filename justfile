# Run with live-reload (watches .go, .templ, .yang, .css files)
dev:
    air

# Generate templ files
generate:
    templ generate

# Build Tailwind CSS
css:
    tailwindcss -i static/css/input.css -o static/css/tailwind.css --minified

# Run the app
run: generate css
    go run ./cmd/yeti/

# Verify the app compiles
build: generate css
    go build -o /dev/null ./cmd/yeti/

# Format code
fmt:
    templ fmt .
    golangci-lint fmt

# Lint code
lint:
    golangci-lint run

# Run tests
test:
    go test ./...

# Check generated templ files are up to date
check:
    templ generate --check
    golangci-lint fmt --diff

# Download YANG models from YangModels/yang repo

# Usage: just fetch-models xr-7112 vendor/cisco/xr/7112
fetch-models name repo_path:
    ./scripts/fetch-models.sh {{ name }} {{ repo_path }}
