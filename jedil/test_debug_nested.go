package main

import (
	"fmt"
	"jedil/pkg/compiler"
	"jedil/pkg/vm"
)

func main() {
	source := `fn add(a, b) {
    return a + b
}
fn sub(a, b) {
    return a - b
}
fn calc() {
    let x = add(10, 5)
    let y = sub(20, 8)
    return x + y
}
return calc()`

	fmt.Println("Source:")
	fmt.Println(source)
	fmt.Println()

	// Compile
	instructions, err := compiler.CompileSource(source)
	if err != nil {
		fmt.Printf("Compile Error: %v\n", err)
		return
	}

	fmt.Println("Bytecode:")
	for i, inst := range instructions {
		fmt.Printf("  [%2d] %s %.0f\n", i, inst.Op, inst.Args)
	}
	fmt.Println()

	// Run
	machine := vm.New(instructions)
	err = machine.Run()
	if err != nil {
		fmt.Printf("Runtime Error: %v\n", err)
		return
	}

	// Get result
	result, err := machine.GetResult()
	if err != nil {
		fmt.Printf("Result Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result.String())
}
