// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import (
	"fmt"
	"github.com/mdhender/ebnf/scanners"
	"github.com/mdhender/ebnf/tokens"
	"log"
)

/*
	╔═══════════════════════════════════════════════════════════╗
	║ syntax     = { production } .                             ║
	║ production = NONTERMINAL EQ [expression] TERMINATOR .     ║
	║ expression = term { OR term } .                           ║
	║ term       = factor { factor } .                          ║
	║ factor     = NONTERMINAL                                  ║
	║            | TERMINAL                                     ║
	║            | START_GROUP      expression END_GROUP        ║
	║            | START_OPTION     expression END_OPTION       ║
	║            | START_REPETITION expression END_REPETITION . ║
	╚═══════════════════════════════════════════════════════════╝

	start             = syntax
	first(syntax)     = first(production), Ɛ
	first(production) = NONTERMINAL
	first(expression) = first(term)
	first(term)       = first(factor)
	first(factor)     = NONTERMINAL, TERMINAL, START_GROUP, START_OPTION, START_REPETITION
*/

func Parse(input []byte) (*Grammar, error) {
	p := &parser{tokens: scanners.Scan(input)}
	p.current = p.tokens[0]
	p.eof = p.tokens[len(p.tokens)-1]
	if p.eof.Kind != tokens.EOF {
		panic("parse: missing EOF token")
	}
	grammar := p.parse()
	if p.current != p.eof {
		// didn't parse all tokens???
		for p.current != p.eof {
			log.Printf("parse: %s\n", p.current.String())
			p.next()
		}
		return grammar, fmt.Errorf("parse: extra tokens")
	}
	return grammar, nil
}

type parser struct {
	current *tokens.Token   // current token in the input
	eof     *tokens.Token   // last token in the input
	errors  []error         // all parsing errors
	tokens  []*tokens.Token // all the tokens in the input
}

// parser parses a grammar file.
func (p *parser) parse() *Grammar {
	grammar := &Grammar{Productions: make(map[string]*Production)}
	syntax := p.ntSyntax()
	for _, production := range syntax.Productions {
		identifier := production.Identifier
		if identifier == nil || identifier.Kind != tokens.NONTERMINAL {
			continue
		}
		if grammar.Start == nil {
			grammar.Start = identifier
		}
		grammar.Productions[string(identifier.Text)] = production
	}
	if grammar.Start == nil {
		grammar.Start = &tokens.Token{Kind: tokens.NONTERMINAL, Text: []byte{'$'}}
	}
	return grammar
}

// syntax recognizes
// --> { production }
func (p *parser) ntSyntax() *Syntax {
	syntax := &Syntax{}
	for p.current != p.eof {
		if p.current.Kind != tokens.NONTERMINAL {
			p.addError(fmt.Errorf("%d:%d: expected identifier, got %s\n", p.current.Line(), p.current.Column(), p.current.Kind))
			p.next()
			continue
		}
		syntax.Productions = append(syntax.Productions, p.ntProduction())
	}
	return syntax
}

// production recognizes
// --> NONTERMINAL EQ [expression] TERMINATOR
func (p *parser) ntProduction() *Production {
	production := &Production{}
	production.Identifier = p.expect(tokens.NONTERMINAL)
	_ = p.expect(tokens.EQ)
	if p.firstExpression(p.current.Kind) {
		production.Expression = p.ntExpression()
	}
	_ = p.expect(tokens.TERMINATOR)
	return production
}

// expression recognizes
// --> term { OR term }
func (p *parser) ntExpression() *Expression {
	expression := &Expression{}
	expression.Terms = append(expression.Terms, p.ntTerm())
	for p.current.Kind == tokens.OR {
		_ = p.expect(tokens.OR)
		expression.Terms = append(expression.Terms, p.ntTerm())
	}
	return expression
}

