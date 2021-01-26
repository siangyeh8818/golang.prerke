#!/bin/bash

#mac 版的binary
go build -o build/mac/prerke_darwin-amd64
 cmd/main.go

#Linux版的binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/prerke_linux-amd64 cmd/main.go 
