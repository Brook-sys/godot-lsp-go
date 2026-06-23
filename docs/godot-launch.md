# Godot launch

`godot-lsp-go` does not launch Godot unless `--launch` is provided.

When launching, it resolves the project in this order:

1. `--project`
2. `GODOT_PROJECT`
3. `OPENCODE_PROJECT_ROOT`
4. current directory, walking upward until `project.godot`

It resolves Godot in this order:

1. `--godot`
2. `GODOT_PATH`
3. `PATH`: `godot`, `godot4`, `godot-editor`, `godot.exe`

Default launch command:

```bash
godot --editor --headless --display-driver headless --audio-driver Dummy --lsp-port 6005 --path /project
```
