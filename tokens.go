// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package ebnf

import (
	"fmt"
)

// Token represents a token scanned from the input
type Token struct {
	Id   int
	Line int
	Kind Kind
	Text []byte
}

func (t *Token) Pos() int {
	return t.Id
}

func (t *Token) String() string {
	switch t.Kind {
	case UNKNOWN:
		return fmt.Sprintf("(%d %q)", t.Line, string(t.Text))
	case END_GROUP:
		return fmt.Sprintf("(%d ')')", t.Line)
	case END_OPTION:
		return fmt.Sprintf("(%d ']')", t.Line)
	case END_REPETITION:
		return fmt.Sprintf("(%d '}')", t.Line)
	case EOF:
		return fmt.Sprintf("(%d eof)", t.Line)
	case EOL:
		return fmt.Sprintf("(%d eol)", t.Line)
	case EQ:
		return fmt.Sprintf("(%d '=')", t.Line)
	case OR:
		return fmt.Sprintf("(%d '|')", t.Line)
	case NONTERMINAL:
		return fmt.Sprintf("(%d %s)", t.Line, string(t.Text))
	case START_GROUP:
		return fmt.Sprintf("(%d '(')", t.Line)
	case START_OPTION:
		return fmt.Sprintf("(%d '[')", t.Line)
	case START_REPETITION:
		return fmt.Sprintf("(%d '{')", t.Line)
	case TERMINAL:
		return fmt.Sprintf("(%d %s)", t.Line, string(t.Text))
	case TERMINATOR:
		return fmt.Sprintf("(%d '.')", t.Line)
	}
	return fmt.Sprintf("(%d %d %q)", t.Line, t.Kind, string(t.Text))
}

type Kind int

func (k Kind) String() string {
	switch k {
	case UNKNOWN:
		return "UNKNOWN"
	case END_GROUP:
		return "END_GROUP"
	case END_OPTION:
		return "END_OPTION"
	case END_REPETITION:
		return "END_REPETITION"
	case EOF:
		return "EOF"
	case EOL:
		return "EOL"
	case EQ:
		return "EQ"
	case OR:
		return "OR"
	case NONTERMINAL:
		return "NONTERMINAL"
	case START_GROUP:
		return "START_GROUP"
	case START_OPTION:
		return "START_OPTION"
	case START_REPETITION:
		return "START_REPETITION"
	case TERMINAL:
		return "TERMINAL"
	case TERMINATOR:
		return "TERMINATOR"
	}
	panic(fmt.Sprintf("assert(kind != %d)", k))
}

const (
	UNKNOWN Kind = iota
	END_GROUP
	END_OPTION
	END_REPETITION
	EOF
	EOL
	EQ
	OR
	NONTERMINAL
	START_GROUP
	START_OPTION
	START_REPETITION
	TERMINAL
	TERMINATOR
)
