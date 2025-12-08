package main

import (
	"fmt"
	"jedil/pkg/bytecode"
	"jedil/pkg/vm"
)

func main() {

	fmt.Println("Hello, Jedil VM!")

	// simple add
	fmt.Println("2 + 3 = ", 2+3)
	code1 := []bytecode.Instruction{
		{Op: bytecode.OP_PUSH, Args: 2.0},
		{Op: bytecode.OP_PUSH, Args: 3.0},
		{Op: bytecode.OP_ADD},
		{Op: bytecode.OP_HALT},
	}

	runTest(code1, 5.0)
}

func runTest(code []bytecode.Instruction, expected float64) {
	v := vm.New(code)
	err := v.Run()
	if err != nil {
		fmt.Printf("VM error: %v\n", err)
		return
	}

	result, err := v.GetResult()
	if err != nil {
		fmt.Printf("Error getting result: %v\n", err)
		return
	}

	if result.AsFloat() == expected {
		fmt.Printf("Test passed: result = %.6f\n", result.AsFloat())
	} else {
		fmt.Printf("Test failed: expected = %.6f, got = %.6f\n", expected, result.AsFloat())
	}
}
