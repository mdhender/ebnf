// Copyright 2023 Michael D Henderson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package ebnf

import (
	"fmt"
	"unicode/utf8"
)

// A Grammar is a set of EBNF productions.
// The map is indexed by production name.
type Grammar map[string]*Production

// ----------------------------------------------------------------------------
// Grammar verification

func isNonTerminal(tok *Token) bool {
	return tok.Kind == NONTERMINAL
}

func isTerminal(tok *Token) bool {
	return tok.Kind == TERMINAL
}

type verifier struct {
	errors   errorList
	worklist []*Production
	reached  Grammar // set of productions reached from (and including) the root production
	grammar  Grammar
}

func (v *verifier) error(format string, args ...any) {
	v.errors = append(v.errors, fmt.Errorf(format, args...))
}

func (v *verifier) push(prod *Production) {
	name := prod.Name.String()
	if _, found := v.reached[name]; !found {
		v.worklist = append(v.worklist, prod)
		v.reached[name] = prod
	}
}

func (v *verifier) verifyChar(x *Token) rune {
	s := string(x.Text)
	if utf8.RuneCountInString(s) != 1 {
		v.error("%d: single char expected, found %q", x.Line, s)
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(s)
	return ch
}

func (v *verifier) verifyExpr(expr Expression, lexical bool) {
	switch x := expr.(type) {
	case nil:
		// empty expression
	case Alternative:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case Sequence:
		for _, e := range x {
			v.verifyExpr(e, lexical)
		}
	case *Name:
		// a production with this name must exist;
		// add it to the worklist if not yet processed
		if prod, found := v.grammar[x.String()]; found {
			v.push(prod)
		} else {
			v.error("%d: missing production %q", x.tok.Line, x.String())
		}
		// within a lexical production references to non-lexical productions are invalid
		if lexical && isTerminal(x.tok) {
			v.error("%d: reference to non-lexical production %q", x.tok.Line, x.String())
		}
	case *Literal:
		// nothing to do for now
	case *Group:
		v.verifyExpr(x.Body, lexical)
	case *Option:
		v.verifyExpr(x.Body, lexical)
	case *Repetition:
		v.verifyExpr(x.Body, lexical)
	case *Bad:
		v.error("%d: %v", x.tok.Line, x.err)
	default:
		panic(fmt.Sprintf("internal error: unexpected type %T", expr))
	}
}

func (v *verifier) verify(grammar Grammar, start string) {
	// find root production
	root, found := grammar[start]
	if !found {
		v.error("%d: no start production %q", 0, start)
		return
	}

	// initialize verifier
	v.worklist = v.worklist[0:0]
	v.reached = make(Grammar)
	v.grammar = grammar

	// work through the worklist
	v.push(root)
	for {
		n := len(v.worklist) - 1
		if n < 0 {
			break
		}
		prod := v.worklist[n]
		v.worklist = v.worklist[0:n]
		v.verifyExpr(prod.Expr, isNonTerminal(prod.Name.tok))
	}

	// check if all productions were reached
	if len(v.reached) < len(v.grammar) {
		for name, prod := range v.grammar {
			if _, found := v.reached[name]; !found {
				v.error("%d: %q is unreachable", prod.Pos(), name)
			}
		}
	}
}

// Verify checks that:
//   - all productions used are defined
//   - all productions defined are used when beginning at start
//   - lexical productions refer only to other lexical productions
//
// Position information is interpreted relative to the file set fset.
func Verify(grammar Grammar, start string) error {
	var v verifier
	v.verify(grammar, start)
	return v.errors.Err()
}
