#!/bin/sh
# XC: build apimain inside alpine with CGO for the s6/alpine runtime.
# The API links mattn/go-sqlite3 which REQUIRES cgo - a CGO_ENABLED=0 build
# produces a stub that only fails at runtime. Built in the same alpine family
# as the production s6 image so the dynamic musl binary matches.
set -e
apk add --no-cache gcc musl-dev git
cd /src
go mod tidy
CGO_ENABLED=1 go build -buildvcs=false -ldflags="-s -w -linkmode external -extldflags -static" -o dist/alpine/apimain ./cmd
echo "=== result ==="
ls -la dist/alpine/apimain
