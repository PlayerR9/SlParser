package highlighter

import (
	"unicode/utf8"

	"github.com/PlayerR9/go-evals/common"
)

// TODO: Delete this once mygo-lib is updated.

// Encode appends the UTF-8 encoding of a slice of runes to the provided byte slice.
//
// Parameters:
//   - data: A pointer to a byte slice where the encoded UTF-8 bytes will be appended.
//   - chars: A slice of runes to encode.
//
// Returns:
//   - error: An error if the provided data parameter is nil.
func Encode(data *[]byte, chars []rune) error {
	if data == nil {
		return common.NewErrNilParam("data")
	}

	for _, c := range chars {
		*data = utf8.AppendRune(*data, c)
	}

	return nil
}
