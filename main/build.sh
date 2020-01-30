#!/bin/bash

#go get -v -t -d ./...
env GOOS=linux GOARCH=amd64 go build -v -o ./services main.go
