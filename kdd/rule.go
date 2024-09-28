package kdd

import (
	"fmt"
	"slices"

	gcslc "github.com/PlayerR9/SlParser/PlayerR9/go-commons/slices"
	"github.com/PlayerR9/SlParser/ast"
	"github.com/PlayerR9/SlParser/grammar"
	"github.com/PlayerR9/go-commons/strings"
	gcstr "github.com/PlayerR9/go-commons/strings"
)

var (
	RuleSet map[TokenType][]*Rule
)

func init() {
	RuleSet = make(map[TokenType][]*Rule)
}

func GetRule(lhs TokenType) ([]*Rule, bool) {
	if len(RuleSet) == 0 {
		return nil, false
	}

	vals, ok := RuleSet[lhs]
	return vals, ok
}

type Field struct {
	Types []TokenType
}

func (f *Field) IsNil() bool {
	return f == nil
}

func NewField(types ...TokenType) (*Field, error) {

	if len(types) == 0 {
		return nil, fmt.Errorf("expected at least one type")
	}

	unique := make([]TokenType, 0, len(types))

	for _, t := range types {
		pos, ok := slices.BinarySearch(unique, t)
		if !ok {
			unique = slices.Insert(unique, pos, t)
		}
	}

	return &Field{
		Types: unique,
	}, nil
}

type Rule struct {
	Fields    []*Field
	Expecteds map[int][]NodeType
	IsLhsRule bool
}

func NewRule(lhs TokenType, is_lhs_rule bool, fields ...*Field) (*Rule, error) {
	fields = gcslc.FilterNonNil(fields)
	if len(fields) == 0 {
		return nil, fmt.Errorf("expected at least one field")
	}

	r := &Rule{
		Fields:    fields,
		IsLhsRule: is_lhs_rule,
		Expecteds: make(map[int][]NodeType),
	}

	prev, ok := RuleSet[lhs]
	if !ok {
		prev = []*Rule{r}
	} else {
		prev = append(prev, r)
	}

	RuleSet[lhs] = prev

	return r, nil
}

func (r *Rule) AddExpected(i int, t NodeType) {
	if r == nil {
		return
	}

	prev, ok := r.Expecteds[i]
	if !ok {
		prev = []NodeType{t}
	} else {
		pos, ok := slices.BinarySearch(prev, t)
		if !ok {
			prev = slices.Insert(prev, pos, t)
		}
	}

	r.Expecteds[i] = prev
}

func (r Rule) CheckExpected(at int, got_type NodeType) error {
	if len(r.Expecteds) == 0 {
		return nil
	}

	prev, ok := r.Expecteds[at]
	if !ok {
		return nil
	}

	_, ok = slices.BinarySearch(prev, got_type)
	if ok {
		return nil
	}

	values := gcstr.SliceOfStringer(prev)
	gcstr.QuoteStrings(values)

	got := strings.Quote(got_type.String())

	msg := gcstr.ExpectedValue("node type", gcstr.EitherOr(values), got)

	return grammar.NewBadParseTree(msg)
}

func (r Rule) ApplyField(children []*grammar.ParseTree[TokenType]) ([]*Node, error) {
	var sub_nodes []*Node

	for i, field := range r.Fields {
		if field.IsTerminal() {
			err := ast.CheckType(children, i, field)
			if err != nil {
				return sub_nodes, err
			}

			continue
		}

		if field.IsLhsRule() {
			rules, ok := GetRule(field)
			if !ok || len(rules) == 0 {
				return sub_nodes, fmt.Errorf("no rules for %q", field.String())
			} else if len(rules) > 1 {
				panic(fmt.Errorf("multiple rules for %q; case not yet implemented", field.String()))
			}

			sub_children, err := ast.LhsToAst(i, children, field, rules[0])
			if err != nil {
				return err
			}

			sub_nodes = append(sub_nodes, sub_children...)
		} else {
			node, err := ast_maker.Convert(children[i])
			if err != nil {
				return sub_nodes, err
			}

			err = r.CheckExpected(i, node.Type)
			if err != nil {
				return sub_nodes, err
			}

			sub_nodes = append(sub_nodes, node)
		}
	}

	return sub_nodes, nil
}
