package pkg

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	gch "github.com/PlayerR9/mygo-lib/runes"
)

func FixUnderscores(name []byte, is_terminal bool) (string, error) {
	chars, err := gch.BytesToUtf8(name)
	if err != nil {
		return "", err
	}

	next_uppercase := true

	var builder strings.Builder

	for _, c := range chars {
		if c == '_' {
			next_uppercase = true

			continue
		}

		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return "", errors.New("name contains invalid characters")
		}

		if next_uppercase {
			builder.WriteRune(unicode.ToUpper(c))
		} else {
			if is_terminal {
				builder.WriteRune(unicode.ToLower(c))
			} else {
				builder.WriteRune(c)
			}
		}

		next_uppercase = false
	}

	if next_uppercase {
		return "", errors.New("name cannot end with an underscore")
	}

	return builder.String(), nil
}

func FixEnumName(name []byte) (string, error) {
	if len(name) == 0 {
		return "", errors.New("name cannot be empty")
	}

	ok := utf8.Valid(name)
	if !ok {
		return "", errors.New("name is not valid utf8")
	}

	literal := string(name)

	// Special case for EOF
	if literal == "EOF" {
		return "EtEOF", nil
	}

	c, _ := utf8.DecodeRune(name)

	if !unicode.IsLetter(c) {
		return "", errors.New("name must start with a letter")
	}

	is_terminal := unicode.IsUpper(c)

	str, err := FixUnderscores(name, is_terminal)
	if err != nil {
		return "", err
	}

	if is_terminal {
		return "Tt" + str, nil
	} else {
		return "Nt" + str, nil
	}
}
