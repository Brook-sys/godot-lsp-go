# Claude Code

Use `godot-lsp-go` as a stdio language server command for GDScript.

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

If Godot is not running, either start it manually or add `--launch`.
