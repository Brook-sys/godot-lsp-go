package bridge

import "github.com/code-xhyun/godot-lsp-go/internal/lsp"

type Queue struct {
	items []lsp.Message
	max   int
}

func NewQueue(max int) *Queue {
	return &Queue{max: max}
}

func (q *Queue) Push(m lsp.Message) {
	if q.max > 0 && len(q.items) >= q.max {
		q.items = q.items[1:]
	}
	q.items = append(q.items, m)
}

func (q *Queue) Drain() []lsp.Message {
	items := q.items
	q.items = nil
	return items
}

func (q *Queue) Clear() {
	q.items = nil
}

func (q *Queue) Len() int {
	return len(q.items)
}
