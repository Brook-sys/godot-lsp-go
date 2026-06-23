# Development

## Test

```bash
go test ./...
go vet ./...
gofmt -s -w .
```

## Build

```bash
go build ./cmd/godot-lsp-go
```

## Release

Create a Git tag like `v0.1.0`. The release workflow builds platform archives and checksums.
