package e2e

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGodotRealLSP(t *testing.T) {
	godotPath := os.Getenv("GODOT_E2E_PATH")
	if godotPath == "" {
		t.Skip("set GODOT_E2E_PATH to run real Godot E2E test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	project := t.TempDir()
	mustWrite(t, filepath.Join(project, "project.godot"), `config_version=5

[application]
config/name="godot-lsp-go-e2e"
`)
	mustWrite(t, filepath.Join(project, "main.gd"), `extends Node

func answer() -> int:
	return 42
`)

	port := freePort(t)
	godotCmd := exec.CommandContext(ctx, godotPath, "--editor", "--headless", "--display-driver", "headless", "--audio-driver", "Dummy", "--lsp-port", strconv.Itoa(port), "--path", project)
	godotCmd.Stdout = io.Discard
	godotCmd.Stderr = io.Discard
	if err := godotCmd.Start(); err != nil {
		t.Fatalf("start godot: %v", err)
	}
	defer func() { _ = godotCmd.Process.Kill() }()

	waitPort(t, ctx, "127.0.0.1", port)

	bridgeCmd := exec.CommandContext(ctx, "go", "run", "./cmd/godot-lsp-go", "--port", strconv.Itoa(port), "--no-reconnect", "--quiet")
	bridgeCmd.Dir = repoRoot(t)
	stdin, err := bridgeCmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := bridgeCmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := bridgeCmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := bridgeCmd.Start(); err != nil {
		t.Fatalf("start bridge: %v", err)
	}
	defer func() { _ = bridgeCmd.Process.Kill() }()
	go io.Copy(io.Discard, stderr)

	reader := bufio.NewReader(stdout)
	writeLSP(t, stdin, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"processId":    nil,
			"rootUri":      fileURI(project),
			"capabilities": map[string]any{},
		},
	})

	msg := readUntilID(t, reader, 1, 20*time.Second)
	if msg["result"] == nil {
		t.Fatalf("initialize response missing result: %#v", msg)
	}

	writeLSP(t, stdin, map[string]any{"jsonrpc": "2.0", "method": "initialized", "params": map[string]any{}})
	writeLSP(t, stdin, map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/didOpen",
		"params": map[string]any{
			"textDocument": map[string]any{
				"uri":        fileURI(filepath.Join(project, "main.gd")),
				"languageId": "plaintext",
				"version":    1,
				"text":       "extends Node\n\nfunc answer() -> int:\n\treturn 42\n",
			},
		},
	})
	writeLSP(t, stdin, map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "textDocument/documentSymbol",
		"params": map[string]any{
			"textDocument": map[string]any{"uri": fileURI(filepath.Join(project, "main.gd"))},
		},
	})

	msg = readUntilID(t, reader, 2, 20*time.Second)
	if _, ok := msg["result"]; !ok {
		t.Fatalf("documentSymbol response missing result: %#v", msg)
	}
}

func writeLSP(t *testing.T, w io.Writer, v any) {
	t.Helper()
	body, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(body), body); err != nil {
		t.Fatal(err)
	}
}

func readUntilID(t *testing.T, r *bufio.Reader, id float64, timeout time.Duration) map[string]any {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		msg := readLSP(t, r)
		if got, ok := msg["id"].(float64); ok && got == id {
			return msg
		}
	}
	t.Fatalf("timeout waiting for id %v", id)
	return nil
}

func readLSP(t *testing.T, r *bufio.Reader) map[string]any {
	t.Helper()
	length := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			parsed, err := strconv.Atoi(value)
			if err != nil {
				t.Fatal(err)
			}
			length = parsed
		}
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(r, body); err != nil {
		t.Fatal(err)
	}
	var msg map[string]any
	if err := json.Unmarshal(body, &msg); err != nil {
		t.Fatalf("json parse %q: %v", body, err)
	}
	return msg
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func waitPort(t *testing.T, ctx context.Context, host string, port int) {
	t.Helper()
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting for %s", addr)
		default:
		}
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("repo root not found")
		}
		wd = parent
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fileURI(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "file://" + path
}

var _ = bytes.Buffer{}
