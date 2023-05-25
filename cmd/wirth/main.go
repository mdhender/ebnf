// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package main

import (
	"fmt"
	"github.com/mdhender/ebnf/tokens"
	"github.com/mdhender/ebnf/wirth"
	"log"
	"os"
	"sort"
)

func main() {
	input, err := os.ReadFile("lua.ebnf")
	if err != nil {
		log.Fatal(err)
	}
	grammar, err := wirth.Parse(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("grammar: start symbol %q\n", string(grammar.Start.Text))
	var identifiers []*tokens.Token
	for _, production := range grammar.Productions {
		identifiers = append(identifiers, production.Identifier)
	}
	sort.Slice(identifiers, func(i, j int) bool {
		return identifiers[i].Line() < identifiers[j].Line()
	})
	for _, identifier := range identifiers {
		name := string(identifier.Text)
		production := grammar.Productions[name]
		fmt.Printf("grammar: %6d %s\n", production.Identifier.Line(), name)
	}
	sort.Slice(identifiers, func(i, j int) bool {
		return string(identifiers[i].Text) < string(identifiers[j].Text)
	})
	for _, identifier := range identifiers {
		name := string(identifier.Text)
		production := grammar.Productions[name]
		fmt.Printf("grammar: %-30s %6d\n", name, production.Identifier.Line())
	}
}
