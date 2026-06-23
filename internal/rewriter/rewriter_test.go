package rewriter

import (
	"strings"
	"testing"
)

func TestNormalizeFileURI(t *testing.T) {
	got := NormalizeFileURI(`file://C:\Project\script.gd`)
	want := "file:///C:/Project/script.gd"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestPatchOpenCodePlaintext(t *testing.T) {
	in := []byte(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///tmp/a.gd","languageId":"plaintext"}}}`)
	out := Rewrite(in, Options{PatchOpenCode: true, Direction: ClientToGodot})
	if !strings.Contains(string(out), `"languageId":"gdscript"`) {
		t.Fatalf("expected gdscript patch, got %s", out)
	}
}

func TestPatchOpenCodeDoesNotPatchTxt(t *testing.T) {
	in := []byte(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///tmp/a.txt","languageId":"plaintext"}}}`)
	out := Rewrite(in, Options{PatchOpenCode: true, Direction: ClientToGodot})
	if string(out) != string(in) {
		t.Fatalf("expected unchanged, got %s", out)
	}
}
