package main

import (
	"fmt"
	"jedil/pkg/compiler"
)

func main() {
	source := "return vec3(1, 2, 3) + vec3(4, 5, 6)"
	parser := compiler.NewParser(source)
	program, err := parser.Parse()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
	} else {
		fmt.Printf("Parse success! Got %d statements\n", len(program.Statements))
	}
}
