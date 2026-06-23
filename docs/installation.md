# Installation

Download the correct archive from the latest GitHub Release:

```text
https://github.com/Brook-sys/godot-lsp-go/releases/latest
```

## Linux x86_64

```bash
curl -L -o godot-lsp-go.tar.gz https://github.com/Brook-sys/godot-lsp-go/releases/latest/download/godot-lsp-go_linux_amd64.tar.gz
tar -xzf godot-lsp-go.tar.gz
chmod +x godot-lsp-go
sudo mv godot-lsp-go /usr/local/bin/
godot-lsp-go --version
```

## Linux ARM64

```bash
curl -L -o godot-lsp-go.tar.gz https://github.com/Brook-sys/godot-lsp-go/releases/latest/download/godot-lsp-go_linux_arm64.tar.gz
tar -xzf godot-lsp-go.tar.gz
chmod +x godot-lsp-go
sudo mv godot-lsp-go /usr/local/bin/
godot-lsp-go --version
```

## macOS Intel

```bash
curl -L -o godot-lsp-go.tar.gz https://github.com/Brook-sys/godot-lsp-go/releases/latest/download/godot-lsp-go_darwin_amd64.tar.gz
tar -xzf godot-lsp-go.tar.gz
chmod +x godot-lsp-go
sudo mv godot-lsp-go /usr/local/bin/
godot-lsp-go --version
```

## macOS Apple Silicon

```bash
curl -L -o godot-lsp-go.tar.gz https://github.com/Brook-sys/godot-lsp-go/releases/latest/download/godot-lsp-go_darwin_arm64.tar.gz
tar -xzf godot-lsp-go.tar.gz
chmod +x godot-lsp-go
sudo mv godot-lsp-go /usr/local/bin/
godot-lsp-go --version
```

If macOS blocks the binary because it is unsigned:

```bash
xattr -d com.apple.quarantine /usr/local/bin/godot-lsp-go
```

## Windows x86_64

Download:

```text
godot-lsp-go_windows_amd64.zip
```

Extract `godot-lsp-go.exe` and place it somewhere in your `PATH`, for example:

```text
C:\Tools\godot-lsp-go\godot-lsp-go.exe
```

Then test from PowerShell:

```powershell
godot-lsp-go.exe --version
```

## Windows ARM64

Download:

```text
godot-lsp-go_windows_arm64.zip
```

Extract `godot-lsp-go.exe` and add its folder to your `PATH`.

## Build from source

```bash
git clone https://github.com/Brook-sys/godot-lsp-go.git
cd godot-lsp-go
go build -o godot-lsp-go ./cmd/godot-lsp-go
./godot-lsp-go --version
```

## Verify checksum

Download `checksums.txt` from the same release and run:

```bash
sha256sum -c checksums.txt
```
