#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

set -e

# Build admin-go
cd ${SCRIPT_DIR}/../admin-container
go mod vendor
export CGO_ENABLED=0
export GOOS=linux
[ ! -f admin-go ] || rm admin-go
go build -o admin-go *.go

# Build container
podman image rm --force admin-container:v1
podman build -t admin-container:v1 .