package highlighter

import (
	"github.com/PlayerR9/go-evals/common"
	"github.com/PlayerR9/mygo-lib/colors"
)

type Builder struct {
	table        map[string]*colors.Style
	defaultStyle *colors.Style
}

func (b *Builder) Register(type_ string, style *colors.Style) error {
	if style == nil {
		return nil
	} else if b == nil {
		return common.ErrNilReceiver
	}

	if b.table == nil {
		b.table = make(map[string]*colors.Style)
	}

	b.table[type_] = style

	return nil
}

func (b Builder) Build() Highlighter {
	h := Highlighter{
		table: make(map[string]*colors.Style, len(b.table)),
	}

	if b.defaultStyle == nil {
		h.defaultStyle = colors.DefaultStyle
	} else {
		h.defaultStyle = b.defaultStyle
	}

	for k, v := range b.table {
		h.table[k] = v
	}

	return h
}

func (b *Builder) Reset() {
	if b == nil {
		return
	}

	if len(b.table) > 0 {
		clear(b.table)
		b.table = nil
	}

	b.defaultStyle = nil
}
