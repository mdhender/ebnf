// Copyright 2023 Michael D Henderson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package ebnf

// ----------------------------------------------------------------------------
// Internal representation

type (
	// A Production node represents an EBNF production.
	Production struct {
		Name *Name
		Expr Expression
	}

	// An Expression node represents a production expression.
	Expression interface {
		// Pos is the index of the first token in the construction
		Pos() int
	}

	// An Alternative node represents a non-empty list of alternative expressions.
	Alternative []Expression // x | y | z

	// A Sequence node represents a non-empty list of sequential expressions.
	Sequence []Expression // x y z

	// A Name node represents a production name.
	Name struct {
		tok *Token
	}

	// A Literal node represents a literal.
	Literal struct {
		tok *Token
	}

	//// A List node represents a range of characters.
	//Range struct {
	//	Begin, End *Literal // begin ... end
	//}

	// A Group node represents a grouped expression.
	Group struct {
		tok  *Token
		Body Expression // (body)
	}

	// An Option node represents an optional expression.
	Option struct {
		tok  *Token
		Body Expression // [body]
	}

	// A Repetition node represents a repeated expression.
	Repetition struct {
		tok  *Token
		Body Expression // {body}
	}

	// A Bad node stands for pieces of source code that lead to a parse error.
	Bad struct {
		tok *Token
		err error // parser error message
	}
)

func (x Alternative) Pos() int { return x[0].Pos() } // the parser always generates non-empty Alternative
func (x Sequence) Pos() int    { return x[0].Pos() } // the parser always generates non-empty Sequences
func (x *Name) Pos() int       { return x.tok.Pos() }
func (x *Literal) Pos() int    { return x.tok.Pos() }
func (x *Group) Pos() int      { return x.tok.Pos() }
func (x *Option) Pos() int     { return x.tok.Pos() }
func (x *Repetition) Pos() int { return x.tok.Pos() }
func (x *Production) Pos() int { return x.Name.Pos() }
func (x *Bad) Pos() int        { return x.Pos() }

func (x *Name) String() string { return string(x.tok.Text) }
