# Architecture

The bridge is split into focused packages:

- `config`: flags, env vars and defaults
- `logging`: stderr/file logging only
- `lsp`: Content-Length reader/writer
- `rewriter`: URI, plain path mapping and OpenCode compatibility rewrites
- `session`: initialize ordering guard
- `connector`: TCP connection and port discovery
- `godot`: project/executable discovery and launching
- `bridge`: stdio ↔ TCP orchestration

Stdout is reserved for LSP messages. Logs must only use stderr or a log file.
