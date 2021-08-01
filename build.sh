#!/usr/bin/env bash

set -e

go get github.com/gorilla/websocket #todo vendor this or learn about the module stuff!
go get -u golang.org/x/lint/golint
go get -u github.com/client9/misspell/cmd/misspell

ls *.md | xargs misspell -error

go test ./...
go vet ./...
go fmt ./...
golint ./...
