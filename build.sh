#!/bin/bash

set -x
set -e

mkdir -p binary

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o binary/gopxy_darwin cmd/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o binary/gopxy_windows.exe cmd/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binary/gopxy_linux cmd/main.go

tar -czf binary.tar.gz ./binary

echo "binary.tar.gz build succ..."
