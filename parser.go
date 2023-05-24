// Copyright 2023 Michael D Henderson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package ebnf

import (
	"fmt"
	"github.com/mdhender/ebnf/scanners"
	"github.com/mdhender/ebnf/tokens"
)

// Parse parses a set of EBNF productions from the input.
// It returns a set of productions.
// Errors are reported for incorrect syntax and if a production
// is declared more than once.
func Parse(input []byte) (Grammar, []error) {
	toks := scanners.Scan(input)

	var p parser
	grammar := p.parse(toks)
	return grammar, p.errors
}

type parser struct {
	pos    int           // token position
	eof    *tokens.Token // end of tokens
	tok    *tokens.Token // one token look-ahead
	lit    string        // token literal
	tokens []*tokens.Token
	errors errorList
}

// parse parses a grammar
// --> grammar     ::= production { production } .
// --> production  ::= NONTERMINAL EQ [ expression ] TERMINATOR .
// --> expression  ::= sequence { OR sequence } .
// --> sequence    ::= term { term } .
// --> term        ::= NONTERMINAL | TERMINAL | group | option | repetition .
// --> group       ::= LPAREN   expression RPAREN   .
// --> option      ::= LBRACKET expression RBRACKET .
// --> repetition  ::= LBRACE   expression RBRACE   .
func (p *parser) parse(toks []*tokens.Token) (grammar Grammar) {
	p.tokens = toks
	if len(toks) > 0 {
		p.eof = toks[len(toks)-1]
	}

	// initializes pos, tok, lit
	p.next()

	grammar = make(Grammar)
	for p.tok != p.eof {
		prod := p.parseProduction()
		name := prod.Name.String()
		if def, found := grammar[name]; found {
			p.error("%d: %s: defined line %d", prod.Name.tok.Line(), def.Name.String(), def.Name.tok.Line())
			continue
		}
		grammar[name] = prod
	}

	return grammar
}

// parseProduction parses
// --> production  ::= NONTERMINAL EQ [ expression ] TERMINATOR .
func (p *parser) parseProduction() *Production {
	name := p.parseNonTerminal()
	p.expect(tokens.EQ)
	var expr Expression
	if p.tok.Kind != tokens.TERMINATOR {
		expr = p.parseExpression()
	}
	p.expect(tokens.TERMINATOR)
	return &Production{Name: name, Expr: expr}
}

// parseExpression parses
// --> expression  ::= sequence { OR sequence } .
func (p *parser) parseExpression() Expression {
	var list Alternative

	list = append(list, p.parseSequence())
	for p.tok.Kind == tokens.OR {
		p.next()
		list = append(list, p.parseSequence())
	}

	// no need for an Alternative node if list.Len() < 2
	if len(list) == 1 {
		return list[0]
	}

	return list
}

// parseSequence parses
// --> sequence    ::= term { term } .
func (p *parser) parseSequence() Expression {
	var list Sequence

	for x := p.parseTerm(); x != nil; x = p.parseTerm() {
		list = append(list, x)
	}

	// it is an error if the list is empty
	if len(list) == 0 {
		p.errorExpected(p.pos, "term", p.tok)
		return &Bad{
			tok: p.tok,
			err: fmt.Errorf("%d: term expected", p.tok.Line()),
		}
	}

	// no need for a sequence if list is just one term.
	if len(list) == 1 {
		return list[0]
	}

	return list
}

// parseTerm parses
// --> term        ::= NONTERMINAL | TERMINAL | group | option | repetition .
// --> group       ::= LPAREN   expression RPAREN   .
// --> option      ::= LBRACKET expression RBRACKET .
// --> repetition  ::= LBRACE   expression RBRACE   .
// Returns nil if no term was found.
func (p *parser) parseTerm() (x Expression) {
	tok := p.tok
	switch p.tok.Kind {
	case tokens.NONTERMINAL:
		x = p.parseNonTerminal()

	case tokens.TERMINAL:
		x = p.parseTerminal()

	case tokens.START_GROUP:
		p.next()
		x = &Group{tok: tok, Body: p.parseExpression()}
		p.expect(tokens.END_GROUP)

	case tokens.START_OPTION:
		p.next()
		x = &Option{tok: tok, Body: p.parseExpression()}
		p.expect(tokens.END_OPTION)

	case tokens.START_REPETITION:
		p.next()
		x = &Repetition{tok: tok, Body: p.parseExpression()}
		p.expect(tokens.END_REPETITION)
	}

	return x
}

// parseNonTerminal parses a NONTERMINAL.
func (p *parser) parseNonTerminal() *Name {
	tok := p.tok
	p.expect(tokens.NONTERMINAL)
	return &Name{tok: tok}
}

// parseTerminal parses a TERMINAL.
func (p *parser) parseTerminal() *Literal {
	tok := p.tok
	if p.tok.Kind == tokens.TERMINAL {
		p.next()
	} else { // didn't find terminal?
		p.expect(tokens.TERMINAL)
	}
	return &Literal{tok: tok}
}

// next returns the next token from the scanner.
// it never advances past the last token in the scanner.
func (p *parser) next() {
	if p.tok != p.eof {
		p.tok = p.tokens[p.pos]
		p.pos = p.pos + 1
	}
	p.lit = string(p.tok.Text)
}

// error appends an error to the parser's list of errors.
func (p *parser) error(format string, args ...any) {
	p.errors = append(p.errors, fmt.Errorf(format, args...))
}

// errorExpected generates an error and appends it to the parser's list of errors.
func (p *parser) errorExpected(pos int, want string, got *tokens.Token) {
	if got == nil {
		p.error("%d: expected %q", pos, want)
		return
	}
	// the error happened at the current position;
	// make the error message more specific
	if len(got.Text) == 0 {
		p.error("%d: expected %q, found %s", pos, want, got.Kind.String())
		return
	}
	p.error("%d: expected %q, found %s: %q", pos, want, got.Kind.String(), string(got.Text))
}

// expect compares the current input token against the expected kind of token.
// if the token doesn't match, it generates an error and appends the error to
// the parser's list of errors.
// always calls next() to advance the input.
func (p *parser) expect(k tokens.Kind) {
	if p.tok.Kind != k {
		p.errorExpected(p.pos, k.String(), p.tok)
	}
	p.next() // make progress in any case
}
