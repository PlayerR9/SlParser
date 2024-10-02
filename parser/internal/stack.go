package internal

import (
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/go-errors/assert"
)

type Stack[T gr.TokenTyper] struct {
	elems  []*gr.ParseTree[T]
	popped []*gr.ParseTree[T]
}

func (s *Stack[T]) Pop() (*gr.ParseTree[T], bool) {
	assert.NotNil(s, "s")

	if len(s.elems) == 0 {
		return nil, false
	}

	top := s.elems[len(s.elems)-1]
	s.elems = s.elems[:len(s.elems)-1]

	s.popped = append(s.popped, top)

	return top, true
}

func (s *Stack[T]) Push(tk *gr.ParseTree[T]) {
	assert.NotNil(s, "s")
	assert.NotNil(tk, "tk")

	s.elems = append(s.elems, tk)
}

func (s Stack[T]) Popped() []*gr.ParseTree[T] {
	popped := make([]*gr.ParseTree[T], len(s.popped))
	copy(popped, s.popped)

	slices.Reverse(popped)

	return popped
}

func (s *Stack[T]) Accept() {
	assert.NotNil(s, "s")

	s.popped = s.popped[:0]
}

func (s *Stack[T]) Refuse() {
	assert.NotNil(s, "s")

	for len(s.popped) > 0 {
		top := s.popped[len(s.popped)-1]
		s.popped = s.popped[:len(s.popped)-1]

		s.elems = append(s.elems, top)
	}
}

func (s Stack[T]) IsEmpty() bool {
	return len(s.elems) == 0
}

func (s *Stack[T]) Reset() {
	if s == nil {
		return
	}

	if len(s.elems) > 0 {
		for i := 0; i < len(s.elems); i++ {
			s.elems[i] = nil
		}

		s.elems = s.elems[:0]
	}

	if len(s.popped) > 0 {
		for i := 0; i < len(s.popped); i++ {
			s.popped[i] = nil
		}

		s.popped = s.popped[:0]
	}
}
