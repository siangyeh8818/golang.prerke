#!/bin/bash

#mac 版的binary
go build -o build/mac/prerke cmd/main.go

#Linux版的binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/prerke cmd/main.go 