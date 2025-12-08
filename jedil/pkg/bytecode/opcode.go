package bytecode

// Opcode represents a single bytecode operation.type
type OpCode uint8

const (
	// Stack operations
	OP_PUSH OpCode = iota // push a constant value onto stack
	OP_POP // pop value from stack

	// Arithmetic operations
	OP_ADD // pop 2 values, add them, push results
	OP_SUB // pop 2 values, subtract them, push results
	OP_MUL // pop 2 values, multiply them, push results
	OP_DIV // pop 2 values, divide them, push results

	// Control flow operations
	OP_HALT // stop execution
)

// Instruction represents a complete bytecode instruction
// Some instructions may require an argument (e.g., OP_PUSH)
type Instruction struct {
	Op   OpCode // op to perform
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
	default:
		return "UNKNOWN_OPCODE"
	}
}

