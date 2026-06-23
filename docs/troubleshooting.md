# Troubleshooting

## Godot not found

Set `GODOT_PATH` or pass `--godot`.

```bash
godot-lsp-go --launch --godot /path/to/godot
```

## Project not found

Run from inside a Godot project or pass `--project`.

```bash
godot-lsp-go --launch --project /path/to/project
```

## No diagnostics or completion

Make sure Godot is open with the target project and its LSP server is listening. Try `--debug` and inspect stderr/log file.

## Windows paths look wrong

URI normalization is enabled by default. If needed, disable it with `--no-normalize-uris`.

## OpenCode opens files as plaintext

The OpenCode patch is enabled by default. If it causes issues, disable it with `--no-patch-opencode`.

## Linux headless

Use `--launch --headless`. If your environment requires a display, install `xvfb-run` and configure `--xvfb auto`.

## Remote Godot finds the wrong files

If OpenCode and Godot see the same project at different absolute paths, configure `--path-map`.

```bash
godot-lsp-go --path-map /home/alpine/game=C:/Users/Murilo/Godot/game
```

Use SSH tunneling instead of exposing the Godot LSP port directly. See [Remote connections](remote.md).
