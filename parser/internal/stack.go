package internal

import (
	"slices"

	gr "github.com/PlayerR9/SlParser/grammar"
	dba "github.com/PlayerR9/go-debug/assert"
)

type Stack[T gr.TokenTyper] struct {
	elems  []*gr.Token[T]
	popped []*gr.Token[T]
}

func (s *Stack[T]) Pop() (*gr.Token[T], bool) {
	dba.AssertNotNil(s, "s")

	if len(s.elems) == 0 {
		return nil, false
	}

	top := s.elems[len(s.elems)-1]
	s.elems = s.elems[:len(s.elems)-1]

	s.popped = append(s.popped, top)

	return top, true
}

func (s *Stack[T]) Push(tk *gr.Token[T]) {
	dba.AssertNotNil(s, "s")
	dba.AssertNotNil(tk, "tk")

	s.elems = append(s.elems, tk)
}

func (s Stack[T]) Popped() []*gr.Token[T] {
	popped := make([]*gr.Token[T], len(s.popped))
	copy(popped, s.popped)

	slices.Reverse(popped)

	return popped
}

func (s *Stack[T]) Accept() {
	dba.AssertNotNil(s, "s")

	s.popped = s.popped[:0]
}

func (s *Stack[T]) Refuse() {
	dba.AssertNotNil(s, "s")

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
