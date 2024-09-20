package backup

type Pair[T any, S any] struct {
	History *History[T]
	Subject S
}

func NewPair[T any, S any](history *History[T], subject S) *Pair[T, S] {
	if history == nil {
		history = &History[T]{}
	}

	return &Pair[T, S]{
		History: history,
		Subject: subject,
	}
}

type Stack[E any, S any] struct {
	elems []*Pair[E, S]
}

func (s *Stack[E, S]) Push(pair *Pair[E, S]) {
	if s == nil || pair == nil {
		return
	}

	s.elems = append(s.elems, pair)
}

func (s *Stack[E, S]) Pop() (*Pair[E, S], bool) {
	if s == nil || len(s.elems) == 0 {
		return nil, false
	}

	top := s.elems[len(s.elems)-1]
	s.elems = s.elems[:len(s.elems)-1]

	return top, true
}
