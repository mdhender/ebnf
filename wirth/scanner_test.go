package wirth_test

import (
	"github.com/mdhender/ebnf/wirth"
	"testing"
)

func TestScanner(t *testing.T) {
	for _, tc := range []struct {
		id     int
		dump   bool // if true, log all tokens
		input  string
		expect []wirth.Kind
	}{
		{id: 1, input: "a6 = B_5 . | ( ) [ ] { } ; comments", expect: []wirth.Kind{
			wirth.NONTERMINAL, wirth.EQ, wirth.TERMINAL, wirth.TERMINATOR, wirth.OR,
			wirth.START_GROUP, wirth.END_GROUP,
			wirth.START_OPTION, wirth.END_OPTION,
			wirth.START_REPETITION, wirth.END_REPETITION,
			wirth.EOF,
		}},
		{id: 2, input: "; comment\n a", expect: []wirth.Kind{
			wirth.NONTERMINAL,
			wirth.EOF,
		}},
		{id: 3, input: "a.b|c(d[e{f}g]h)i;", expect: []wirth.Kind{
			wirth.NONTERMINAL, wirth.TERMINATOR,
			wirth.NONTERMINAL, wirth.OR,
			wirth.NONTERMINAL, wirth.START_GROUP,
			wirth.NONTERMINAL, wirth.START_OPTION,
			wirth.NONTERMINAL, wirth.START_REPETITION,
			wirth.NONTERMINAL, wirth.END_REPETITION,
			wirth.NONTERMINAL, wirth.END_OPTION,
			wirth.NONTERMINAL, wirth.END_GROUP,
			wirth.NONTERMINAL,
			wirth.EOF,
		}},
		{id: 4, input: "b@t\na\nJab+Ba;\nb", expect: []wirth.Kind{
			wirth.UNKNOWN,
			wirth.NONTERMINAL,
			wirth.UNKNOWN,
			wirth.NONTERMINAL,
			wirth.EOF,
		}},
	} {
		tokens := wirth.Scan([]byte(tc.input))
		if tc.dump {
			for _, token := range tokens {
				t.Logf("%d: %d:%d: %s\n", tc.id, token.Line(), token.Column(), token.String())
			}
		}
		expects := tc.expect
		pos := 1
		for len(tokens) != 0 && len(expects) != 0 {
			expect, got := expects[0], tokens[0].Kind
			expects, tokens = expects[1:], tokens[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
		for len(tokens) != 0 {
			expect, got := wirth.EOF, tokens[0].Kind
			tokens = tokens[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
		for len(expects) != 0 {
			expect, got := expects[0], wirth.EOF
			expects = expects[1:]
			if expect != got {
				t.Errorf("%d: %d: want %s, got %s\n", tc.id, pos, expect.String(), got.String())
			}
			pos++
		}
	}
}
