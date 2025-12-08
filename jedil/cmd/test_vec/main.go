package main

import (
	"fmt"
	"jedil/pkg/bytecode"
	"jedil/pkg/vm"
)

func main() {
	fmt.Println("=== JEDIL Vector Math Test ===\n")

	// Test 1: Vector Addition
	// v1(1, 2, 3) + v2(4, 5, 6) = v3(5, 7, 9)
	fmt.Println("Test 1: Vector Addition")
	code1 := []bytecode.Instruction{
		// Push Vector 1 (1, 2, 3)
		{Op: bytecode.OP_PUSH, Args: 1.0}, // x
		{Op: bytecode.OP_PUSH, Args: 2.0}, // y
		{Op: bytecode.OP_PUSH, Args: 3.0}, // z
		{Op: bytecode.OP_VEC3},            // Stack: [v1]

		// Push Vector 2 (4, 5, 6)
		{Op: bytecode.OP_PUSH, Args: 4.0},
		{Op: bytecode.OP_PUSH, Args: 5.0},
		{Op: bytecode.OP_PUSH, Args: 6.0},
		{Op: bytecode.OP_VEC3}, // Stack: [v1, v2]

		{Op: bytecode.OP_VADD}, // Stack: [v1 + v2]
		{Op: bytecode.OP_HALT},
	}
	runVectorTest(code1, "v3d(5.000000, 7.000000, 9.000000)")

	// Test 2: Dot Product
	// v(1, 0, 0) . v(0, 1, 0) = 0 (Orthogonal)
	// v(2, 0, 0) . v(2, 0, 0) = 4
	fmt.Println("\nTest 2: Dot Product")
	code2 := []bytecode.Instruction{
		{Op: bytecode.OP_PUSH, Args: 2.0},
		{Op: bytecode.OP_PUSH, Args: 0.0},
		{Op: bytecode.OP_PUSH, Args: 0.0},
		{Op: bytecode.OP_VEC3},

		{Op: bytecode.OP_PUSH, Args: 2.0},
		{Op: bytecode.OP_PUSH, Args: 0.0},
		{Op: bytecode.OP_PUSH, Args: 0.0},
		{Op: bytecode.OP_VEC3},

		{Op: bytecode.OP_VMUL}, // Dot product returns a float
		{Op: bytecode.OP_HALT},
	}
	runFloatTest(code2, 4.0)

	// Test 3: Scalar Multiplication
	// v(1, 2, 3) * 2.0 = v(2, 4, 6)
	fmt.Println("\nTest 3: Scalar Scale")
	code3 := []bytecode.Instruction{
		{Op: bytecode.OP_PUSH, Args: 1.0},
		{Op: bytecode.OP_PUSH, Args: 2.0},
		{Op: bytecode.OP_PUSH, Args: 3.0},
		{Op: bytecode.OP_VEC3}, // Stack: [v]

		{Op: bytecode.OP_PUSH, Args: 2.0}, // Stack: [v, 2.0]
		{Op: bytecode.OP_VSCALE},          // Stack: [v * 2.0]
		{Op: bytecode.OP_HALT},
	}
	runVectorTest(code3, "v3d(2.000000, 4.000000, 6.000000)")

	// Test 4: Magnitude
	// |v(3, 4, 0)| = 5
	fmt.Println("\nTest 4: Magnitude")
	code4 := []bytecode.Instruction{
		{Op: bytecode.OP_PUSH, Args: 3.0},
		{Op: bytecode.OP_PUSH, Args: 4.0},
		{Op: bytecode.OP_PUSH, Args: 0.0},
		{Op: bytecode.OP_VEC3},

		{Op: bytecode.OP_VMAG},
		{Op: bytecode.OP_HALT},
	}
	runFloatTest(code4, 5.0)

	fmt.Println("\n=== All Tests Complete ===")
}

func runVectorTest(code []bytecode.Instruction, expectedStr string) {
	v := vm.New(code)
	err := v.Run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	result, _ := v.GetResult()

	// Check type
	if !result.IsVec3() {
		fmt.Printf("FAIL: Expected Vector, got %v\n", result.Type)
		return
	}

	str := result.String()
	if str == expectedStr {
		fmt.Printf("PASS: %s\n", str)
	} else {
		fmt.Printf("FAIL: Expected %s, got %s\n", expectedStr, str)
	}
}

func runFloatTest(code []bytecode.Instruction, expected float64) {
	v := vm.New(code)
	err := v.Run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	result, _ := v.GetResult()

	if !result.IsFloat() {
		fmt.Printf("FAIL: Expected Float, got %v\n", result.Type)
		return
	}

	if result.AsFloat() == expected {
		fmt.Printf("PASS: %.2f\n", result.AsFloat())
	} else {
		fmt.Printf("FAIL: Expected %.2f, got %.2f\n", expected, result.AsFloat())
	}
}
