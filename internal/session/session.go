package session

import (
	"bytes"
	"encoding/json"

	"github.com/code-xhyun/godot-lsp-go/internal/lsp"
)

type Guard struct {
	waiting bool
	id      any
	buffer  []lsp.Message
	ready   bool
}

func (g *Guard) TrackClientMessage(msg lsp.Message) {
	var m map[string]any
	if json.Unmarshal(msg.Body, &m) != nil {
		return
	}
	if method, _ := m["method"].(string); method == "initialize" {
		if id, ok := m["id"]; ok {
			g.waiting = true
			g.ready = false
			g.id = id
			g.buffer = nil
		}
	}
}

func (g *Guard) HandleServerMessage(msg lsp.Message) []lsp.Message {
	if !g.waiting {
		return []lsp.Message{msg}
	}
	var m map[string]any
	if json.Unmarshal(msg.Body, &m) != nil {
		return []lsp.Message{msg}
	}
	if id, ok := m["id"]; ok && sameID(id, g.id) && m["result"] != nil {
		out := []lsp.Message{msg}
		out = append(out, g.buffer...)
		g.waiting = false
		g.ready = true
		g.buffer = nil
		return out
	}
	if _, hasMethod := m["method"]; hasMethod {
		if _, hasID := m["id"]; !hasID {
			g.buffer = append(g.buffer, msg)
			return nil
		}
	}
	return []lsp.Message{msg}
}

func (g *Guard) Initialized() bool {
	return g.ready
}

func (g *Guard) ResetConnection() {
	g.waiting = false
	g.id = nil
	g.buffer = nil
}

func sameID(a, b any) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return bytes.Equal(aj, bj)
}
