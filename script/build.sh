#!/bin/bash
if [ ! -d "./admin-container" ]; 
    then echo "Please run this script from root folder of git project" 
    exit 1 
fi

# Build admin-go
cd admin-container
go mod tidy
export CGO_ENABLED=0
export GOOS=linux
go build -o admin-go admin.go

# Build container
podman build -t admin-container:v1 .