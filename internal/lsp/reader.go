package lsp

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) ReadMessage() (Message, error) {
	var contentLength int
	for {
		line, err := r.r.ReadBytes('\n')
		if err != nil {
			return Message{}, err
		}
		line = bytes.TrimRight(line, "\r\n")
		if len(line) == 0 {
			if contentLength > 0 {
				break
			}
			continue
		}
		parts := bytes.SplitN(line, []byte(": "), 2)
		if len(parts) == 2 && string(parts[0]) == "Content-Length" {
			n, err := strconv.Atoi(string(parts[1]))
			if err == nil {
				contentLength = n
			}
		}
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(r.r, body); err != nil {
		return Message{}, err
	}
	return Message{Body: body}, nil
}
