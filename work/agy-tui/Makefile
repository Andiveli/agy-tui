BINARY_NAME := agy-tui
BUILD_DIR := ./build
PLUGIN_DIR := $(HOME)/.gemini/antigravity-cli/plugins/agy-tui
INSTALL_DIR := $(HOME)/.local/bin

VERSION := 0.1.0
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build install dev clean plugin-install plugin-uninstall

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/agy-tui

dev:
	go run ./cmd/agy-tui

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"

plugin-install: build
	@echo "Installing plugin to $(PLUGIN_DIR)..."
	@mkdir -p $(PLUGIN_DIR)
	cp plugin/plugin.json $(PLUGIN_DIR)/plugin.json
	@[ -f plugin/hooks.json ] && cp plugin/hooks.json $(PLUGIN_DIR)/hooks.json || true
	@echo "Plugin installed. Run 'agy plugin list' to verify."

plugin-uninstall:
	@echo "Removing plugin..."
	rm -rf $(PLUGIN_DIR)
	@echo "Plugin removed."

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

lint:
	golangci-lint run ./...
