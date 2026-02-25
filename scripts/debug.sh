#!/bin/bash

cd "$(dirname "$0")"

SOURCE_DIR=${SOURCE_DIR:-"../"}
ARTIFACT_NAME=${ARTIFACT_NAME:-"server"}
ARTIFACT_OUTPUT_DIR=${ARTIFACT_OUTPUT_DIR:-"../build/debug"}
ARTIFACT_OUTPUT_PATH="$ARTIFACT_OUTPUT_DIR/$ARTIFACT_NAME"

echo "building debug executable $ARTIFACT_OUTPUT_PATH..."

if [ ! -d "$ARTIFACT_OUTPUT_DIR" ]; then
  mkdir -p "$ARTIFACT_OUTPUT_DIR"
fi

go build -o "$ARTIFACT_OUTPUT_PATH" "$SOURCE_DIR/cmd/$ARTIFACT_NAME/main.go"
if [ $? -ne 0 ]; then
  echo "could not execute 'go build' command"
  exit 1
fi
echo "build successful: ${ARTIFACT_OUTPUT_PATH}"

cd "$ARTIFACT_OUTPUT_DIR"
./"$ARTIFACT_NAME"
