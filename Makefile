.PHONY: fmt test tidy build install uninstall register-handler all

INSTALL_BIN   ?= /usr/local/bin/dms-profiler
DESKTOP_DIR   := $(HOME)/.local/share/applications
CONFIG_DIR    := $(HOME)/.config/dms-profiler
DESKTOP_FILE  := linux/io.github.dms-profiler.desktop

fmt:
	go fmt ./...

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o dms-profiler ./cmd/dms-profiler

install: build
	sudo install -m 0755 dms-profiler "$(INSTALL_BIN)"
	mkdir -p "$(DESKTOP_DIR)"
	@if [ -f "$(DESKTOP_DIR)/io.github.dms-profiler.desktop" ]; then \
		cp "$(DESKTOP_DIR)/io.github.dms-profiler.desktop" "$(DESKTOP_DIR)/io.github.dms-profiler.desktop.bak"; \
		echo "Backed up existing desktop entry to $(DESKTOP_DIR)/io.github.dms-profiler.desktop.bak"; \
	fi
	cp "$(DESKTOP_FILE)" "$(DESKTOP_DIR)/"
	@if [ ! -f "$(CONFIG_DIR)/config.toml" ]; then \
		mkdir -p "$(CONFIG_DIR)"; \
		cp config.example.toml "$(CONFIG_DIR)/config.toml"; \
		echo "Config created at $(CONFIG_DIR)/config.toml — edit it to match your Chrome profiles."; \
	else \
		echo "Config already exists at $(CONFIG_DIR)/config.toml — not overwritten."; \
	fi
	bash linux/set-default-http-handler.sh

uninstall:
	sudo rm -f "$(INSTALL_BIN)"
	rm -f "$(DESKTOP_DIR)/io.github.dms-profiler.desktop"
	xdg-mime default "" x-scheme-handler/http  || true
	xdg-mime default "" x-scheme-handler/https || true
	@echo "Config at $(CONFIG_DIR)/ was not removed. Delete manually if desired."

register-handler:
	bash linux/set-default-http-handler.sh

all: fmt test
