// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import "github.com/mdhender/ebnf/tokens"

type Node struct {
	// Identifier is the name of the node.
	// If the node is a terminal, it is the name of the terminal.
	// Otherwise, it is the name of the production.
	Identifier string
	// type of node, terminal or non-terminal.
	Type Type
	// position of the node in the input.
	Pos tokens.Position
	// Successor and Alternative are set only for non-terminal nodes.
	Successor   *Node
	Alternative *Node
	// Errors contains any errors parsing or verifying the node.
	Errors []error
}

func (n *Node) IsNonTerminal() bool {
	return n.Type == NONTERMINALNODE
}

func (n *Node) IsTerminal() bool {
	return n.Type == TERMINALNODE
}

type Type int

const (
	BADNODE Type = iota
	NONTERMINALNODE
	TERMINALNODE
)
