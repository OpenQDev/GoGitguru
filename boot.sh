#!/bin/bash

docker stop postgres
docker rm postgres

docker run --name postgres -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres

# Wait for PostgreSQL to become accessible
until PGPASSWORD= psql -h "localhost" -U "postgres" -p 5432 -c '\q'; do
	echo "Postgres is unavailable - sleeping"
	sleep 1
done

echo "Postgres is up - executing command"

cd ./sql/schema
goose postgres "postgres://postgres:@localhost:5432/postgres?sslmode=disable" up
