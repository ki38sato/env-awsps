#!/bin/bash

rm -rf vendor/*
##https://github.com/Masterminds/glide/issues/654

echo "execute build.sh using golang:1.8"
docker run --rm -v "$(pwd)":/go/src/github.com/ki38sato/env-awsps -w /go/src/github.com/ki38sato/env-awsps golang:1.8 bash build.sh
