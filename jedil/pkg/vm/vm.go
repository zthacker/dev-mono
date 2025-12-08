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
		case bytecode.OP_ADD:
			// pop 2 values, add them, push result
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("ADD failed: %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("ADD failed: %v", err)
			}

			// type check
			if !a.IsFloat() || !b.IsFloat() {
				return fmt.Errorf("ADD requires two float values")
			}

			result := NewFloat(a.AsFloat() + b.AsFloat())
			err = vm.stack.Push(result)
			if err != nil {
				return fmt.Errorf("ADD failed: %v", err)
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

			if !a.IsFloat() || !b.IsFloat() {
				return fmt.Errorf("SUB requires float operands")
			}

			result := a.AsFloat() - b.AsFloat()
			vm.stack.Push(NewFloat(result))

		case bytecode.OP_MUL:
			b, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("MUL failed (operand b): %v", err)
			}
			a, err := vm.stack.Pop()
			if err != nil {
				return fmt.Errorf("MUL failed (operand a): %v", err)
			}

			if !a.IsFloat() || !b.IsFloat() {
				return fmt.Errorf("MUL requires float operands")
			}

			result := a.AsFloat() * b.AsFloat()
			vm.stack.Push(NewFloat(result))

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
