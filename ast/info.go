package ast

/*
type CheckDataFn func(data string) error

type CheckChildrenFn[N interface {
	Child() iter.Seq[N]
	GetData() string

	Noder
}] func(children []N) error

type NodeInfo[N interface {
	Child() iter.Seq[N]
	GetData() string

	Noder
}] struct {
	check_data     CheckDataFn
	check_children CheckChildrenFn[N]
}

func NewNodeInfo[N interface {
	Child() iter.Seq[N]
	GetData() string

	Noder
}](check_data CheckDataFn, check_children CheckChildrenFn[N]) NodeInfo[N] {
	if check_data == nil {
		check_data = func(data string) error {
			if data != "" {
				return fmt.Errorf("expected an empty string, got %q", data)
			}

			return nil
		}
	}

	if check_children == nil {
		check_children = func(children []N) error {
			if len(children) != 0 {
				str := strconv.Itoa(len(children))

				return fmt.Errorf("expected no children, got %s", str)
			}

			return nil
		}
	}

	return NodeInfo[N]{
		check_data:     check_data,
		check_children: check_children,
	}
}

func (ni NodeInfo[N]) CheckNode(node N) ([]N, error) {
	if node.IsNil() {
		return nil, ErrNilNode
	}

	children := slices.Collect(node.Child())

	err := ni.check_data(node.GetData())
	if err != nil {
		return children, NewErrBadData(err)
	}

	err = ni.check_children(children)
	if err != nil {
		return children, NewErrBadChildren(err)
	}

	return children, nil
}
*/
