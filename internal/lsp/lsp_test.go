package lsp

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	input := "Content-Length: 2\r\n\r\n{}"
	r := NewReader(bytes.NewBufferString(input))
	msg, err := r.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if string(msg.Body) != "{}" {
		t.Errorf("expected {}, got %s", msg.Body)
	}
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	err := w.WriteMessage(Message{Body: []byte("{}")})
	if err != nil {
		t.Fatal(err)
	}
	expected := "Content-Length: 2\r\n\r\n{}"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}
