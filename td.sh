#!/bin/bash

cd ./sql/schema
while true; do
	goose postgres "postgres://postgres:@localhost:5432/postgres?sslmode=disable" down || break
done
