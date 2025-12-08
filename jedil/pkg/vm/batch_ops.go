package vm

import (
	"fmt"
	"jedil/pkg/types"
)

func (vm *VM) popBatch() (types.Vec3Batch, error) {
	val, err := vm.stack.Pop()
	if err != nil {
		return types.Vec3Batch{}, err
	}
	if !val.IsVec3Batch() {
		return types.Vec3Batch{}, fmt.Errorf("expected Vec3Batch on stack, got %s", val.String())
	}
	return val.AsVec3Batch(), nil
}

// OP_BATCH_PACK: Pop 4 Vec3s, Push 1 Vec3Batch
// Stack before: [v0, v1, v2, v3] (top is v3)
// Stack after: [Batch{v0, v1, v2, v3}]
func (vm *VM) opBatchPack() error {

	batch := types.NewVec3Batch()

	// Pop in reverse order...

	// v3
	v3, err := vm.popVector()
	if err != nil {
		return err
	}
	batch.Set(3, v3.X, v3.Y, v3.Z)

	// v2
	v2, err := vm.popVector()
	if err != nil {
		return err
	}
	batch.Set(2, v2.X, v2.Y, v2.Z)

	// v1
	v1, err := vm.popVector()
	if err != nil {
		return err
	}
	batch.Set(1, v1.X, v1.Y, v1.Z)

	// v0
	v0, err := vm.popVector()
	if err != nil {
		return err
	}
	batch.Set(0, v0.X, v0.Y, v0.Z)

	// Push batch
	return vm.stack.Push(NewVec3Batch(batch))
}

// OP_BATCH_VADD: Batch + Batch
func (vm *VM) opBatchVAdd() error {
	b, err := vm.popBatch()
	if err != nil {
		return err
	}
	a, err := vm.popBatch()
	if err != nil {
		return err
	}

	result := a.Add(b)
	return vm.stack.Push(NewVec3Batch(result))
}

// OP_BATCH_VSUB: Batch - Batch
func (vm *VM) opBatchVSub() error {
	b, err := vm.popBatch()
	if err != nil {
		return err
	}
	a, err := vm.popBatch()
	if err != nil {
		return err
	}

	result := a.Sub(b)
	return vm.stack.Push(NewVec3Batch(result))
}

func (vm *VM) opBatchVMul() error {
	b, err := vm.popBatch()
	if err != nil {
		return err
	}
	a, err := vm.popBatch()
	if err != nil {
		return err
	}

	dots := a.Dot(b)

	// push back onto the stack: d0, d1, d2, d3
	for i := 0; i < 4; i++ {
		vm.stack.Push(NewFloat(dots[i]))
	}

	return nil
}
