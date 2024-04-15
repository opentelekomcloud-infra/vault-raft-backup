GO111MODULE=on
export PATH := /usr/local/go/bin:$(PATH)

exec_path := /usr/local/bin/
exec_name := vault-raft-backup
VERSION := 0.0.1

.PHONY: default fmt lint vet test build install release

default: lint build

fmt:
	@echo "Running go fmt"
	@go fmt ./...

lint:
	@echo "Running go lint"
	@golangci-lint run ./...

vet:
	@echo "Running go vet"
	@go vet ./...

build: fmt vet
	@echo "Building vault-raft-backup for Linux"
	@GOOS=linux go build -ldflags="-X 'main.version=$(VERSION)'" -o bin/$(exec_name)

install: build
	@echo "Installing vault-raft-backup to $(exec_path)"
	@cp ./bin/$(exec_name) $(exec_path)
	@echo "Vault-Raft-Backup installed to $(exec_path)$(exec_name)"

release:
	@echo "Packaging release for Linux"
	@tar -czf "./bin/$(exec_name)-$(VERSION)-linux-amd64.tgz" -C "./bin/" $(exec_name)
	@echo "Release version is built and packaged"
