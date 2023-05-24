// Copyright 2023 Michael D Henderson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ebnf

import (
	"github.com/mdhender/ebnf/scanners"
	"github.com/mdhender/ebnf/tokens"
	"testing"
)

var goodGrammars = []string{
	`program = .`,
	`program = Foo . ; end`,
	`program = Foo .`,
	`program = A | B C .`,
	`program = AtoZ .`,
	`program = song .
	 song = { note } .
	 note = Do | (Re Mi | Fa | So La) | ti .
	 ti = Ti .`,
	`program=song.song={note}.note=Do|(Re Mi|Fa|So La)|ti.ti=Ti.`,
}

var badParse = []string{
	`program = | .`,
	`program = | b .`,
	`program = a B ( .`,
	`program = A B ] .`,
	`program = B } .`,
	`program = = .`,
	`program = () .`,
	`program = [] .`,
	`program = {} .`,
	`program = song .
	 song Do | Ti .`,
	`program = song .
	 song = { note } .
	 note = Do | Ti .
	 note = Fa | La .`,
	`program = b59$ && foo .`,
}

var badVerify = []string{
	`program = a B .`,
	`start = a B .`,
	`program = A .
	 a = A .`,
}

func checkGood(t *testing.T, src string) {
	input := []byte(src)
	grammar, err := Parse(input)
	if err != nil {
		t.Errorf("Parse(%q) failed: %v", src, err)
		return
	}
	if err = Verify(grammar, "program"); err != nil {
		t.Errorf("Verify(%q) failed: %v", src, err)
	}
}

func checkBadParse(t *testing.T, src string) {
	input := []byte(src)
	if _, err := Parse(input); err == nil {
		t.Errorf("Parse(%q) should have failed", src)
	}
}

func checkBadVerify(t *testing.T, src string) {
	input := []byte(src)
	if grammar, err := Parse(input); err != nil {
		t.Errorf("Parse(%q) failed: %v", src, err)
	} else if err = Verify(grammar, "program"); err == nil {
		t.Errorf("Verify(%q) should have failed", src)
	}
}

func TestGrammars(t *testing.T) {
	debug := len(goodGrammars) == 0
	t.Logf("running goodGrammars\n")
	for _, src := range goodGrammars {
		if debug {
			for _, token := range scanners.Scan([]byte(src)) {
				t.Logf("%s\n", token)
				if token.Kind == tokens.EOF {
					break
				}
			}
		}
		checkGood(t, src)
	}

	t.Logf("running badParse\n")
	for _, src := range badParse {
		if debug {
			for _, token := range scanners.Scan([]byte(src)) {
				t.Logf("%s\n", token)
				if token.Kind == tokens.EOF {
					break
				}
			}
		}
		checkBadParse(t, src)
	}

	t.Logf("running badVerify\n")
	for _, src := range badVerify {
		if debug {
			for _, token := range scanners.Scan([]byte(src)) {
				t.Logf("%s\n", token)
				if token.Kind == tokens.EOF {
					break
				}
			}
		}
		checkBadVerify(t, src)
	}
}
