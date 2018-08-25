#!/usr/bin/env bash

GOOS=js GOARCH=wasm go build -o cmd/web/html/test.wasm main.go
