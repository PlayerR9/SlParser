package lexer

import (
	"errors"
	"strings"

	gr "github.com/PlayerR9/SlParser/grammar"
)

// LexOption is a lexer option.
//
// Parameters:
//   - options: the lexer options.
type LexOption func(*LexOptions)

// WithLexMany lexes one or more fragments.
//
// Parameters:
//   - should_lex_many: true if the lexer should lex one or more fragments, false otherwise.
//
// Returns:
//   - LexOption: the lexer option.
func WithLexMany(should_lex_many bool) LexOption {
	return func(options *LexOptions) {
		options.lex_many = should_lex_many
	}
}

// WithAllowOptional allows optional fragments.
//
// Parameters:
//   - allow_optional: true if the lexer should allow optional fragments, false otherwise.
//
// Returns:
//   - LexOption: the lexer option.
func WithAllowOptional(allow_optional bool) LexOption {
	return func(options *LexOptions) {
		options.allow_optional = allow_optional
	}
}

// LexOptions is the lexer options.
type LexOptions struct {
	// lex_many is true if the lexer should lex one or more fragments.
	lex_many bool

	// allow_optional is true if the lexer should allow optional fragments.
	allow_optional bool
}

// FragWithOptions lexes a fragment with options.
//
// Parameters:
//   - frag_fn: the fragment function.
//   - options: the lexer options.
//
// Returns:
//   - LexFragment: the lexer fragment.
//
// By default, the lexer does allow optional fragments and only lexes once.
//
// If 'frag_fn' is nil, then a function that returns an error is returned.
func FragWithOptions[T gr.TokenTyper](frag_fn LexFragment[T], options ...LexOption) LexFragment[T] {
	if frag_fn == nil {
		return func(lexer *Lexer[T]) (string, error) {
			return "", errors.New("no fragment function provided")
		}
	}

	settings := LexOptions{
		lex_many:       false,
		allow_optional: true,
	}

	for _, opt := range options {
		opt(&settings)
	}

	var fn LexFragment[T]

	if settings.lex_many {
		if settings.allow_optional {
			fn = func(lexer *Lexer[T]) (string, error) {
				str, err := frag_fn(lexer)
				if err != nil {
					return "", err
				}

				var builder strings.Builder

				builder.WriteString(str)

				for {
					str, err := frag_fn(lexer)
					if err == NotFound {
						break
					} else if err != nil {
						return "", err
					}

					builder.WriteString(str)
				}

				return builder.String(), nil
			}
		} else {
			fn = func(lexer *Lexer[T]) (string, error) {
				var builder strings.Builder

				for {
					str, err := frag_fn(lexer)
					if err == NotFound {
						break
					} else if err != nil {
						return "", err
					}

					builder.WriteString(str)
				}

				return builder.String(), nil
			}
		}
	} else {
		if settings.allow_optional {
			fn = func(lexer *Lexer[T]) (string, error) {
				str, err := frag_fn(lexer)
				if err == NotFound {
					return "", nil
				} else if err != nil {
					return "", err
				}

				return str, nil
			}
		} else {
			fn = func(lexer *Lexer[T]) (string, error) {
				str, err := frag_fn(lexer)
				if err != nil {
					return "", err
				}

				return str, nil
			}
		}
	}

	return fn
}
