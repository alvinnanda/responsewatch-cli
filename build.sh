#!/bin/bash

# Build script for ResponseWatch CLI
# Creates binaries for multiple platforms

set -e

APP_NAME="rwcli"
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
BUILD_DIR="./dist"

# Clean build directory
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

echo "Building $APP_NAME version $VERSION..."

# Build for multiple platforms
platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    output="$BUILD_DIR/${APP_NAME}_${GOOS}_${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output="${output}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X main.version=$VERSION" -o "$output" .
done

echo "Build complete! Binaries in $BUILD_DIR:"
ls -la "$BUILD_DIR"
