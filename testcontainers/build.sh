#!/bin/bash
set -e

# build protos
protoc --proto_path=internal/api internal/api/*.proto --go_out=internal/api --go-grpc_out=internal/api

# build binaries
go build .

# build docker images
docker build -t exchange -f ExchangeDockerfile .
docker build -t account -f AccountDockerfile .
