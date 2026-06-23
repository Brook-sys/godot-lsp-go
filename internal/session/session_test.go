package session

import (
	"strings"
	"testing"

	"github.com/Brook-sys/godot-lsp-go/internal/lsp"
)

func TestInitializeGuardBuffersNotifications(t *testing.T) {
	var g Guard
	g.TrackClientMessage(lsp.Message{Body: []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`)})
	out := g.HandleServerMessage(lsp.Message{Body: []byte(`{"jsonrpc":"2.0","method":"window/logMessage","params":{}}`)})
	if len(out) != 0 {
		t.Fatalf("expected buffered notification")
	}
	out = g.HandleServerMessage(lsp.Message{Body: []byte(`{"jsonrpc":"2.0","id":1,"result":{}}`)})
	if len(out) != 2 {
		t.Fatalf("expected response plus buffered notification, got %d", len(out))
	}
	if !strings.Contains(string(out[0].Body), `"result"`) {
		t.Fatalf("initialize response should be first")
	}
}
