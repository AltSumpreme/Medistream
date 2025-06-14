#!/bin/bash

set -a
source ../.env 

set +a

export GOOSE_DBSTRING=$DATABASE_URL
export GOOSE_DRIVER=postgres
export GOOSE_DIR=../migrations

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
echo "Running migrations with Goose..."
goose -dir $GOOSE_DIR up  "$GOOSE_DRIVER" "$GOOSE_DBSTRING" 