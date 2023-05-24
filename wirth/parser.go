// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import "github.com/mdhender/ebnf/tokens"

func Parse(input []byte) (*Syntax, error) {
	p := &parser{tokens: Scan(input)}
	p.eof = p.tokens[len(p.tokens)-1]

	syntax := p.parse()
	return syntax, nil
}

type parser struct {
	tokens []*tokens.Token // all the tokens in the input
	eof    *tokens.Token   // last token in the input
	peek   *tokens.Token   // one token look-ahead
}

// parser parses a grammar file.
func (p *parser) parse() (syntax *Syntax) {
	syntax = &Syntax{Productions: make(map[string]*Production)}

	return syntax
}
