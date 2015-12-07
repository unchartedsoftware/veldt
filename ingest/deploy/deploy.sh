#!/bin/bash

TARGET_ARCH="amd64"
TARGET_OS="linux"

TARGET_USER="etl"
TARGET_HOST="10.63.254.14"
TARGET_DIR="~/es-ingest/"

TARGET_SH="es-ingest.sh"
TARGET_EXE="es-ingest.bin"

# build target executable
GOARCH=$TARGET_ARCH GOOS=$TARGET_OS go build -o ./$TARGET_EXE ../main.go

# ensure dest directory exists
ssh $TARGET_USER@$TARGET_HOST "mkdir $TARGET_DIR"

# scp executable and shell script into target directory
scp $TARGET_SH $TARGET_USER@$TARGET_HOST:$TARGET_DIR
scp $TARGET_EXE $TARGET_USER@$TARGET_HOST:$TARGET_DIR
