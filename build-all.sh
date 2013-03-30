#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <app>"
    exit 1
fi

# Load crosscompile environment (https://github.com/davecheney/golang-crosscompile)
source /Users/pii/scripts/golang-crosscompile/crosscompile.bash

PLATFORMS="darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 linux/386 linux/amd64 linux/arm linux/arm5 windows/386 windows/amd64"
APP_NAME=$1

# Remove old binaries
rm -rf bin/*

# Build binary for each platform in parallel
for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_NAME=$APP_NAME
    BIN_PATH="bin/${GOOS}-${GOARCH}"

    if [ $GOOS == "windows" ]; then
        BIN_NAME="${BIN_NAME}.exe"
    fi

    if [ $GOARCH == "arm5" ]; then
        export GOARM=5
        GOARCH="arm"
    else
        unset GOARM
    fi

    # Create platform directory and copy in example config file
    mkdir -p ${BIN_PATH}
    cp config.json ${BIN_PATH}/config.json

    BUILD_CMD="go-${GOOS}-${GOARCH} build -o ${BIN_PATH}/${BIN_NAME} ${APP_NAME}.go"

    echo "Building ${BIN_PATH}/${BIN_NAME}"
    $BUILD_CMD &
done

# Wait for builds to complete
for job in $(jobs -p); do
    wait $job
done

echo "All done"
