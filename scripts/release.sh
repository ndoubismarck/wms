#!/bin/bash

cd "$(dirname "$0")"

GOOS=${GOOS:-"linux"}
GOARCH=${GOARCH:-"amd64"}
SOURCE_DIR=${SOURCE_DIR:-"../"}
ARTIFACT_NAME=${ARTIFACT_NAME:-"server"}
ARTIFACT_FILE_NAME="$ARTIFACT_NAME-$GOOS-$GOARCH"
ARTIFACT_OUTPUT_DIR=${ARTIFACT_OUTPUT_DIR:-"../build/release"}
ARTIFACT_OUTPUT_PATH="${ARTIFACT_OUTPUT_DIR}/${ARTIFACT_FILE_NAME}"
WEB_FRONTEND_BUILD_SCRIPT_PATH="../frontend/web/scripts/build.sh"
DESKTOP_FRONTEND_BUILD_SCRIPT_PATH="../frontend/desktop/scripts/build.sh"

echo "building release executable ${ARTIFACT_OUTPUT_PATH}..."

if [ ! -d "$ARTIFACT_OUTPUT_DIR" ]; then
  mkdir -p "$ARTIFACT_OUTPUT_DIR"
fi

if [ -f "$WEB_FRONTEND_BUILD_SCRIPT_PATH" ]; then
  bash "$WEB_FRONTEND_BUILD_SCRIPT_PATH"
fi

if [ -f "$DESKTOP_FRONTEND_BUILD_SCRIPT_PATH" ]; then
  bash "$DESKTOP_FRONTEND_BUILD_SCRIPT_PATH"
fi

go mod tidy

env GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build \
    -trimpath \
    -gcflags="-trimpath=$(dirname "$SOURCE_DIR"):$GOPATH" \
    -asmflags="-trimpath=$(dirname "$SOURCE_DIR"):$GOPATH" \
    -o "$ARTIFACT_OUTPUT_PATH" "$SOURCE_DIR/cmd/$ARTIFACT_NAME/main.go"
if [ $? -ne 0 ]; then
  echo "could not execute 'go build' command"
  exit 1
fi
echo "build successful: $ARTIFACT_OUTPUT_PATH"
