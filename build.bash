#!/usr/bin/env bash

set -euf -o pipefail

export GOARCH=amd64
export CGO_ENABLED=0

BUILD_DATE="$(date -u)"
VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)

go vet

if [[ " $@ " =~ " -release" ]]; then
  mkdir -p bin
  GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Commit=$COMMIT -X main.Version=$VERSION -X 'main.BuildDate=$BUILD_DATE'" -a -o bin/collector .
else
  VERSION="$VERSION-$(hostname -f)-local"
  COMMIT="$COMMIT-$(git rev-parse --abbrev-ref HEAD)"
  go build -ldflags "-s -w -X main.Commit=$COMMIT -X 'main.Version=$VERSION' -X 'main.BuildDate=$BUILD_DATE'" -a -o collector .
fi
