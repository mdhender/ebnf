// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package scanners_test

import (
	"github.com/mdhender/ebnf/scanners"
	"github.com/mdhender/ebnf/tokens"
	"testing"
)

func TestScanner(t *testing.T) {
	for _, tc := range []struct {
		id     int
		dump   bool // if true, log all tokens
		input  string
		expect []tokens.Kind
	}{
		{id: 1, input: "a6 = B_5 . | ( ) [ ] { } ; comments", expect: []tokens.Kind{
			tokens.NONTERMINAL, tokens.EQ, tokens.TERMINAL, tokens.TERMINATOR, tokens.OR,
			tokens.START_GROUP, tokens.END_GROUP,
			tokens.START_OPTION, tokens.END_OPTION,
			tokens.START_REPETITION, tokens.END_REPETITION,
			tokens.EOF,
		}},
		{id: 2, input: "; comment\n a", expect: []tokens.Kind{
			tokens.NONTERMINAL,
			tokens.EOF,
		}},
		{id: 3, input: "a.b|c(d[e{f}g]h)i;", expect: []tokens.Kind{
			tokens.NONTERMINAL, tokens.TERMINATOR,
			tokens.NONTERMINAL, tokens.OR,
			tokens.NONTERMINAL, tokens.START_GROUP,
			tokens.NONTERMINAL, tokens.START_OPTION,
			tokens.NONTERMINAL, tokens.START_REPETITION,
			tokens.NONTERMINAL, tokens.END_REPETITION,
			tokens.NONTERMINAL, tokens.END_OPTION,
			tokens.NONTERMINAL, tokens.END_GROUP,
			tokens.NONTERMINAL,
			tokens.EOF,
		}},
		{id: 4, input: "b@t\na\nJab+Ba;\nb", expect: []tokens.Kind{
			tokens.UNKNOWN,
			tokens.NONTERMINAL,
			tokens.UNKNOWN,
			tokens.NONTERMINAL,
			tokens.EOF,
		}},
	} {
		toks := scanners.Scan([]byte(tc.input))
		if tc.dump {
			for _, token := range toks {
				t.Logf("%d: %d:%d: %s\n", tc.id, token.Line(), token.Column(), token.String())
			}
		}
		expects := tc.expect
		pos := 1
		for len(toks) != 0 && len(expects) != 0 {
			expect, got := expects[0], toks[0].Kind
			expects, toks = expects[1:], toks[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
		for len(toks) != 0 {
			expect, got := tokens.EOF, toks[0].Kind
			toks = toks[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
		for len(expects) != 0 {
			expect, got := expects[0], tokens.EOF
			expects = expects[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
	}
}
