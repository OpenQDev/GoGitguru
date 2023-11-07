# Gitguru

## Getting Started

```bash
cp .env.sample .env
./boot.sh

# With different port if 5432 is already in use (by existing local postgres)
./boot.sh 5433
```

## Testing

### Test

```bash
cd ./dir/where/tests/live
go test -v
```

`-v` - Run in verbose mode, showing all tests as they run

### Run One Test

```bash
go test -v -run TestUploadTarballToS3
```

### Run All Tests from Root

```bash
go test ./...
```

### Test with Coverage

```bash
go test -cover -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Benchmark Tests

```bash
go test -bench=BenchmarkProcessRepo
```

`-cover` - Provide coverage percentage to stdout
`-coverprofile=coverage.out` - Provide a coverage report to `coverage.out` that can be opened in the browser

## Local Setup Commands

`docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres`

`psql "postgres://postgres:@localhost:5432/gitguru?sslmode=disable"`

`cd sql/schema && goose postgres "postgres://postgres:@localhost:5432/postgres?sslmode=disable" up && cd ~/Gitguru`