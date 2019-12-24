#!/usr/bin/env bash

GOCMD="go"
GOBUILD="go build"
GOGET="go get"
BINARY_NAME="sahc"



build()
{
  $( ${GOBUILD} -o bin/${BINARY_NAME} -v ./cmd/sahc/main.go)
}

go-test()
{
  go test  -v ./...
}

go-run(){
  SAHC_CONFIG=services.yaml go run ./cmd/sahc/main.go
}

go-clean(){
  go clean && rm -rf bin
}


if [[ "$1" == "build" ]]; then
    build
elif [[ "$1" == "test" ]]; then
    go-test
elif [[ "$1" == "run" ]]; then
    go-run
elif [[ "$1" == "run" ]]; then
    go-clean
fi
