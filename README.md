# godot-lsp-go

A robust single-binary stdio-to-TCP bridge for Godot's GDScript Language Server.

Godot exposes its GDScript LSP over TCP, while tools like OpenCode, Cursor, Claude Code, Neovim and other LSP clients commonly expect stdio. This bridge converts between both protocols without requiring Node.js, npm or npx.

## Features

- stdio ↔ TCP bridge for Godot LSP
- single Go binary
- port discovery: `6005`, `6007`, `6008`
- reconnect loop
- initialize response ordering guard for Godot's non-standard notifications
- Windows file URI normalization
- OpenCode `plaintext` → `gdscript` patch for `.gd` and `.gdshader`
- optional Godot auto-launch with `--launch`
- remote Godot support with bidirectional `--path-map`
- stderr/file logging only, never stdout

## Install

Download a release binary from [GitHub Releases](https://github.com/Brook-sys/godot-lsp-go/releases/latest) and place it in your `PATH`.

Linux x86_64:

```bash
curl -L -o godot-lsp-go.tar.gz https://github.com/Brook-sys/godot-lsp-go/releases/latest/download/godot-lsp-go_linux_amd64.tar.gz
tar -xzf godot-lsp-go.tar.gz
chmod +x godot-lsp-go
sudo mv godot-lsp-go /usr/local/bin/
```

Windows users should download `godot-lsp-go_windows_amd64.zip`, extract `godot-lsp-go.exe`, and add the folder to `PATH`.

See [Installation](docs/installation.md) for all platforms.

## Usage

Verify installation:

```bash
godot-lsp-go --version
```


Default mode does not launch Godot automatically:

```bash
godot-lsp-go
```

Start Godot yourself, open your project, then start your LSP client.

To allow the bridge to launch Godot:

```bash
godot-lsp-go --launch
```

## OpenCode

```json
{
  "lsp": {
    "gdscript": {
      "command": ["godot-lsp-go"],
      "extensions": [".gd", ".gdshader"]
    }
  }
}
```

With auto-launch:

```json
{
  "lsp": {
    "gdscript": {
      "command": ["godot-lsp-go", "--launch"],
      "extensions": [".gd", ".gdshader"]
    }
  }
}
```

## Important default

`godot-lsp-go` uses the conservative default: it connects/reconnects but does not spawn Godot unless `--launch` is provided.

## Documentation

- [Installation](docs/installation.md)
- [Configuration](docs/configuration.md)
- [OpenCode](docs/opencode.md)
- [Godot launch](docs/godot-launch.md)
- [Remote connections](docs/remote.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Architecture](docs/architecture.md)
- [Real Godot E2E tests](docs/e2e.md)
- [Development](docs/development.md)

## License

MIT
