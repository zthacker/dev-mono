package main

import (
	"fmt"
	"jedil/pkg/bytecode"
	"jedil/pkg/vm"
)

func main() {
	fmt.Println("=== JEDIL SIMD Batch Test ===\n")

	// Goal: Add 4 vectors to 4 other vectors in ONE instruction
	// Batch A: [v(1,1,1), v(2,2,2), v(3,3,3), v(4,4,4)]
	// Batch B: [v(10,10,10), v(20,20,20), v(30,30,30), v(40,40,40)]
	// Result:  [v(11,11,11), v(22,22,22), v(33,33,33), v(44,44,44)]

	code := []bytecode.Instruction{
		// --- Create Batch A ---
		// Push 4 vectors onto the stack
		// Vector 0
		{Op: bytecode.OP_PUSH, Args: 1.0}, {Op: bytecode.OP_PUSH, Args: 1.0}, {Op: bytecode.OP_PUSH, Args: 1.0}, {Op: bytecode.OP_VEC3},
		// Vector 1
		{Op: bytecode.OP_PUSH, Args: 2.0}, {Op: bytecode.OP_PUSH, Args: 2.0}, {Op: bytecode.OP_PUSH, Args: 2.0}, {Op: bytecode.OP_VEC3},
		// Vector 2
		{Op: bytecode.OP_PUSH, Args: 3.0}, {Op: bytecode.OP_PUSH, Args: 3.0}, {Op: bytecode.OP_PUSH, Args: 3.0}, {Op: bytecode.OP_VEC3},
		// Vector 3
		{Op: bytecode.OP_PUSH, Args: 4.0}, {Op: bytecode.OP_PUSH, Args: 4.0}, {Op: bytecode.OP_PUSH, Args: 4.0}, {Op: bytecode.OP_VEC3},

		// Stack has: [v0, v1, v2, v3]
		// Compress them into one Batch value
		{Op: bytecode.OP_BATCH_PACK}, // Stack: [BatchA]

		// --- Create Batch B ---
		// Vector 0
		{Op: bytecode.OP_PUSH, Args: 10.0}, {Op: bytecode.OP_PUSH, Args: 10.0}, {Op: bytecode.OP_PUSH, Args: 10.0}, {Op: bytecode.OP_VEC3},
		// Vector 1
		{Op: bytecode.OP_PUSH, Args: 20.0}, {Op: bytecode.OP_PUSH, Args: 20.0}, {Op: bytecode.OP_PUSH, Args: 20.0}, {Op: bytecode.OP_VEC3},
		// Vector 2
		{Op: bytecode.OP_PUSH, Args: 30.0}, {Op: bytecode.OP_PUSH, Args: 30.0}, {Op: bytecode.OP_PUSH, Args: 30.0}, {Op: bytecode.OP_VEC3},
		// Vector 3
		{Op: bytecode.OP_PUSH, Args: 40.0}, {Op: bytecode.OP_PUSH, Args: 40.0}, {Op: bytecode.OP_PUSH, Args: 40.0}, {Op: bytecode.OP_VEC3},

		{Op: bytecode.OP_BATCH_PACK}, // Stack: [BatchA, BatchB]

		// --- The SIMD Magic ---
		{Op: bytecode.OP_BATCH_VADD}, // Adds 12 floats (4 vectors) in one go!
		{Op: bytecode.OP_HALT},
	}

	v := vm.New(code)
	err := v.Run()
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	result, _ := v.GetResult()

	// Verify result is a Batch
	if !result.IsVec3Batch() {
		fmt.Printf("❌ FAIL: Expected Batch, got %v\n", result.Type)
		return
	}

	fmt.Println("Result Batch:")
	fmt.Println(result.String())

	// Manual check of first value
	batch := result.AsVec3Batch()
	first := batch.Get(0)
	if first.X == 11.0 {
		fmt.Println("✅ PASS: SIMD Addition verified!")
	} else {
		fmt.Printf("❌ FAIL: Expected 11.0, got %f\n", first.X)
	}
}
