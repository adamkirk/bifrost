#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    go fmt ./...
    go build -o ./bin/bifrost-server ./cmd/server/*.go
)