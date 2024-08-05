@echo off

SET APP=%1
SET POSTGRES_PORT=5433
IF "%APP%"=="" SET APP="server"

docker inspect gitguru-postgres >nul 2>nul
if errorlevel 1 (
    REM Remove existing postgres container
    docker stop gitguru-postgres >nul 2>nul
    docker rm gitguru-postgres >nul 2>nul

    REM Start new container
    echo Starting new container: 
    docker run --name gitguru-postgres -d -e POSTGRES_HOST_AUTH_METHOD=trust -p %POSTGRES_PORT%:5432 postgres
) else (
    echo Postgres container already exists, you need to start it if it is not running
)
:retry
psql -h "localhost" -U "postgres" -p %POSTGRES_PORT% -c \q >nul 2>nul
if errorlevel 1 goto retry
echo Postgres is up - running migrations


cd .\sql\schema
goose postgres "postgres://postgres:@localhost:%POSTGRES_PORT%/postgres?sslmode=disable" up