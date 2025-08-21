#!/bin/bash

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

DATABASE_URL="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable"

USER_ID=$(id -u)
GROUP_ID=$(id -g)

if [ "$1" = "create" ]; then
  shift
  docker run --rm \
    --user $USER_ID:$GROUP_ID \
    -v $(pwd)/migrations:/migrations \
    migrate/migrate \
    create -dir /migrations "$@"
else
  docker run --rm \
    --user $USER_ID:$GROUP_ID \
    -v $(pwd)/migrations:/migrations \
    --network host \
    migrate/migrate \
    -path=/migrations \
    -database "$DATABASE_URL" \
    "$@"
fi
