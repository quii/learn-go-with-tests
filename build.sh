#!/usr/bin/env bash

set -e

go install github.com/client9/misspell/cmd/misspell@latest
go install github.com/po3rin/gofmtmd/cmd/gofmtmd@latest

MD_FILES="./*.md ./go-fundamentals/*.md ./testing-fundamentals/*.md ./build-an-application/*.md ./questions-and-answers/*.md ./meta/*.md"

ls $MD_FILES | xargs misspell -error

for md_file in $MD_FILES; do
    echo "formatting  file: $md_file"
    gofmtmd  "$md_file" -r
done

go test ./...
go vet ./...
go fmt ./...
