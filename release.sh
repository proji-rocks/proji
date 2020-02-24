#!/usr/bin/env bash

OUTPUT_PATH=./bin
RELEASE_PATH=./bin/release/

MAC_BIN_PATH=${OUTPUT_PATH}/darwin/proji
MAC_SHA_PATH=${OUTPUT_PATH}/darwin/proji-mac-sha256.txt
MAC_RLS_PATH=${RELEASE_PATH}/proji-mac.tar.gz

LNX_BIN_PATH=${OUTPUT_PATH}/linux/proji
LNX_SHA_PATH=${OUTPUT_PATH}/linux/proji-linux-sha256.txt
LNX_RLS_PATH=${RELEASE_PATH}/proji-linux.tar.gz

WIN_BIN_PATH=${OUTPUT_PATH}/windows/proji.exe
WIN_SHA_PATH=${OUTPUT_PATH}/windows/proji-windows-sha256.txt
WIN_RLS_PATH=${RELEASE_PATH}/proji-windows.zip

echo
echo " -----------------------------"
echo " |||  Generating Releases  |||"
echo " -----------------------------"

# Build the binaries
./build.sh

# Generating SHA256 hashes
echo
echo " -> Generating SHA256 Hashes of Binaries"
sha256sum ${MAC_BIN_PATH} | awk '{print $1;}' > ${MAC_SHA_PATH}
sha256sum ${LNX_BIN_PATH} | awk '{print $1;}' > ${LNX_SHA_PATH}
sha256sum ${WIN_BIN_PATH} | awk '{print $1;}' > ${WIN_SHA_PATH}

# Compressing binaries
echo
echo " -> Creating Compressed Archives"
tar -I pigz -cf ${MAC_RLS_PATH} -C ${OUTPUT_PATH}/darwin .
tar -I pigz -cf ${LNX_RLS_PATH} -C ${OUTPUT_PATH}/linux .
7z a ${WIN_RLS_PATH} ${WIN_BIN_PATH} ${WIN_SHA_PATH} > /dev/null

# Done
echo
echo " -> All Done!"