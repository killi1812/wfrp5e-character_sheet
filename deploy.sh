#!/bin/sh
set -e

# Gather build information
BUILD="prod"
VERSION=$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
if [ -z "$VERSION" ]; then
  VERSION="0.0.0"
fi
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null)
if [ -z "$COMMIT_HASH" ]; then
  COMMIT_HASH="n/a"
fi
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

echo "Building all docker images with build args:"
echo "  BUILD=$BUILD"
echo "  VERSION=$VERSION"
echo "  COMMIT_HASH=$COMMIT_HASH"
echo "  BUILD_TIMESTAMP=$BUILD_TIMESTAMP"

docker compose -f deploy.yaml build \
  --build-arg BUILD="$BUILD" \
  --build-arg VERSION="$VERSION" \
  --build-arg COMMIT_HASH="$COMMIT_HASH" \
  --build-arg BUILD_TIMESTAMP="$BUILD_TIMESTAMP"

echo "Starting 4-container stack (database, server, client, caddy)..."
docker compose -f deploy.yaml up -d

echo "Deployment complete!"
