package ast

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/mygo-lib/common"
)

func Quote(b *strings.Builder, type_ string) error {
	if b == nil {
		return common.NewErrNilParam("b")
	}

	_, _ = b.WriteString(strconv.Quote(type_))

	return nil
}

type Noder interface {
	IsNil() bool
	String() string
	GetType() string
}
