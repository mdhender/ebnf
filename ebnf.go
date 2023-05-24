// Copyright 2023 Michael D Henderson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

// Package ebnf is a library for EBNF grammars.
// The input is text ([]byte) satisfying the following grammar (represented itself in EBNF):
//
//	grammar     = production { production } .
//	production  = NONTERMINAL EQ [ expression ] TERMINATOR .
//	expression  = sequence { OR sequence } .
//	sequence    = term { term } .
//	term        = NONTERMINAL | TERMINAL | group | option | repetition .
//	group       = START_GROUP      expression END_GROUP      .
//	option      = START_OPTION     expression END_OPTION     .
//	repetition  = START_REPETITION expression END_REPETITION .
//
// A NONTERMINAL denotes a non-terminal production.
// A TERMINAL denotes a token returned from the scanner
//
//		NONTERMINAL      = LOWERLETTER { LETTER | DIGIT | UNDERSCORE }
//		TERMINAL         = UPPERLETTER { LETTER | DIGIT | UNDERSCORE }
//		EQ               = "="
//		OR               = "|"
//		START_GROUP      = "("
//		END_GROUP        = ")"
//		START_OPTION     = "["
//		END_OPTION       = "]"
//		START_REPETITION = "{"
//		END_REPETITION   = "}"
//	 COMMENT          = "//" ... EOL
//		TERMINATOR       = "."
//
// The scanner treats spaces, invalid runes, and comments as delimiters
// that separate tokens.
package ebnf
