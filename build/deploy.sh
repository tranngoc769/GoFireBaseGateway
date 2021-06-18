#!/bin/sh
echo "Copy environment file"
yes | cp -rf build/go-firebase-gateway-env /root/go/env/go-firebase-gateway-env
echo "Build go application"
go mod tidy
go get .
GOOS=linux GOARCH=amd64 go build -o go-firebase-gateway-api main.go
echo "Restart service"
systemctl restart go-firebase-gateway
