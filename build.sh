#!/usr/bin/env /bin/sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o output/mcbot-linux-amd64 main.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o output/mcbot-linux-arm64 main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o output/mcbot-windows-amd64.exe main.go