// Copyright 2023 Michael D Henderson.
// Use of this source code is governed by a BSD-style
// license that can be found in the COPYING file.

package main

import (
	"fmt"
	"github.com/mdhender/ebnf/scanners"
	"log"
	"os"
)

func main() {
	input, err := os.ReadFile("lua.ebnf")
	if err != nil {
		log.Fatal(err)
	}
	tokens := scanners.Scan(input)
	for id, token := range tokens {
		fmt.Printf("%4d: %4d: %3d: %s\n", id, token.Line(), token.Column(), token.String())
	}
}
