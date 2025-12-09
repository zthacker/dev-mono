package main

import (
	"fmt"
	"jedil/pkg/compiler"
)

func main() {
	source := "return 2.0 + 3.0"
	fmt.Printf("Source: %s\n", source)
	
	instructions, err := compiler.CompileSource(source)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success! %d instructions\n", len(instructions))
		for i, inst := range instructions {
			fmt.Printf("  [%d] %s %.2f\n", i, inst.Op, inst.Args)
		}
	}
}
