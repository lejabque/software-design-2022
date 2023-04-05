#!/bin/bash
set -e

export GOARCh=amd64 GOOS=linux

# build protos
protoc --proto_path=internal/api internal/api/*.proto --go_out=internal/api --go-grpc_out=internal/api

# build binaries
go build .

# build docker images
docker build --platform linux/amd64 -t exchange -f ExchangeDockerfile .
