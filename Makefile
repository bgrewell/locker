# Pull versioning info from git (or set defaults if desired).
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE  ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "0000000")
BRANCH      ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

.PHONY: all clean proto install_protoc install_deps pam_locker

all: proto locker lockerd pam_locker

## 1) Install the required protoc plugins.
install_protoc:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## 2) Compile the protobuf file.
proto: install_protoc
	@echo "Compiling locker.proto..."
	protoc --proto_path=api/. --go_out=paths=source_relative:api/go/ \
	       --go-grpc_out=paths=source_relative:api/go/ locker.proto

## 3) Build targets, injecting version information via -ldflags.
locker:
	go build -ldflags \
		"-X 'main.version=$(VERSION)' \
		 -X 'main.builddate=$(BUILD_DATE)' \
		 -X 'main.commit=$(COMMIT)' \
		 -X 'main.branch=$(BRANCH)'" \
		-o bin/locker cmd/locker/main.go

lockerd:
	go build -ldflags \
		"-X 'main.version=$(VERSION)' \
		 -X 'main.builddate=$(BUILD_DATE)' \
		 -X 'main.commit=$(COMMIT)' \
		 -X 'main.branch=$(BRANCH)'" \
		-o bin/lockerd cmd/lockerd/main.go

## 4) Ensure libpam0g-dev is installed, then build the module.
pam_locker: install_deps
	go build -buildmode=c-shared -o bin/pam_locker.so cmd/module/main.go

## 5) Check dependencies and install them if missing.
install_deps:
	@dpkg -s libpam0g-dev >/dev/null 2>&1 || ( \
		echo "Installing libpam0g-dev..." && \
		sudo apt-get update && \
		sudo apt-get install -y libpam0g-dev \
	)

## 6) Clean up binaries.
clean:
	rm -f bin/locker bin/lockerd bin/pam_locker.so bin/pam_locker.h
