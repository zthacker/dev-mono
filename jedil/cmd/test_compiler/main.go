package main

import (
	"fmt"
	"jedil/pkg/compiler"
	"jedil/pkg/vm"
)

func main() {
	fmt.Println("=== JEDIL Compiler Test ===\n")

	// Test 1: Simple vector addition
	fmt.Println("Test 1: Vector Addition")
	source1 := `
		return vec3(1, 2, 3) + vec3(4, 5, 6)
	`

	instructions1, err := compiler.CompileSource(source1)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")
		fmt.Printf("   Generated %d instructions\n", len(instructions1))

		// Execute
		v := vm.New(instructions1)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsVec3() {
				vec := result.AsVec3()
				if vec.X == 5.0 && vec.Y == 7.0 && vec.Z == 9.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Println("   FAIL: incorrect result\n")
				}
			} else {
				fmt.Println("   FAIL: not a vector\n")
			}
		}
	}

	// Test 2: Vector cross product
	fmt.Println("Test 2: Cross Product")
	source2 := `
		return cross(vec3(1, 0, 0), vec3(0, 1, 0))
	`

	instructions2, err := compiler.CompileSource(source2)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")

		v := vm.New(instructions2)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsVec3() {
				vec := result.AsVec3()
				// i × j = k (0, 0, 1)
				if vec.X == 0.0 && vec.Y == 0.0 && vec.Z == 1.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Println("   FAIL: incorrect result\n")
				}
			} else {
				fmt.Println("   FAIL: not a vector\n")
			}
		}
	}

	// Test 3: Dot product
	fmt.Println("Test 3: Dot Product")
	source3 := `
		return dot(vec3(1, 2, 3), vec3(4, 5, 6))
	`

	instructions3, err := compiler.CompileSource(source3)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")

		v := vm.New(instructions3)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsFloat() {
				// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
				if result.AsFloat() == 32.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Printf("   FAIL: expected 32.0, got %f\n\n", result.AsFloat())
				}
			} else {
				fmt.Println("   FAIL: not a float\n")
			}
		}
	}

	// Test 4: Vector magnitude
	fmt.Println("Test 4: Vector Magnitude")
	source4 := `
		return mag(vec3(3, 4, 0))
	`

	instructions4, err := compiler.CompileSource(source4)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")

		v := vm.New(instructions4)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsFloat() {
				// sqrt(3^2 + 4^2) = sqrt(9 + 16) = 5
				if result.AsFloat() == 5.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Printf("   FAIL: expected 5.0, got %f\n\n", result.AsFloat())
				}
			} else {
				fmt.Println("   FAIL: not a float\n")
			}
		}
	}

	// Test 5: Arithmetic expression
	fmt.Println("Test 5: Scalar Arithmetic")
	source5 := `
		return 2.0 + 3.0 * 4.0
	`

	instructions5, err := compiler.CompileSource(source5)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")

		v := vm.New(instructions5)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsFloat() {
				if result.AsFloat() == 14.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Printf("   FAIL: expected 14.0, got %f\n\n", result.AsFloat())
				}
			} else {
				fmt.Println("   FAIL: not a float\n")
			}
		}
	}

	fmt.Println("=== All Compiler Tests Complete! ===")

	fmt.Println("Test 6: Variable Load")
	source6 := `
		let x = 5.0
		let y = 10.0
		return x + y
	
	`

	instructions6, err := compiler.CompileSource(source6)
	if err != nil {
		fmt.Printf("Compilation failed: %v\n\n", err)
	} else {
		fmt.Println("Compilation succeeded")

		v := vm.New(instructions6)
		if err := v.Run(); err != nil {
			fmt.Printf("Execution failed: %v\n\n", err)
		} else {
			result, _ := v.GetResult()
			fmt.Printf("   Result: %s\n", result.String())
			if result.IsFloat() {
				if result.AsFloat() == 15.0 {
					fmt.Println("   PASS\n")
				} else {
					fmt.Printf("   FAIL: expected 15.0, got %f\n\n", result.AsFloat())
				}
			} else {
				fmt.Println("   FAIL: not a float\n")
			}
		}
	}

}
