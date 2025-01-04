# Define variables
BINARY_NAME := sway-windows
INSTALL_DIR := $(HOME)/.local/share/pop-launcher/plugins/$(BINARY_NAME)

# Default target: Build the binary
.PHONY: all
all: build

build:
	@echo "Building the plugin..."
	go build -o $(BINARY_NAME)

install: build
	@echo "Installing the binary and plugin definition into..."
	@echo $(INSTALL_DIR)
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)
	cp plugin.ron $(INSTALL_DIR)

clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(BINARY_NAME)

uninstall:
	@echo "Uninstalling the plugin..."
	rm -rf $(INSTALL_DIR)

# Phony targets
.PHONY: build install uninstall clean