// term recognizes
// --> factor { factor }
func (p *parser) ntTerm() *Term {
	term := &Term{}
	term.Factors = append(term.Factors, p.ntFactor())
	for p.firstFactor(p.current.Kind) {
		term.Factors = append(term.Factors, p.ntFactor())
	}
	return term
}

/*
factor recognizes

	--> NONTERMINAL
	  | TERMINAL
	  | START_GROUP expression END_GROUP
	  | START_OPTION expression END_OPTION
	  | START_REPETITION expression END_REPETITION
*/
func (p *parser) ntFactor() *Factor {
	factor := &Factor{}
	if p.current.Kind == tokens.NONTERMINAL {
		factor.NonTerminal = p.expect(tokens.NONTERMINAL)
	} else if p.current.Kind == tokens.TERMINAL {
		factor.Terminal = p.expect(tokens.TERMINATOR)
	} else if p.firstGroup(p.current.Kind) {
		factor.Group = p.ntGroup()
	} else if p.firstOption(p.current.Kind) {
		factor.Option = p.ntOption()
	} else if p.firstRepetition(p.current.Kind) {
		factor.Repetition = p.ntRepetition()
	} else { // should never happen
		p.addError(fmt.Errorf("%d:%d: expected identifier, '(', '[', or '[', got %s", p.current.Line(), p.current.Column(), p.current.String()))
		p.next()
	}
	return factor
}

// group recognizes
// --> START_GROUP expression END_GROUP
func (p *parser) ntGroup() *Group {
	group := &Group{}
	group.Start = p.expect(tokens.START_GROUP)
	group.Expression = p.ntExpression()
	group.End = p.expect(tokens.END_GROUP)
	return group
}

// option recognizes
// --> START_OPTION expression END_OPTION
func (p *parser) ntOption() *Option {
	option := &Option{}
	option.Start = p.expect(tokens.START_OPTION)
	option.Expression = p.ntExpression()
	option.End = p.expect(tokens.END_OPTION)
	return option
}

// repetition recognizes
// --> START_REPETITION expression END_REPETITION
func (p *parser) ntRepetition() *Repetition {
	repetition := &Repetition{}
	repetition.Start = p.expect(tokens.START_REPETITION)
	repetition.Expression = p.ntExpression()
	repetition.End = p.expect(tokens.END_REPETITION)
	return repetition
}

func (p *parser) firstExpression(k tokens.Kind) bool {
	return p.firstTerm(k)
}

func (p *parser) firstTerm(k tokens.Kind) bool {
	return p.firstFactor(k)
}

func (p *parser) firstFactor(k tokens.Kind) bool {
	return k == tokens.NONTERMINAL || k == tokens.TERMINAL || k == tokens.START_GROUP || k == tokens.START_OPTION || k == tokens.START_REPETITION
}

func (p *parser) firstGroup(k tokens.Kind) bool {
	return k == tokens.START_GROUP
}

func (p *parser) firstOption(k tokens.Kind) bool {
	return k == tokens.START_OPTION
}

func (p *parser) firstRepetition(k tokens.Kind) bool {
	return k == tokens.START_REPETITION
}

func (p *parser) addError(err error) {
	p.errors = append(p.errors, err)
}

// expect reads the next token from the input.
// if the kind matches the expected kind, the token is returned.
// otherwise, the both the token and an error are returned.
func (p *parser) expect(k tokens.Kind) *tokens.Token {
	token := p.current
	p.next()
	if token.Kind != k {
		if len(token.Text) == 0 {
			p.addError(fmt.Errorf("%d:%d: expect %s, got %s", token.Line(), token.Column(), k.String(), token.String()))
		} else {
			p.addError(fmt.Errorf("%d:%d: expect %s, got %s: %q", token.Line(), token.Column(), k.String(), token.String(), string(token.Text)))
		}
	}
	return token
}

func (p *parser) next() {
	if p.current != p.eof {
		p.current, p.tokens = p.tokens[1], p.tokens[1:]
	}
}

func (p *parser) peek() *tokens.Token {
	return p.tokens[0]
}
