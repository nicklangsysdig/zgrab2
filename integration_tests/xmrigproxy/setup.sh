#!/usr/bin/env bash

echo "xmrigproxy/setup: Tests setup for xmrigproxy"

CONTAINER_TAG="zgrab_xmrigproxy"
CONTAINER_NAME="zgrab_xmrigproxy"

# If the container is already running, use it.
if docker ps --filter "name=$CONTAINER_NAME" | grep -q $CONTAINER_NAME; then
    echo "xmrigproxy/setup: Container $CONTAINER_NAME already running -- nothing to setup"
    exit 0
fi

DOCKER_RUN_FLAGS="--rm --name $CONTAINER_NAME -td"

# If it is not running, try launching it -- on success, use that. 
echo "xmrigproxy/setup: Trying to launch $CONTAINER_NAME..."
if ! docker run $DOCKER_RUN_FLAGS $CONTAINER_TAG; then
    echo "xmrigproxy/setup: Building docker image $CONTAINER_TAG..."
    # If it fails, build it from ./container/Dockerfile
    docker build -t $CONTAINER_TAG ./container
    # Try again
    echo "xmrigproxy/setup: Launching $CONTAINER_NAME..."
    docker run -p 3333:3333 $DOCKER_RUN_FLAGS $CONTAINER_TAG
fi
