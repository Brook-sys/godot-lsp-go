# Cursor

Configure a custom language server command that points to the binary:

```json
{
  "gdscript": {
    "command": ["godot-lsp-go"],
    "filetypes": ["gdscript"]
  }
}
```

Use `--launch` only if you want the bridge to start Godot automatically.
