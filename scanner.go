// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package ebnf

//import (
//	"bytes"
//	"github.com/mdhender/ebnf/tokens"
//	"unicode"
//	"unicode/utf8"
//)
//
//func Scan(input []byte) []*tokens.Token {
//	if len(input) == 0 {
//		return []*tokens.Token{&tokens.Token{Kind: tokens.EOF}}
//	}
//	var toks []*tokens.Token
//	var token *tokens.Token
//	pos := tokens.Position{Line: 1, Col: 1}
//	for len(input) != 0 {
//		token, input = next(input)
//		if token == nil {
//			continue
//		}
//		token.Pos = pos
//		if token.Kind == tokens.UNKNOWN && bytes.Equal(token.Text, []byte{'\n'}) {
//			pos.Line++
//			continue
//		}
//		toks = append(toks, token)
//	}
//	return append(toks, &tokens.Token{Pos: pos, Kind: tokens.EOF})
//}
//
//var (
//	delims = []byte("()[]{}.=|; \f\n\r\t\v")
//)
//
//// next returns the next token from the input, skipping spaces, comments, and invalid runes.
//// returns nil only if the input is empty.
//func next(buffer []byte) (*tokens.Token, []byte) {
//	// skip spaces, invalid runes, and comments
//	for len(buffer) != 0 && buffer[0] != '\n' {
//		if buffer[0] == ';' {
//			if eol := bytes.IndexByte(buffer, '\n'); eol == -1 {
//				buffer = nil
//			} else {
//				buffer = buffer[eol:]
//			}
//		} else if r, w := utf8.DecodeRune(buffer); r == utf8.RuneError || unicode.IsSpace(r) {
//			buffer = buffer[w:]
//		} else {
//			break
//		}
//	}
//
//	if len(buffer) == 0 {
//		return nil, nil
//	} else if buffer[0] == '\n' {
//		return &tokens.Token{Kind: tokens.UNKNOWN, Text: buffer[:1]}, buffer[1:]
//	} else if buffer[0] == '=' {
//		return &tokens.Token{Kind: tokens.EQ}, buffer[1:]
//	} else if buffer[0] == ')' {
//		return &tokens.Token{Kind: tokens.END_GROUP}, buffer[1:]
//	} else if buffer[0] == ']' {
//		return &tokens.Token{Kind: tokens.END_OPTION}, buffer[1:]
//	} else if buffer[0] == '}' {
//		return &tokens.Token{Kind: tokens.END_REPETITION}, buffer[1:]
//	} else if buffer[0] == '|' {
//		return &tokens.Token{Kind: tokens.OR}, buffer[1:]
//	} else if buffer[0] == '(' {
//		return &tokens.Token{Kind: tokens.START_GROUP}, buffer[1:]
//	} else if buffer[0] == '[' {
//		return &tokens.Token{Kind: tokens.START_OPTION}, buffer[1:]
//	} else if buffer[0] == '{' {
//		return &tokens.Token{Kind: tokens.START_REPETITION}, buffer[1:]
//	} else if buffer[0] == '.' {
//		return &tokens.Token{Kind: tokens.TERMINATOR}, buffer[1:]
//	}
//
//	// token continues until a delimiter is reached.
//	// delimiters are spaces, invalid runes, comments, or any symbol like '(' or '|'.
//	start, length := buffer, 0
//	r, w := utf8.DecodeRune(buffer)
//	buffer, length = buffer[w:], length+w
//
//	var kind tokens.Kind
//	if unicode.IsLower(r) {
//		kind = tokens.NONTERMINAL
//	} else if unicode.IsUpper(r) {
//		kind = tokens.TERMINAL
//	} else {
//		kind = tokens.UNKNOWN
//	}
//
//	for len(buffer) != 0 && bytes.IndexByte(delims, buffer[0]) == -1 {
//		r, w = utf8.DecodeRune(buffer)
//		if r == utf8.RuneError || unicode.IsSpace(r) {
//			break
//		} else if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
//			kind = tokens.UNKNOWN
//		}
//		buffer, length = buffer[w:], length+w
//	}
//
//	return &tokens.Token{Kind: kind, Text: start[:length]}, buffer
//}
