.PHONY: run build test clean install uninstall

INSTALL_DIR ?= $(HOME)/.local/bin

run:
	go run ./cmd/timer

build:
	go build -o timer ./cmd/timer

test:
	go test ./...

install:
	go build -o $(INSTALL_DIR)/timer ./cmd/timer

uninstall:
	rm -f $(INSTALL_DIR)/timer

clean:
	rm -f timer
