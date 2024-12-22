package matcher

import (
	"fmt"
	"io"

	"github.com/PlayerR9/SlParser/mygo-lib/common"
	assert "github.com/PlayerR9/go-verify"
)

type Matcher interface {
	Match(char rune) error
	Matched() []rune
	Close() error
}

func do(matcher Matcher, scanner io.RuneScanner) (bool, error) {
	char, _, err := scanner.ReadRune()
	if err == io.EOF {
		return false, nil
	}

	if err != nil {
		err := fmt.Errorf("unable to read rune: %w", err)
		return false, err
	}

	err = matcher.Match(char)
	if err == nil {
		return true, nil
	}

	assert.Err(scanner.UnreadRune(), "scanner.UnreadRune()")

	if err == ErrMatchDone {
		return false, nil
	}

	err = fmt.Errorf("unable to match: %w", err)
	return false, err
}

func Match(matcher Matcher, scanner io.RuneScanner) ([]rune, error) {
	if matcher == nil {
		err := common.NewErrNilParam("matcher")
		return nil, err
	} else if scanner == nil {
		err := common.NewErrNilParam("scanner")
		return nil, err
	}

	var stop bool
	var err error

	for !stop && err == nil {
		stop, err = do(matcher, scanner)
	}

	if err != nil {
		return nil, err
	}

	err = matcher.Close()
	if err != nil {
		err := fmt.Errorf("unable to close matcher: %w", err)
		return nil, err
	}

	matched := matcher.Matched()

	return matched, nil
}
