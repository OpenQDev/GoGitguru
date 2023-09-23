# Gitguru

## Testing

### Test

```bash
cd ./dir/where/tests/live
go test -v
```

`-v` - Run in verbose mode, showing all tests as they run

### Test with Coverage

```bash
go test -cover -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

`-cover` - Provide coverage percentage to stdout
`-coverprofile=coverage.out` - Provide a coverage report to `coverage.out` that can be opened in the browser