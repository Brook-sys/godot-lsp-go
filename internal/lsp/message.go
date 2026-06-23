package lsp

import "encoding/json"

type Message struct {
	Body []byte
}

func (m Message) JSON(v any) error {
	return json.Unmarshal(m.Body, v)
}

func NewMessage(v any) (Message, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return Message{}, err
	}
	return Message{Body: b}, nil
}
