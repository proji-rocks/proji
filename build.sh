#!/usr/bin/env bash

PACKAGE="./cmd/proji/"
PACKAGE_SPLIT_NAME=(${PACKAGE//\// })
PACKAGE_NAME="${PACKAGE_SPLIT_NAME[-1]}"
OUTPUT_PATH="bin/"
PLATFORMS=("darwin/amd64" "linux/amd64" "windows/amd64")

echo
echo " -> Building Cross-Platform Binaries"
echo

for PLATFORM in "${PLATFORMS[@]}"; do
    PLATFORM_SPLIT_NAME=(${PLATFORM//\// })
    GOOS="${PLATFORM_SPLIT_NAME[0]}"
    GOARCH="${PLATFORM_SPLIT_NAME[1]}"
    BIN_PATH="${OUTPUT_PATH}${GOOS}/${PACKAGE_NAME}"

    echo " ---> Building ${GOARCH} Binary for ${GOOS}..."
    CMD=""

    if [ $GOOS = "linux" ]; then
        # Linux build
        CMD="CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -ldflags '-s -w' -o ${BIN_PATH} ${PACKAGE}"
    elif [ $GOOS = "windows" ]; then
        # Windows build
        BIN_PATH+='.exe'
        CMD="CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} CC=x86_64-w64-mingw32-gcc go build -a -ldflags '-s -w' -o ${BIN_PATH} ${PACKAGE}"
    elif [ $GOOS = "darwin" ]; then
        # Mac build
        CMD="CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} CC=o64-clang  go build -a -ldflags '-s -w' -o ${BIN_PATH} ${PACKAGE}"
    else
        echo "Error: OS not support!"
        exit 1
    fi

    eval ${CMD} &
    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done
wait