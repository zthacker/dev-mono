package vm

import (
	"jedil/pkg/bytecode"
	"jedil/pkg/types"
	"testing"
)

// -----------------------------------------------------------------------------
// 1. VM EXECUTION BENCHMARKS (Full System Overhead)
// -----------------------------------------------------------------------------

// BenchmarkVMScalar measures the VM executing 4 separate VADD instructions
func BenchmarkVMScalar(b *testing.B) {
	// Pre-load the stack so we only measure the ADD instruction
	// We need 8 vectors on stack (4 pairs) to do 4 ADDs
	// But to keep it simple and avoid underflow, we'll just re-use the VM setup
	// and reset the IP.

	code := []bytecode.Instruction{
		// 4 Adds
		{Op: bytecode.OP_VADD},
		{Op: bytecode.OP_VADD},
		{Op: bytecode.OP_VADD},
		{Op: bytecode.OP_VADD},
		{Op: bytecode.OP_HALT},
	}

	vm := New(code)

	// Pre-populate stack with 8 vectors so VADD has data
	// (In a real app, we'd push these, but here we cheat for speed)
	v := NewVec3(types.NewVec3(1, 2, 3))
	for i := 0; i < 8; i++ {
		vm.stack.Push(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.ip = 0        // Reset Instruction Pointer
		vm.stack.top = 8 // Reset Stack Pointer (pretend data is still there)
		vm.Run()
	}
}

// BenchmarkVMBatch measures the VM executing 1 BATCH_VADD instruction
func BenchmarkVMBatch(b *testing.B) {
	code := []bytecode.Instruction{
		// 1 Batch Add (does the work of 4 scalar adds)
		{Op: bytecode.OP_BATCH_VADD},
		{Op: bytecode.OP_HALT},
	}

	vm := New(code)

	// Pre-populate stack with 2 batches
	batch := types.NewVec3Batch()
	batch.Set(0, 1, 2, 3)
	vBatch := NewVec3Batch(batch)

	vm.stack.Push(vBatch)
	vm.stack.Push(vBatch)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.ip = 0
		vm.stack.top = 2
		vm.Run()
	}
}
