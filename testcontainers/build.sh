#!/bin/bash

# build protos
protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

# build binaries
go build .

# build docker images
docker build -t exchange exchange/
docker build -t account account/
