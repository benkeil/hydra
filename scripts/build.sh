#!/bin/bash

APP_VERSION=$1
go build -ldflags "-X main.version=${APP_VERSION}" -o hydra cmd/*.go