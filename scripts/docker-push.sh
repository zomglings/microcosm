#!/bin/sh

DOCKER_IMAGE=${1:-fuzzyfrog/microcosm}

echo "Pushing $DOCKER_IMAGE to Docker Hub..."

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --pasword-stdin

docker push $DOCKER_IMAGE

echo "Done!"
