package bytecode

// Opcode represents a single bytecode operation.type
type OpCode uint8

const (
	// Stack operations
	OP_PUSH OpCode = iota // push a constant value onto stack
	OP_POP                // pop value from stack

	// Arithmetic operations
	OP_ADD // pop 2 values, add them, push results
	OP_SUB // pop 2 values, subtract them, push results
	OP_MUL // pop 2 values, multiply them, push results
	OP_DIV // pop 2 values, divide them, push results

	// Vector operations
	OP_VADD   // Start at a new range (optional, just keeps it clean)
	OP_VSUB   // Vector subtraction
	OP_VMUL   // Vector dot product (returns float)
	OP_VSCALE // Vector scalar multiplication (vec * float)
	OP_VCROSS // Vector cross product
	OP_VMAG   // Vector magnitude (returns float)
	OP_VEC3   // Pop 3 floats and push as vec3

	// Batch (SIMD) Operations
	OP_BATCH_PACK   // Pop 4 Vec3s, Push 1 batch
	OP_BATCH_VADD   // Batch + Batch
	OP_BATCH_VSUB   // Batch - Batch
	OP_BATCH_VMUL   // Batch . Batch (Dot Product)
	OP_BATCH_VSCALE // Batch * ScalarArray

	// Control flow operations
	OP_HALT // stop execution
)

// Instruction represents a complete bytecode instruction
// Some instructions may require an argument (e.g., OP_PUSH)
type Instruction struct {
	Op   OpCode  // op to perform
	Args float64 // arg for instructions that need it
}

// String returns a human-readable representation of the OpCode
func (op OpCode) String() string {
	switch op {
	case OP_PUSH:
		return "OP_PUSH"
	case OP_POP:
		return "OP_POP"
	case OP_ADD:
		return "OP_ADD"
	case OP_SUB:
		return "OP_SUB"
	case OP_MUL:
		return "OP_MUL"
	case OP_DIV:
		return "OP_DIV"
	case OP_HALT:
		return "OP_HALT"
	case OP_VADD:
		return "VADD"
	case OP_VSUB:
		return "VSUB"
	case OP_VMUL:
		return "VMUL" // Dot product
	case OP_VSCALE:
		return "VSCALE"
	case OP_VCROSS:
		return "VCROSS"
	case OP_VMAG:
		return "VMAG"
	case OP_BATCH_PACK:
		return "BATCH_PACK"
	case OP_BATCH_VADD:
		return "BATCH_VADD"
	case OP_BATCH_VSUB:
		return "BATCH_VSUB"
	case OP_BATCH_VMUL:
		return "BATCH_VMUL"
	case OP_BATCH_VSCALE:
		return "BATCH_VSCALE"
	default:
		return "UNKNOWN_OPCODE"
	}
}
