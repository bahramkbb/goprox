#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -o goprox *.go

docker-compose down
docker-compose up -d
