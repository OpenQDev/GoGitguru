#!/bin/bash

cd ./sql/schema
goose postgres "postgres://postgres:@localhost:5432/postgres?sslmode=disable" down
