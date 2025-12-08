package vm

import (
	"fmt"
	"jedil/pkg/types"
)

// Pop a vector from the stack
func (vm *VM) popVector() (types.Vec3, error) {
	v, err := vm.stack.Pop()
	if err != nil {
		return types.Vec3{}, err
	}
	if !v.IsVec3() {
		return types.Vec3{}, fmt.Errorf("expected vector, got %v", v.Type)
	}
	return v.AsVec3(), nil
}

// Pop a float from the stack
func (vm *VM) popFloat() (float64, error) {
	v, err := vm.stack.Pop()
	if err != nil {
		return 0, err
	}
	if !v.IsFloat() {
		return 0, fmt.Errorf("expected float, got %v", v.Type)
	}
	return v.AsFloat(), nil
}

// OP_VADD: v1 + v2
func (vm *VM) opVAdd() error {
	b, err := vm.popVector() // Second operand
	if err != nil {
		return err
	}
	a, err := vm.popVector() // First operand
	if err != nil {
		return err
	}

	result := a.Add(b)
	return vm.stack.Push(NewVec3(result))
}

// OP_VSUB: v1 - v2
func (vm *VM) opVSub() error {
	b, err := vm.popVector()
	if err != nil {
		return err
	}
	a, err := vm.popVector()
	if err != nil {
		return err
	}

	result := a.Sub(b)
	return vm.stack.Push(NewVec3(result))
}

// OP_VMUL: Dot Product (v1 . v2) -> Float
func (vm *VM) opVMul() error {
	b, err := vm.popVector()
	if err != nil {
		return err
	}
	a, err := vm.popVector()
	if err != nil {
		return err
	}

	result := a.Dot(b)
	return vm.stack.Push(NewFloat(result))
}

// OP_VSCALE: Vector * Float -> Vector
func (vm *VM) opVScale() error {
	s, err := vm.popFloat() // Pop the scalar (top of stack)
	if err != nil {
		return err
	}
	v, err := vm.popVector() // Pop the vector
	if err != nil {
		return err
	}

	result := v.Scale(s)
	return vm.stack.Push(NewVec3(result))
}

// OP_VCROSS: Cross Product (v1 x v2) -> Vector
func (vm *VM) opVCross() error {
	b, err := vm.popVector()
	if err != nil {
		return err
	}
	a, err := vm.popVector()
	if err != nil {
		return err
	}

	result := a.Cross(b)
	return vm.stack.Push(NewVec3(result))
}

// OP_VMAG: Magnitude(v) -> Float
func (vm *VM) opVMag() error {
	v, err := vm.popVector()
	if err != nil {
		return err
	}

	result := v.Magnitude()
	return vm.stack.Push(NewFloat(result))
}

// OP_VEC3: Pop 3 floats, make Vec3
func (vm *VM) opVec3() error {
	z, err := vm.popFloat()
	if err != nil {
		return err
	}
	y, err := vm.popFloat()
	if err != nil {
		return err
	}
	x, err := vm.popFloat()
	if err != nil {
		return err
	}

	v := types.NewVec3(x, y, z)
	return vm.stack.Push(NewVec3(v))
}
