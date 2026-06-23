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
- stderr/file logging only, never stdout

## Install

Download a release binary from GitHub Releases and place it in your `PATH`.

## Usage

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

- [Configuration](docs/configuration.md)
- [OpenCode](docs/opencode.md)
- [Godot launch](docs/godot-launch.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Architecture](docs/architecture.md)
- [Development](docs/development.md)

## License

MIT
