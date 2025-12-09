package vm

import (
	"fmt"
	"jedil/pkg/bytecode"
)

// VM represents the virtual machine.
type VM struct {
	stack *Stack                 // the VM stack
	code  []bytecode.Instruction // bytecode instructions to execute
	ip    int                    // instruction pointer
}

// NewVM creates and initializes a new VM with the given bytecode.
func New(code []bytecode.Instruction) *VM {
	return &VM{
		stack: NewStack(),
		code:  code,
		ip:    0,
	}
}

// Run executes the VM.
func (vm *VM) Run() error {
	for vm.ip < len(vm.code) {
		// fetch current instruction
		inst := vm.code[vm.ip]
		vm.ip++

		// execute instruction
		switch inst.Op {
		case bytecode.OP_PUSH:
			// push constant value onto stack
			err := vm.stack.Push(NewFloat(inst.Args))
			if err != nil {
				return fmt.Errorf("PUSH failed: %v", err)
			}
		case bytecode.OP_POP:
			// pop value from stack
			_, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("POP failed: %v", err)
			}
		case bytecode.OP_LOAD:
			// load variable from stack (Args is index)
			// get the offset
			offset := int(inst.Args)

			// get the value at the offset -- bounds check dont in Get()
			val, err := vm.stack.Get(offset)
			if err != nil {
				return fmt.Errorf("LOAD failed: %v", err)
			}

			// push the value onto the top of the stack
			err = vm.stack.Push(val)
			if err != nil {
				return fmt.Errorf("LOAD failed: %v", err)
			}
		case bytecode.OP_ADD:
			// pop 2 values, add them, push result (polymorphic: float or vec3)
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("ADD failed: %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("ADD failed: %v", err)
			}

			// Polymorphic: handle both floats and vectors
			if a.IsFloat() && b.IsFloat() {
				result := NewFloat(a.AsFloat() + b.AsFloat())
				err = vm.stack.Push(result)
				if err != nil {
					return fmt.Errorf("ADD failed: %v", err)
				}
			} else if a.IsVec3() && b.IsVec3() {
				// Vector addition
				v1 := a.AsVec3()
				v2 := b.AsVec3()
				result := NewVec3(v1.Add(v2))
				err = vm.stack.Push(result)
				if err != nil {
					return fmt.Errorf("ADD failed: %v", err)
				}
			} else {
				return fmt.Errorf("ADD requires two floats or two vec3s")
			}
		case bytecode.OP_SUB:
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("SUB failed (operand b): %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("SUB failed (operand a): %v", err)
			}

			// Polymorphic: handle both floats and vectors
			if a.IsFloat() && b.IsFloat() {
				result := a.AsFloat() - b.AsFloat()
				vm.stack.Push(NewFloat(result))
			} else if a.IsVec3() && b.IsVec3() {
				v1 := a.AsVec3()
				v2 := b.AsVec3()
				result := NewVec3(v1.Sub(v2))
				vm.stack.Push(result)
			} else {
				return fmt.Errorf("SUB requires two floats or two vec3s")
			}

		case bytecode.OP_MUL:
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("MUL failed (operand b): %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("MUL failed (operand a): %v", err)
			}

			// Polymorphic: scalar * scalar, or vec3 * scalar (scaling)
			if a.IsFloat() && b.IsFloat() {
				result := a.AsFloat() * b.AsFloat()
				vm.stack.Push(NewFloat(result))
			} else if a.IsVec3() && b.IsFloat() {
				// vec * scalar
				v := a.AsVec3()
				s := b.AsFloat()
				result := NewVec3(v.Scale(s))
				vm.stack.Push(result)
			} else if a.IsFloat() && b.IsVec3() {
				// scalar * vec
				s := a.AsFloat()
				v := b.AsVec3()
				result := NewVec3(v.Scale(s))
				vm.stack.Push(result)
			} else {
				return fmt.Errorf("MUL requires floats or vec3*float")
			}

		case bytecode.OP_DIV:
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("DIV failed (operand b): %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("DIV failed (operand a): %v", err)
			}

			if !a.IsFloat() || !b.IsFloat() {
				return fmt.Errorf("DIV requires float operands")
			}

			if b.AsFloat() == 0.0 {
				return fmt.Errorf("division by zero")
			}

			result := a.AsFloat() / b.AsFloat()
			vm.stack.Push(NewFloat(result))
		case bytecode.OP_VADD:
			if err := vm.opVAdd(); err != nil {
				return fmt.Errorf("VADD failed: %v", err)
			}

		case bytecode.OP_VSUB:
			if err := vm.opVSub(); err != nil {
				return fmt.Errorf("VSUB failed: %v", err)
			}

		case bytecode.OP_VMUL:
			if err := vm.opVMul(); err != nil {
				return fmt.Errorf("VMUL failed: %v", err)
			}

		case bytecode.OP_VSCALE:
			if err := vm.opVScale(); err != nil {
				return fmt.Errorf("VSCALE failed: %v", err)
			}

		case bytecode.OP_VCROSS:
			if err := vm.opVCross(); err != nil {
				return fmt.Errorf("VCROSS failed: %v", err)
			}

		case bytecode.OP_VMAG:
			if err := vm.opVMag(); err != nil {
				return fmt.Errorf("VMAG failed: %v", err)
			}
		case bytecode.OP_VEC3:
			if err := vm.opVec3(); err != nil {
				return fmt.Errorf("VEC3 failed: %v", err)
			}
		case bytecode.OP_BATCH_PACK:
			if err := vm.opBatchPack(); err != nil {
				return fmt.Errorf("BATCH_PACK failed: %v", err)
			}
		case bytecode.OP_BATCH_VADD:
			if err := vm.opBatchVAdd(); err != nil {
				return fmt.Errorf("BATCH_VADD failed: %v", err)
			}
		case bytecode.OP_BATCH_VSUB:
			if err := vm.opBatchVSub(); err != nil {
				return fmt.Errorf("BATCH_VSUB failed: %v", err)
			}
		case bytecode.OP_BATCH_VMUL:
			if err := vm.opBatchVMul(); err != nil {
				return fmt.Errorf("BATCH_VMUL failed: %v", err)
			}
		case bytecode.OP_HALT:
			// Stop execution
			return nil

		default:
			return fmt.Errorf("unknown opcode: %d", inst.Op)
		}
	}

	return nil
}

// GetResult returns the current state of the VM stack.
func (vm *VM) GetResult() (Value, error) {
	return vm.stack.Peek()
}

func (vm *VM) Reset() {
	vm.stack.Reset()
	vm.ip = 0
}

func (vm *VM) Debug() string {
	return fmt.Sprintf("IP: %d, Stack: %s", vm.ip, vm.stack.String())
}
