package lexer

import (
	"errors"
	"io"
	"unicode/utf8"

	gcers "github.com/PlayerR9/go-commons/errors"
)

// Stream is a stream of runes.
type Stream struct {
	// data is the current input data.
	data []byte

	// prev is the previously read rune.
	prev *rune

	// prev_size is the size of the previously read rune.
	prev_size int

	// has_unread is true if the stream has unread data, false otherwise.
	has_unread bool
}

// ReadRune implements io.RuneScanner.
func (s *Stream) ReadRune() (rune, int, error) {
	if s == nil {
		return 0, 0, gcers.NilReceiver
	}

	if s.has_unread {
		s.has_unread = false

		c := *s.prev

		return c, s.prev_size, nil
	}

	if len(s.data) == 0 {
		return 0, 0, io.EOF
	}

	c, size := utf8.DecodeRune(s.data)
	if c == utf8.RuneError {
		return 0, 0, errors.New("invalid UTF-8")
	}

	s.data = s.data[size:]

	s.prev = &c
	s.prev_size = size

	return c, size, nil
}

// UnreadRune implements io.RuneScanner.
func (s *Stream) UnreadRune() error {
	if s == nil {
		return gcers.NilReceiver
	}

	if s.prev == nil {
		return errors.New("nothing to unread")
	}

	s.has_unread = true

	return nil
}

// NewStream creates a new stream.
//
// Returns:
//   - *Stream: the new stream. Never returns nil.
func NewStream() *Stream {
	return &Stream{}
}

// FromString initializes the stream from a string.
//
// Returns:
//   - *Stream: the new stream. Nil only if the receiver is nil.
func (s *Stream) FromString(data string) *Stream {
	if s == nil {
		return nil
	}

	s.data = []byte(data)

	return s
}

// FromBytes initializes the stream from a byte slice.
//
// Returns:
//   - *Stream: the new stream. Nil only if the receiver is nil.
func (s *Stream) FromBytes(data []byte) *Stream {
	if s == nil {
		return nil
	}

	s.data = data

	return s
}
