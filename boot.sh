#!/bin/bash

echo "Removing existing postgres container"
docker stop gitguru-postgres >/dev/null 2>/dev/null
docker rm gitguru-postgres >/dev/null 2>/dev/null

POSTGRES_PORT=5432
APP=${1:-"reposync"}

echo "Starting new container: $(docker run --name gitguru-postgres -d -e POSTGRES_HOST_AUTH_METHOD=trust -p $POSTGRES_PORT:5432 -v ./data:/var/lib/postgresql/data postgres)"

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
air ./$APP
