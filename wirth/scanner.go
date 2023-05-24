// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package wirth

import (
	"bytes"
	"github.com/mdhender/ebnf/tokens"
	"unicode"
	"unicode/utf8"
)

// Scan returns a slice containing all the tokens in the input.
// It always adds an end of input token to that slice.
func Scan(input []byte) []*tokens.Token {
	pos := tokens.Position{Line: 1, Col: 1}
	s := &scanner{
		line:   pos.Line,
		col:    pos.Col,
		buffer: input,
		// delimiters are spaces, comments, any single character terminal, or invalid runes.
		delims: []byte(" \f\n\n\t\v;()[]{}.=|"),
	}
	var toks []*tokens.Token
	for token := s.next(); token != nil; token = s.next() {
		toks = append(toks, token)
		pos = token.Pos
	}
	return append(toks, &tokens.Token{Pos: pos, Kind: tokens.EOF})
}

type scanner struct {
	line, col int
	buffer    []byte
	delims    []byte
}

func (s *scanner) getch() rune {
	if s.iseof() {
		return utf8.RuneError
	}
	r, w := utf8.DecodeRune(s.buffer)
	s.buffer = s.buffer[w:]
	if r == '\n' {
		s.line, s.col = s.line+1, 0
	}
	s.col++
	return r
}

func (s *scanner) iseof() bool {
	return len(s.buffer) == 0
}

// next returns the next token from the input, skipping spaces, comments, and invalid runes.
// returns nil only if the input is empty.
func (s *scanner) next() *tokens.Token {
	// skip spaces, invalid runes, and comments
	for !s.iseof() {
		r := s.peekch()
		if r == ';' {
			if eol := bytes.IndexByte(s.buffer, '\n'); eol == -1 {
				s.buffer = nil
			} else {
				s.buffer = s.buffer[eol:]
			}
		} else if r == utf8.RuneError || unicode.IsSpace(r) {
			s.getch()
		} else {
			break
		}
	}

	if s.iseof() {
		return nil
	}

	tok := &tokens.Token{Pos: tokens.Position{Line: s.line, Col: s.col}}
	start := s.buffer
	r := s.getch()

	switch r {
	case '=':
		tok.Kind = tokens.EQ
	case ')':
		tok.Kind = tokens.END_GROUP
	case ']':
		tok.Kind = tokens.END_OPTION
	case '}':
		tok.Kind = tokens.END_REPETITION
	case '|':
		tok.Kind = tokens.OR
	case '(':
		tok.Kind = tokens.START_GROUP
	case '[':
		tok.Kind = tokens.START_OPTION
	case '{':
		tok.Kind = tokens.START_REPETITION
	case '.':
		tok.Kind = tokens.TERMINATOR
	default:
		if unicode.IsLower(r) {
			tok.Kind = tokens.NONTERMINAL
		} else if unicode.IsUpper(r) {
			tok.Kind = tokens.TERMINAL
		} else {
			tok.Kind = tokens.UNKNOWN
		}

		// token continues until a delimiter is reached.
		for !s.iseof() && bytes.IndexByte(s.delims, s.buffer[0]) == -1 {
			if r = s.peekch(); r == utf8.RuneError || unicode.IsSpace(r) {
				break
			} else if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
				tok.Kind = tokens.UNKNOWN
			}
			s.getch()
		}
		tok.Text = start[:len(start)-len(s.buffer)]
	}

	return tok
}

func (s *scanner) peekch() rune {
	if s.iseof() {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRune(s.buffer)
	return r
}
