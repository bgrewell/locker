.PHONY: all clean proto install_protoc

all: proto locker lockerd pam_locker

# Install the required protoc plugins.
install_protoc:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Compile the protobuf file.
proto: install_protoc
	@echo "Compiling locker.proto..."
	protoc --proto_path=api/. --go_out=paths=source_relative:api/go/ --go-grpc_out=paths=source_relative:api/go/ locker.proto

locker:
	go build -o bin/locker cmd/locker/main.go

lockerd:
	go build -o bin/lockerd cmd/lockerd/main.go

pam_locker:
	go build -buildmode=c-shared -o bin/pam_locker.so cmd/module/main.go

clean:
	rm -f bin/locker bin/lockerd bin/pam_locker.so bin/pam_locker.h