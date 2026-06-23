# Configuration

Priority order:

```text
flags > environment variables > defaults
```

## Defaults

```text
host: 127.0.0.1
ports: 6005,6007,6008
launch: false
headless: true
reconnect: true
connect-timeout: 2s
startup-timeout: 30s
reconnect-delay: 5s
warmup-delay: 1s
max-buffer-size: 10485760
max-pending-messages: 1000
normalize-uris: true
patch-opencode: true
```

## Flags

```bash
godot-lsp-go --host 127.0.0.1 --ports 6005,6007,6008
godot-lsp-go --port 6008
godot-lsp-go --launch --project /path/to/project --godot /path/to/godot
godot-lsp-go --no-patch-opencode --no-normalize-uris
godot-lsp-go --debug --log-file /tmp/godot-lsp-go.log
```

## Environment variables

```text
GODOT_LSP_HOST
GODOT_LSP_PORT
GODOT_LSP_PORTS
GODOT_LSP_PATH_MAP
GODOT_LSP_BRIDGE_DEBUG
GODOT_LSP_BRIDGE_LOG
GODOT_PATH
GODOT_PROJECT
OPENCODE_PROJECT_ROOT
```

## Remote path mapping

```bash
godot-lsp-go --path-map /client/project=/godot/project
```

`--path-map` can be repeated. See [Remote connections](remote.md).
