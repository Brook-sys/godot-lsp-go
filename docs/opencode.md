# OpenCode

## Manual Godot startup

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

## Auto-launch Godot

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

## Custom Godot binary

```json
{
  "lsp": {
    "gdscript": {
      "command": ["godot-lsp-go", "--launch", "--godot", "/path/to/godot"],
      "extensions": [".gd", ".gdshader"]
    }
  }
}
```
