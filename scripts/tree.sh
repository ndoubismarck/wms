#!/bin/bash

cd "$(dirname "$0")"

SOURCE_DIR="$(realpath "../")"

cd "$SOURCE_DIR"

tree -L 10 -I 'public|node_modules|vendor|dist|build|.dev|.git'