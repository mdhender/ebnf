package main

import (
	"fmt"
	"github.com/mdhender/ebnf"
	"log"
	"os"
)

func main() {
	src := "lua.ebnf"
	if input, err := os.ReadFile(src); err != nil {
		log.Fatal(err)
	} else if grammar, errors := ebnf.Parse(input); errors != nil {
		for _, err := range errors {
			fmt.Printf("Parse(%s) failed: %v\n", src, err)
		}
	} else if errors = ebnf.Verify(grammar, "chunk"); errors != nil {
		for _, err := range errors {
			fmt.Printf("Verify(%s) failed: %v\n", src, err)
		}
	}
}
