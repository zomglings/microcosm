#!/bin/sh

DOCKER_IMAGE=${1:-fuzzyfrog/microcosm}

echo "Pushing $DOCKER_IMAGE to Docker Hub..."

docker login -u "$DOCKER_USERNAME" --pasword "$DOCKER_PASSWORD"

docker push $DOCKER_IMAGE

echo "Done!"
