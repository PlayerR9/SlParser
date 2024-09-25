package kdd

import (
	"fmt"

	gerr "github.com/PlayerR9/go-errors/error"
	"github.com/dustin/go-humanize"
)

//go:generate stringer -type=ErrorCode

type ErrorCode int

const (
	// InvalidSyntax occurs when the AST is invalid or not
	// as it should be.
	InvalidSyntax ErrorCode = iota
)

// Int implements the error.ErrorCoder interface.
func (e ErrorCode) Int() int {
	return int(e)
}

/* // NewErrSyntax returns a new error.Err error representing a
// syntax error.
//
// Parameters:
//   - msg: The reason why the syntax is wrong.
//
// Returns:
//   - *error.Err: A pointer to the newly created error. Never returns nil.
//
// If the message is not specified, the string "the AST is invalid is used instead".
func NewErrSyntax(msg string) *gerr.Err {
	if msg == "" {
		msg = "the AST is invalid"
	}

	err := gerr.New(InvalidSyntax, msg)

	return err
} */

// check_ast_rec is a recursive function that checks if the given ast node is valid according to the following rules:
//
//   - All children must be valid.
//   - Source node must not be flagged as terminal and must have at least one children of type RuleNode.
//   - Rule node must not be flagged as terminal and must have at least two children of type RhsNode.
//   - Rhs node must not be empty and must have no children.
//
// Parameters:
//   - node: The node to check.
//   - depth: The search depth. -1 means no limit.
//
// Returns:
//   - *error.Err: A pointer to the newly created error.
func check_ast_rec(node *Node, depth int) *gerr.Err {
	if node == nil {
		err := gerr.New(InvalidSyntax, "node must not be nil")

		return err
	}

	if depth == 0 {
		return nil
	}

	if depth > 0 {
		depth--
	}

	children := node.GetChildren()

	// 0. All children must be valid.
	for _, child := range children {
		err := check_ast_rec(child, depth)
		if err == nil {
			continue
		}

		msg := fmt.Sprintf("%s child is invalid", humanize.Ordinal(len(children)))

		outer_err := gerr.New(InvalidSyntax, msg)
		outer_err.AddFrame("", node.Type.String())
		outer_err.SetInner(err)

		return outer_err
	}

	switch node.Type {
	case SourceNode:
		// 1. All children must be rule nodes.
		// 2. At least one children is expected.
		// 3. Must not flagged as terminal.

		if node.IsTerminal {
			err := gerr.New(InvalidSyntax, "source node must not be flagged as terminal")

			return err
		}

		if len(children) == 0 {
			err := gerr.New(InvalidSyntax, "at least one rule is expected")

			return err
		}

		for i, child := range children {
			if child.Type == RuleNode {
				continue
			}

			msg := fmt.Sprintf("expected %s child to be of type %q, got %q instead",
				humanize.Ordinal(i+1), RuleNode.String(), child.Type.String(),
			)

			err := gerr.New(InvalidSyntax, msg)

			return err
		}
	case RuleNode:
		// 1. All children must be rhs nodes.
		// 2. At least two children are expected.
		// 3. Must not flagged as terminal.

		if node.IsTerminal {
			err := gerr.New(InvalidSyntax, "rule node must not be flagged as terminal")

			return err
		}

		if len(children) == 0 {
			err := gerr.New(InvalidSyntax, "missing LHS node")
			return err
		} else if len(children) == 1 {
			err := gerr.New(InvalidSyntax, "missing RHS node")
			return err
		}

		for i, child := range children {
			if child.Type == RhsNode {
				continue
			}

			msg := fmt.Sprintf("expected %s child to be of type %q, got %q instead",
				humanize.Ordinal(i+1), RhsNode.String(), child.Type.String(),
			)

			err := gerr.New(InvalidSyntax, msg)
			return err
		}
	case RhsNode:
		// 1. No children are expected.
		// 2. Data must not be empty.

		if node.Data == "" {
			err := gerr.New(InvalidSyntax, "missing identifier")
			return err
		}

		if len(children) != 0 {
			err := gerr.New(InvalidSyntax, fmt.Sprintf("expected no children, got %d instead", len(children)))

			return err
		}
	default:
		err := gerr.New(InvalidSyntax, fmt.Sprintf("type %q is not supported", node.Type.String()))
		return err
	}

	return nil
}

// CheckASTWithLimit checks if the given ast node is valid, up to a given
// limit depth. If limit is negative, it will check all the way down to the
// leaves. On the other hand, if limit is 0, it will only check if the node
// is nil or not.
//
// Parameters:
//   - node: The node to check.
//   - limit: The maximum depth to check. If negative, it will check all the
//     way down to the leaves.
//
// Returns:
//   - error: If the node is not valid, an error describing the problem will
//     be returned. Otherwise, nil is returned.
//
// What Counts as a Valid Node?
//
// Overall Rules:
//  1. All children must be valid.
//  2. A node cannot be nil.
//
// SourceNode:
//  1. Must not be flagged as terminal.
//  2. Must have at least one children.
//  3. All children must be of type RuleNode.
//
// RuleNode:
//  1. Must not be flagged as terminal.
//  2. Must have at least two children. (The first is the LHS while the rest are the RHSs).
//  3. All children must be of type RhsNode.
//
// RhsNode:
//  1. No children are expected.
//  2. Data must not be empty.
func CheckASTWithLimit(node *Node, limit int) error {
	if limit < 0 {
		limit = -1
	}

	err := check_ast_rec(node, limit)
	if err != nil {
		return err
	}

	return nil
}

// CheckNode checks if the given ast node is valid. (See CheckASTWithLimit for
// more details).
//
// Parameters:
//   - node: The node to check.
//
// Returns:
//   - error: An error if the node is invalid. Otherwise, nil.
func CheckNode(node *Node) error {
	err := check_ast_rec(node, 1)
	if err != nil {
		return err
	}

	return nil
}

// CheckAST checks the given AST is valid or not in a recursive way in a DFS manner.
// (See CheckASTWithLimit for more details).
//
// Parameters:
//   - root: The root of the AST to check.
//
// Returns:
//   - error: An error describing why the AST is invalid, or nil if the AST is valid.
func CheckAST(root *Node) error {
	err := check_ast_rec(root, -1)
	if err != nil {
		return err
	}

	return nil
}
