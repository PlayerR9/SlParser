package internal

import (
	"testing"
)

func TestMakeLiteral(t *testing.T) {
	const (
		Expected string = "TtNewLine"
	)

	literal, err := MakeLiteral(TerminalTk, "NEW_LINE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if literal != Expected {
		t.Fatalf("expected %q, got %q instead", Expected, literal)
	}
}
