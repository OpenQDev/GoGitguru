#!/bin/bash

echo "Removing existing postgres container"
docker stop gitguru-postgres >/dev/null 2>/dev/null
docker rm gitguru-postgres >/dev/null 2>/dev/null

POSTGRES_PORT=${1:-5432}

echo "Checking port $POSTGRES_PORT"
if nc -z localhost $POSTGRES_PORT; then
  echo "Something is already running on port $POSTGRES_PORT."
  echo "Please stop it or use a different port and try again."$'\n'
  echo "Usage: ./boot.sh <port>"
  exit 1
fi

echo "Starting new container: $(docker run --name gitguru-postgres -d -e POSTGRES_HOST_AUTH_METHOD=trust -p $POSTGRES_PORT:5432 postgres)"

# Wait for PostgreSQL to become accessible
until PGPASSWORD= psql -h "localhost" -U "postgres" -p $POSTGRES_PORT -c '\q' 2>/dev/null; do
	echo "Waiting for postgres"
	sleep 1
done

echo "Postgres is up - running migrations"

cd ./sql/schema
goose postgres "postgres://postgres:@localhost:$POSTGRES_PORT/postgres?sslmode=disable" up

cd ../..

psql -h "localhost" -U "postgres" -p $POSTGRES_PORT -f ./repos.sql

# To install go air, a Go runtime with live reloads on code changes,
# run the following command: go get -u github.com/cosmtrek/air
go run .
