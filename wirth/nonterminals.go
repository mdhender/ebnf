// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import "github.com/mdhender/ebnf/tokens"

type Syntax struct {
	Start       *tokens.Token
	Productions map[string]*Production
}

type Production struct {
	Name       *tokens.Token
	Expression *Expression
}

type Expression struct {
	Terms []*Term
}

type Term struct {
	Factors []*Factor
}

type Factor struct {
	NonTerminal *tokens.Token
	Terminal    *tokens.Token
	Group       *Expression
	Option      *Expression
	Repetition  *Expression
}
