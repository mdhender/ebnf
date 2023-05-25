// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

// Package tokens defines tokens returned from scanning.
package tokens

import "fmt"

// Position is line and column in the input
type Position struct {
	Line, Col int
}

// Token represents a token scanned from the input
type Token struct {
	Pos  Position
	Kind Kind
	Text []byte
}

func (t *Token) Column() int {
	if t == nil {
		return 0
	}
	return t.Pos.Col
}

func (t *Token) Line() int {
	if t == nil {
		return 0
	}
	return t.Pos.Line
}

func (t *Token) String() string {
	switch t.Kind {
	case UNKNOWN:
		return fmt.Sprintf("(%d %q)", t.Line(), string(t.Text))
	case NONTERMINAL:
		return fmt.Sprintf("(%d %s)", t.Line(), string(t.Text))
	case EQ:
		return fmt.Sprintf("(%d '=')", t.Line())
	case TERMINATOR:
		return fmt.Sprintf("(%d '.')", t.Line())
	case OR:
		return fmt.Sprintf("(%d '|')", t.Line())
	case TERMINAL:
		return fmt.Sprintf("(%d %s)", t.Line(), string(t.Text))
	case START_GROUP:
		return fmt.Sprintf("(%d '(')", t.Line())
	case END_GROUP:
		return fmt.Sprintf("(%d ')')", t.Line())
	case START_OPTION:
		return fmt.Sprintf("(%d '[')", t.Line())
	case END_OPTION:
		return fmt.Sprintf("(%d ']')", t.Line())
	case START_REPETITION:
		return fmt.Sprintf("(%d '{')", t.Line())
	case END_REPETITION:
		return fmt.Sprintf("(%d '}')", t.Line())
	case EOF:
		return fmt.Sprintf("(%d $)", t.Line())
	}
	return fmt.Sprintf("(%d %d %q)", t.Line(), t.Kind, string(t.Text))
}

type Kind int

func (k Kind) String() string {
	switch k {
	case UNKNOWN:
		return "UNKNOWN"
	case NONTERMINAL:
		return "NONTERMINAL"
	case EQ:
		return "EQ"
	case TERMINATOR:
		return "TERMINATOR"
	case OR:
		return "OR"
	case TERMINAL:
		return "TERMINAL"
	case START_GROUP:
		return "START_GROUP"
	case END_GROUP:
		return "END_GROUP"
	case START_OPTION:
		return "START_OPTION"
	case END_OPTION:
		return "END_OPTION"
	case START_REPETITION:
		return "START_REPETITION"
	case END_REPETITION:
		return "END_REPETITION"
	case EOF:
		return "EOF"
	}
	panic(fmt.Sprintf("assert(kind != %d)", k))
}

const (
	UNKNOWN Kind = iota
	NONTERMINAL
	EQ
	TERMINAL
	OR
	TERMINATOR
	START_GROUP
	END_GROUP
	START_OPTION
	END_OPTION
	START_REPETITION
	END_REPETITION
	EOF
)
