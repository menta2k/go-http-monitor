APP_NAME   := go-http-monitor
GO         := go
NPM        := npm
FRONTEND   := frontend
DIST       := $(FRONTEND)/dist
BUILD_DIR  := bin
BINARY     := $(BUILD_DIR)/$(APP_NAME)
GOFLAGS    ?=
LDFLAGS    ?= -s -w

.PHONY: all build clean frontend frontend-install frontend-clean backend run dev test lint install uninstall help

## all: Build everything (frontend + backend) into a single binary
all: build

## build: Build frontend then compile Go binary with embedded SPA
build: frontend backend

## backend: Compile Go binary (expects frontend/dist to exist)
backend: | $(BUILD_DIR)
	@test -d $(DIST) || (echo "Error: $(DIST) not found. Run 'make frontend' first." && exit 1)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) .
	@echo "Built $(BINARY)"

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

## frontend-install: Install frontend npm dependencies
frontend-install:
	cd $(FRONTEND) && $(NPM) install

## frontend: Install deps and build the Vue SPA into frontend/dist
frontend: frontend-install
	cd $(FRONTEND) && $(NPM) run build

## frontend-clean: Remove frontend build artifacts and node_modules
frontend-clean:
	rm -rf $(DIST) $(FRONTEND)/node_modules

## clean: Remove all build artifacts
clean: frontend-clean
	rm -rf $(BUILD_DIR)

## run: Build and run the application (requires JWT_SECRET and ADMIN_PASSWORD env vars)
run: build
	$(BINARY)

## dev: Run frontend dev server and Go backend concurrently
dev:
	@echo "Starting Go backend on :8080..."
	$(GO) run . &
	@echo "Starting Vite dev server..."
	cd $(FRONTEND) && $(NPM) run dev

## test: Run Go tests
test:
	$(GO) test ./auth/... ./checker/... ./config/... ./database/... ./domain/... ./monitor/... ./notification/... ./notifier/... ./result/... ./response/...

## lint: Run go vet
lint:
	$(GO) vet ./auth/... ./checker/... ./config/... ./database/... ./domain/... ./monitor/... ./notification/... ./notifier/... ./result/... ./response/...

## install: Install binary and systemd service (run as root)
install: build
	install -m 755 $(BINARY) /usr/local/bin/$(APP_NAME)
	install -m 644 deploy/go-http-monitor.service /usr/lib/systemd/system/
	mkdir -p /usr/share/$(APP_NAME)
	install -m 644 deploy/env.example /usr/share/$(APP_NAME)/
	bash deploy/postinstall.sh

## uninstall: Remove binary and systemd service (run as root)
uninstall:
	bash deploy/preremove.sh
	rm -f /usr/local/bin/$(APP_NAME)
	rm -f /usr/lib/systemd/system/go-http-monitor.service
	systemctl daemon-reload

## help: Show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //' | column -t -s ':'
