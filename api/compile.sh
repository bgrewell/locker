#!/bin/bash

# Go
protoc --proto_path=. --go_out=paths=source_relative:go/. --go-grpc_out=paths=source_relative:go/. ./locker.proto