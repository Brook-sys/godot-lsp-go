# Real Godot E2E tests

The default test suite does not require Godot:

```bash
go test ./...
```

To run the real Godot compatibility test, provide a Godot 4.x executable with GDScript LSP support:

```bash
GODOT_E2E_PATH=/path/to/godot go test ./test/e2e -run TestGodotRealLSP -count=1 -v
```

Recommended target versions:

- Godot `4.7-stable`
- latest Godot `4.6.x-stable`

The test creates a temporary Godot project, launches Godot headless, starts `godot-lsp-go`, sends real LSP messages and verifies responses from the real Godot LSP server.

No Godot binary is committed to this repository.

## Linux notes

Official Godot Linux binaries require glibc. They do not run on musl-based distributions such as Alpine without compatibility layers. Run real E2E tests on Ubuntu/Debian/Fedora/Arch or another glibc environment.
