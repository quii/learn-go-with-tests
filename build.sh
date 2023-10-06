#!/usr/bin/env bash

set -e

go install github.com/client9/misspell/cmd/misspell@latest
go install github.com/po3rin/gofmtmd/cmd/gofmtmd@latest

ls *.md | xargs misspell -error

for md_file in ./*.md; do
    echo "formatting  file: $md_file"
    gofmtmd  "$md_file" -r
done

go test ./...
go vet ./...
go fmt ./...
