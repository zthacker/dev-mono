package main

import (
	"fmt"
	"jedil/pkg/compiler"
	"jedil/pkg/vm"
)

func runTest(name string, source string, expected string) {
	fmt.Printf("\nTest: %s\n", name)
	fmt.Printf("Source:\n%s\n", source)

	// Compile
	instructions, err := compiler.CompileSource(source)
	if err != nil {
		fmt.Printf("Compile Error: %v\n", err)
		return
	}

	fmt.Printf("Compiled: %d instructions\n", len(instructions))

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
	fmt.Printf("   Expected: %s\n", expected)
}

func main() {
	fmt.Println("========================================")
	fmt.Println("   JEDIL FUNCTION TESTS")
	fmt.Println("========================================")

	// Test 1: Simple function call
	runTest(
		"Simple function call",
		`fn double(x) {
    return x * 2
}
return double(5)`,
		"10.0",
	)

	// Test 2: Multiple parameters
	runTest(
		"Multiple parameters",
		`fn add3(a, b, c) {
    return a + b + c
}
return add3(1, 2, 3)`,
		"6.0",
	)

	// Test 3: Function with local variables
	runTest(
		"Function with local variables",
		`fn compute(x) {
    let a = x + 1
    let b = a * 2
    return b
}
return compute(5)`,
		"12.0",
	)

	// Test 4: Nested function calls
	runTest(
		"Nested function calls",
		`fn double(x) {
    return x * 2
}
fn quad(x) {
    return double(double(x))
}
return quad(3)`,
		"12.0",
	)

	// Test 5: Function with arithmetic
	runTest(
		"Function with arithmetic",
		`fn average(a, b) {
    return (a + b) / 2
}
return average(10, 20)`,
		"15.0",
	)

	// Test 6: Function calling built-in
	runTest(
		"Function calling built-in (mag)",
		`fn magnitude(v) {
    return mag(v)
}
return magnitude(vec3(3, 4, 0))`,
		"5.0",
	)

	// Test 7: Vector operations in functions
	runTest(
		"Vector operations in functions",
		`fn scale_vec(v, s) {
    return v * s
}
return scale_vec(vec3(1, 2, 3), 2)`,
		"vec3(2.00, 4.00, 6.00)",
	)

	// Test 8: Multiple function definitions
	runTest(
		"Multiple function definitions",
		`fn add(a, b) {
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
return calc()`,
		"27.0",
	)

	// Test 9: Function with no parameters
	runTest(
		"Function with no parameters",
		`fn getAnswer() {
    return 42
}
return getAnswer()`,
		"42.0",
	)

	// Test 10: Implicit return (function without explicit return)
	runTest(
		"Implicit return",
		`fn noReturn() {
    let x = 5
}
return noReturn()`,
		"0.0",
	)

	fmt.Println("\n========================================")
	fmt.Println("   ALL TESTS COMPLETE")
	fmt.Println("========================================")
}
