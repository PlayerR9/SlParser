package lexer

import (
	"testing"
	"unicode"

	emtch "github.com/PlayerR9/go-evals/matcher"
)

func TestSuccessMatching(t *testing.T) {
	type args struct {
		Stream   []rune
		Matcher  emtch.Matcher[rune]
		Expected string
	}

	tests := []args{
		// Test newline
		{
			Stream:   []rune{'\r', '\n'},
			Matcher:  Newline(),
			Expected: "\r\n",
		},
		{
			Stream:   []rune{'\n'},
			Matcher:  Newline(),
			Expected: "\n",
		},

		// Test single rune
		{
			Stream:   []rune{'a'},
			Matcher:  One('a'),
			Expected: "a",
		},

		// Test match group
		{
			Stream:   []rune{'a'},
			Matcher:  Predicate("letter", unicode.IsLetter),
			Expected: "a",
		},

		// Test literal
		{
			Stream:   []rune{'f', 'o', 'o'},
			Matcher:  Literal("foo"),
			Expected: "foo",
		},

		// Test match many
		{
			Stream:   []rune{'a', 'a', 'a'},
			Matcher:  Many(One('a')),
			Expected: "aaa",
		},

		// Test match sequence
		{
			Stream:   []rune{'a', 'b', 'c'},
			Matcher:  Sequence(One('a'), One('b'), One('c')),
			Expected: "abc",
		},
	}

	for i, arg := range tests {
		valid := true

		for _, c := range arg.Stream {
			err := arg.Matcher.Match(c)
			if err != nil {
				t.Errorf("test %d, want no error, got %v", i, err)
				valid = false
				break
			}
		}

		if !valid {
			continue
		}

		err := arg.Matcher.Close()
		if err != nil {
			t.Errorf("test %d, want no error, got %v", i, err)
		}

		data := string(arg.Matcher.Matched())
		if data != arg.Expected {
			t.Errorf("test %d, want %s, got %s", i, arg.Expected, data)
		}
	}
}
