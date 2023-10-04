#!/bin/bash

echo "Removing existing postgres container"
docker stop gitguru-postgres > /dev/null 2>/dev/null
docker rm gitguru-postgres > /dev/null 2>/dev/null

echo "Starting new container: $(docker run --name gitguru-postgres -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres)"

# Wait for PostgreSQL to become accessible
until PGPASSWORD= psql -h "localhost" -U "postgres" -p 5432 -c '\q' 2>/dev/null; do
	echo "Waiting for postgres"
	sleep 1
done

echo "Postgres is up - running migrations"

cd ./sql/schema
goose postgres "postgres://postgres:@localhost:5432/postgres?sslmode=disable" up

cd ../..

psql -h "localhost" -U "postgres" -p 5432 -f ./repos.sql

go run .
