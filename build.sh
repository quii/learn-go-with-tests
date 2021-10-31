#!/usr/bin/env bash

set -e

go get github.com/gorilla/websocket #todo vendor this or learn about the module stuff!
go get -u golang.org/x/lint/golint
go get -u github.com/client9/misspell/cmd/misspell
go get github.com/po3rin/gofmtmd/cmd/gofmtmd

ls *.md | xargs misspell -error
for md_file in ./*.md; do
    echo "formatting  file: $md_file"
    gofmtmd "$md_file" -r
done

go test ./...
go vet ./...
go fmt ./...
golint ./...
