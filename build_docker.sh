#!/usr/bin/env bash
set -e

# This script builds the Docker image for the leaderboard-api project using the Go SDK Docker image.
# Usage: ./build_docker.sh [tag]

TAG=${1:-latest}
full_tag="leaderboard-api:$TAG"

echo "Building Docker image with tag: $full_tag"
docker buildx build -t "$full_tag" .
echo "Docker image built successfully: $full_tag"
