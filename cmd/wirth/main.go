package main

import (
	"fmt"
	"github.com/mdhender/ebnf/wirth"
	"log"
	"os"
)

func main() {
	input, err := os.ReadFile("lua.ebnf")
	if err != nil {
		log.Fatal(err)
	}
	tokens := wirth.Scan(input)
	for id, token := range tokens {
		fmt.Printf("%4d: %4d: %3d: %s\n", id, token.Line(), token.Column(), token.String())
	}
}
