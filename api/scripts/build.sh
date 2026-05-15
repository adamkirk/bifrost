#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    go fmt ./...
    go build -o ./bin/server ./cmd/server/*.go
)