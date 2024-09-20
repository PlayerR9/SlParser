package ast

import (
	gr "github.com/PlayerR9/SlParser/grammar"
)

type Builder[N interface {
	AddChildren(children []N)
}, T gr.TokenTyper] struct {
	table map[T]ToAstFunc[N, T]
	// make_fake_node func(root *gr.ParseTree[T]) N
}

func NewBuilder[N interface {
	AddChildren(children []N)
}, T gr.TokenTyper]() Builder[N, T] {
	return Builder[N, T]{
		table: make(map[T]ToAstFunc[N, T]),
		// make_fake_node: nil,
	}
}

// func (b *Builder[N, T]) SetMakeFakeNode(make_fake_node func(root *gr.ParseTree[T]) N) {
// 	if b == nil {
// 		return
// 	}

// 	b.make_fake_node = make_fake_node
// }

/*
func (am Builder[N]) Convert(root *gr.Token[T]) (N, error) {
	if root == nil {
		return *new(N), util.NewErrNilParameter("root")
	}

	type_ := root.Type

	var node N
	var err error

	fn, ok := am.table[type_]
	if !ok {
		err = fmt.Errorf("type is not registered")
	} else {
		node, err = fn(root)
	}

	if err != nil {
		node = TransformFakeNode(root)
		err = NewErrIn(type_, err)
	}

	return node, nil
} */

/* func (am *Builder) Forward(type_ gr.TokenType) {
	if am == nil {
		return
	}

	fn := func(tk *gr.Token[T]) (*Node, error) {
		children := tk.Children

		if len(children) != 1 {
			return nil, fmt.Errorf("expected one children, got %d instead", len(children))
		}

		type_ := children[0].Type

		fn, ok := am.table[type_]
		if !ok {
			return nil, fmt.Errorf("invalid token type %q", type_.String())
		}

		node, err := fn(children[0])
		return node, err
	}

	am.table[type_] = fn
} */

func (am *Builder[N, T]) Register(type_ T, fn ToAstFunc[N, T]) {
	if am == nil || fn == nil {
		return
	}

	am.table[type_] = fn
}

func (am Builder[N, T]) Build() *AstMaker[N, T] {
	var table map[T]ToAstFunc[N, T]

	if len(am.table) > 0 {
		table = make(map[T]ToAstFunc[N, T], len(am.table))
		for k, v := range am.table {
			table[k] = v
		}
	}

	// fn := am.make_fake_node

	return &AstMaker[N, T]{
		table: am.table,
		// make_fake_node: fn,
	}
}

func (am *Builder[N, T]) Reset() {
	if am == nil {
		return
	}

	if len(am.table) > 0 {
		for k := range am.table {
			am.table[k] = nil
			delete(am.table, k)
		}

		am.table = make(map[T]ToAstFunc[N, T])
	}

	// am.make_fake_node = nil
}
