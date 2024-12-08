package trav

/////////////////////////////////////////////////////////

type Info interface {
}

type Node interface {
	IsNil() bool
}

type Pair[N Node] struct {
	node        N
	info        Info
	is_critical bool
}

func NewPair[N Node](node N, info Info, is_critical bool) Pair[N] {
	return Pair[N]{
		node:        node,
		info:        info,
		is_critical: is_critical,
	}
}
