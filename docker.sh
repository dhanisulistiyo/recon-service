#!/bin/bash

set -e

APP_NAME="reconciliation-service"
IMAGE_NAME="reconciliation-service:latest"
CONTAINER_NAME="reconciliation-service"
PORT=8081

function build_image() {
  echo "docker build -t reconciliation-service:latest .Building Docker image..."
  docker build -t ${IMAGE_NAME} .
  echo "Image built: ${IMAGE_NAME}"
}

function stop_container() {
  if [ "$(docker ps -aq -f name=${CONTAINER_NAME})" ]; then
    echo "Stopping existing container..."
    docker rm -f ${CONTAINER_NAME}
    echo "Old container removed"
  fi
}

function run_container() {
  echo "Running container..."
  docker run -d \
    --name ${CONTAINER_NAME} \
    -p ${PORT}:${PORT} \
    --env-file .env \
    ${IMAGE_NAME}

  echo "Container running at http://localhost:${PORT}"
}

case "$1" in
  build)
    build_image
    ;;
  run)
    stop_container
    run_container
    ;;
  up)
    build_image
    stop_container
    run_container
    ;;
  *)
    echo "Usage:"
    echo "  ./scripts/docker.sh build   # build image"
    echo "  ./scripts/docker.sh run     # run container"
    echo "  ./scripts/docker.sh up      # build + run"
    exit 1
    ;;

  esac