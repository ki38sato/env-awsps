#!/bin/bash

VERSION=$(cat ./VERSION)
#PACKAGE="github.com/ki38sato/env-awsps"

go get -v github.com/Masterminds/glide
go install github.com/Masterminds/glide
export GO15VENDOREXPERIMENT=1
glide up

export GOOS=linux
export GOARCH=amd64
go build -v -ldflags "-X main.version=${VERSION}" -o build/env-awsps