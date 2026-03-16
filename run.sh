#!/usr/bin/env bash
# Run the pipeline in Docker. Replaces func.ps1 for Linux/macOS/WSL.
# Usage: ./run.sh <URL>   or   ./run.sh --rebuild <URL>

set -e
IMAGE="video-forge"
OUTPUT="output"

REBUILD=""
URL=""
for arg in "$@"; do
  if [ "$arg" = "--rebuild" ] || [ "$arg" = "-r" ]; then
    REBUILD=1
  else
    URL="$arg"
  fi
done

if [ -z "$URL" ]; then
  echo "Usage: $0 [--rebuild] <URL_OR_SOURCE>"
  echo "  URL: YouTube URL, HTTP URL, file path, or - for stdin"
  exit 1
fi

if [ -n "$REBUILD" ]; then
  echo "Forcing image rebuild..."
  docker rmi "$IMAGE" 2>/dev/null || true
fi

if ! docker image inspect "$IMAGE" &>/dev/null; then
  echo "Building Docker image..."
  docker build -t "$IMAGE" .
fi

mkdir -p "$OUTPUT"/temp
echo "Running pipeline..."
docker run --rm \
  -v "$(pwd)/$OUTPUT:/app/output" \
  --add-host=host.docker.internal:host-gateway \
  "$IMAGE" "$URL"
