// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import (
	"fmt"
	"github.com/mdhender/ebnf/scanners"
	"github.com/mdhender/ebnf/tokens"
)

func Parse(input []byte) (*Syntax, error) {
	p := &parser{tokens: scanners.Scan(input)}
	p.eof = p.tokens[len(p.tokens)-1]

	syntax := p.parse()
	return syntax, nil
}

type parser struct {
	tokens  []*tokens.Token // all the tokens in the input
	current int             // index of the current token
	eof     *tokens.Token   // last token in the input
	errors  []error         // all parsing errors
}

// parser parses a grammar file.
func (p *parser) parse() (syntax *Syntax) {
	syntax = &Syntax{Productions: make(map[string]*Production)}

	return syntax
}

// syntax recognizes
// --> { production }
func (p *parser) ntSyntax() (*Node, error) {
	for n, err := p.ntProduction(); n != nil; n, err = p.ntProduction() {
		if err != nil {
			p.addError(err)
		}
		// ... //
	}
	return nil, nil
}

// production recognizes
// --> NONTERMINAL EQ expression TERMINATOR
func (p *parser) ntProduction() (*Production, error) {
	identifier, ok := p.accept(tokens.NONTERMINAL)
	if !ok {
		return nil, nil
	}
	production := &Production{
		Identifier: identifier,
	}
	_, err := p.expect(tokens.EQ)
	if err != nil {
		p.addError(err)
	}
	production.Expression, err = p.ntExpression()
	if err != nil {
		p.addError(err)
	}
	if _, err = p.expect(tokens.TERMINATOR); err != nil {
		p.addError(err)
	}
	return production, nil
}

// expression recognizes
// --> term { OR term }
func (p *parser) ntExpression() (*Expression, error) {
	term, err := p.ntTerm()
	if err != nil {
		p.addError(err)
		return nil, fmt.Errorf("expected term")
	}
	expression := &Expression{
		Terms: []*Term{term},
	}
	for _, ok := p.accept(tokens.OR); ok; _, ok = p.accept(tokens.OR) {
		term, err = p.ntTerm()
		if err != nil {
			p.addError(err)
			return nil, fmt.Errorf("expected term")
		}
		expression.Terms = append(expression.Terms)
	}
	return expression, nil
}

// term recognizes
// --> factor { factor }
func (p *parser) ntTerm() (*Term, error) {
	factor, err := p.ntFactor()
	if err != nil {
		p.addError(err)
	}
	term := &Term{
		Factors: []*Factor{
			factor,
		},
	}
	for factor, err = p.ntFactor(); factor != nil && err != nil; factor, err = p.ntFactor() {
		if err != nil {
			p.addError(err)
		}
		if factor != nil {
			term.Factors = append(term.Factors, factor)
		}
	}
	return term, err
}

// factor recognizes
// * --> NONTERMINAL
// *   | TERMINAL
// *   | START_GROUP expression END_GROUP
// *   | START_OPTION expression END_OPTION
// *   | START_REPETITION expression END_REPETITION
func (p *parser) ntFactor() (*Factor, error) {
	if token, ok := p.accept(tokens.NONTERMINAL); ok {
		return &Factor{NonTerminal: token}, nil
	}
	if token, ok := p.accept(tokens.TERMINAL); ok {
		return &Factor{Terminal: token}, nil
	}
	if start, ok := p.accept(tokens.START_GROUP); ok {
		expression, err := p.ntExpression()
		if err != nil {
			p.addError(err)
		}
		end, err := p.expect(tokens.END_GROUP)
		return &Factor{
			Group: &Group{
				Start:      start,
				Expression: expression,
				End:        end},
		}, err
	}
	if start, ok := p.accept(tokens.START_OPTION); ok {
		expression, err := p.ntExpression()
		if err != nil {
			p.addError(err)
		}
		end, err := p.expect(tokens.END_OPTION)
		return &Factor{
			Option: &Option{
				Start:      start,
				Expression: expression,
				End:        end},
		}, err
	}
	if start, ok := p.accept(tokens.START_REPETITION); ok {
		expression, err := p.ntExpression()
		if err != nil {
			p.addError(err)
		}
		end, err := p.expect(tokens.END_REPETITION)
		return &Factor{
			Repetition: &Repetition{
				Start:      start,
				Expression: expression,
				End:        end},
		}, err
	}
	return nil, nil
}

func (p *parser) addError(err error) {
	p.errors = append(p.errors, err)
}

// accept returns the next terminal in the input if it matches,
// otherwise, it returns the token and false.
func (p *parser) accept(k tokens.Kind) (*tokens.Token, bool) {
	if token := p.peek(); token.Kind != k {
		return token, false
	}
	return p.next(), true
}

// expect reads the next token from the input.
// if the kind matches the expected kind, the token is returned.
// otherwise, the both the token and an error are returned.
func (p *parser) expect(k tokens.Kind) (token *tokens.Token, err error) {
	token = p.next()
	if token.Kind != k {
		if len(token.Text) == 0 {
			err = fmt.Errorf("%d:%d: expect %s, got %s", token.Line(), token.Column(), k.String(), token.String())
		} else {
			err = fmt.Errorf("%d:%d: expect %s, got %s: %q", token.Line(), token.Column(), k.String(), token.String(), string(token.Text))
		}
	}
	return token, err
}

func (p *parser) next() *tokens.Token {
	if p.current+1 >= len(p.tokens) {
		panic("assert(next != nil)")
	}
	p.current++
	return p.tokens[p.current]
}

func (p *parser) peek() *tokens.Token {
	if p.current+1 >= len(p.tokens) {
		panic("assert(peek != nil)")
	}
	return p.tokens[p.current+1]
}
