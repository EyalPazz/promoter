INSTALL_PATH = /usr/local/bin/promoter
VERSION ?= v0.9.0
LDFLAGS = -ldflags="-X 'promoter/cmd.Version=${VERSION}'"

install:
	@echo "Installing promoter binary..."
	./install.sh

uninstall:
	@echo "Uninstalling promoter binary..."
	rm -f $(INSTALL_PATH)

build_macos_intel:
	GOOS=darwin GOARCH=amd64 go build -o promoter_darwin_amd64 -ldflags="-X 'promoter/cmd.Version=${VERSION}'"

build_macos_silicon:
	GOOS=darwin GOARCH=arm64 go build -o promoter_darwin_arm64 -ldflags="-X 'promoter/cmd.Version=${VERSION}'"

build_linux_intel:
	GOOS=linux GOARCH=amd64 go build -o promoter_linux_amd64 -ldflags="-X 'promoter/cmd.Version=${VERSION}'"

build_linux_arm:
	GOOS=linux GOARCH=arm64 go build -o promoter_linux_arm64 -ldflags="-X 'promoter/cmd.Version=${VERSION}'"


build: build_macos_intel build_macos_silicon build_linux_intel build_linux_arm

reinstall: uninstall install

.PHONY: uninstall install reinstall build_macos_intel build_macos_silicon build_linux_intel build_linux_arm build

.DEFAULT_GOAL := build
