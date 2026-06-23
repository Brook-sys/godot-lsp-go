# Remote connections

`godot-lsp-go` supports translating file paths bidirectionally, allowing you to run OpenCode (or any LSP client) on one machine while Godot is running on another machine.

This solves the problem where the LSP client expects a file at `/home/alpine/game` but the Godot LSP server expects it at `C:/Users/Murilo/Godot/game`.

## How it works

The `--path-map` flag intercepts JSON-RPC messages and translates:
- `file:///` URIs
- The `rootPath` initialize parameter
- Plain path fields used by Godot (like `gdscript_client/changeWorkspace`)

Translations are applied bidirectionally:
- **Client → Godot:** `/home/alpine/game` is rewritten as `C:/Users/Murilo/Godot/game` before transmission.
- **Godot → Client:** `C:/Users/Murilo/Godot/game` is rewritten back as `/home/alpine/game` before transmission.

## Example: Linux to Windows

Godot is running on Windows:
```text
C:/Users/Murilo/Godot/game
```

OpenCode is running on an Alpine Linux machine:
```text
/home/alpine/game
```

**Step 1:** Ensure your project files are synchronized between the two machines. You can use Syncthing, SSHFS, NFS, rsync, etc.

**Step 2:** Open a secure SSH tunnel from Alpine to Windows, forwarding the LSP port:
```bash
ssh -N -L 6005:127.0.0.1:6005 murilo@windows-pc
```

**Step 3:** Configure OpenCode on Alpine:

```json
{
  "lsp": {
    "gdscript": {
      "command": [
        "godot-lsp-go",
        "--host", "127.0.0.1",
        "--port", "6005",
        "--path-map", "/home/alpine/game=C:/Users/Murilo/Godot/game"
      ],
      "extensions": [".gd", ".gdshader"]
    }
  }
}
```

## Security Warning

**Never expose port 6005 publicly.** The Godot LSP does not have authentication. Always use an SSH tunnel, Tailscale, VPN, or a secure local network to bridge the connection. 

This is why the example connects to `127.0.0.1` locally via SSH instead of directly connecting to `192.168.x.x`.

## Multiple mappings

You can specify `--path-map` multiple times. The most specific mapping is applied first.

```bash
godot-lsp-go \
  --path-map /home/alpine/game=C:/Game \
  --path-map /home/alpine/game/addons=C:/SharedAddons
```

Or via the environment variable (separated by `;`):

```bash
GODOT_LSP_PATH_MAP="/home/alpine/game=C:/Game;/home/alpine/game/addons=C:/SharedAddons"
```
