package lsp

import (
	"fmt"
	"io"
	"sync"
)

type Writer struct {
	w  io.Writer
	mu sync.Mutex
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteMessage(m Message) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(m.Body))
	if _, err := io.WriteString(w.w, header); err != nil {
		return err
	}
	if _, err := w.w.Write(m.Body); err != nil {
		return err
	}
	return nil
}
