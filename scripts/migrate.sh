#!/bin/bash

set -a
source .env
set +a

GOOSE_DBSTRING="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
export GOOSE_DBSTRING
export GOOSE_DRIVER=postgres
export GOOSE_DIR=./migrations

ACTION=${1:-"up"}

if [ -z "$GOOSE_DBSTRING" ]; then
    echo "Error: GOOSE_DBSTRING is not set. Please set the DATABASE_URL in .env file."
    exit 1
fi
if [ -z "$GOOSE_DRIVER" ]; then
    echo "Error: GOOSE_DRIVER could not be determined from DATABASE_URL."
    exit 1
fi
if [ -z "$GOOSE_DIR" ]; then
    echo "Error: GOOSE_DIR is not set. Please set the migrations directory."
    exit 1
fi

case "$ACTION" in
up)
    echo "Running migrations UP..."
    goose -dir "$GOOSE_DIR" up "$GOOSE_DRIVER" "$GOOSE_DBSTRING"
    ;;
down)
    echo "Rolling back last migration..."
    goose -dir "$GOOSE_DIR" down "$GOOSE_DRIVER" "$GOOSE_DBSTRING"
    ;;
reset)
    echo "Rolling back all migrations..."
    goose -dir "$GOOSE_DIR" reset "$GOOSE_DRIVER" "$GOOSE_DBSTRING"
    ;;
*)
    echo "Unknown action: $ACTION"
    echo "Usage: ./migrate.sh [up|down|reset]"
    exit 1
    ;;
esac
