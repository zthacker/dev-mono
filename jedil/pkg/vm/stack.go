package vm

import "fmt"

// STACK_MAX defines the maximum stack size.
const STACK_MAX = 256

// Stack represents the VM stack.
type Stack struct {
	values [STACK_MAX]Value // fixed-size array to hold stack values
	top    int              // index of the next free slot in the stack
}

// NewStack creates and initializes a new Stack.
func NewStack() *Stack {
	return &Stack{top: 0}
}

// Push adds a value to the top of the stack.
func (s *Stack) Push(v Value) error {
	if s.top >= STACK_MAX {
		return fmt.Errorf("stack overflow")
	}

	s.values[s.top] = v
	s.top++
	return nil
}

// Pop removes and returns the value from the top of the stack.
func (s *Stack) Pop() (Value, error) {
	if s.top == 0 {
		return Value{}, fmt.Errorf("stack underflow")
	}

	s.top--
	return s.values[s.top], nil
}

// Peek returns the value at the top of the stack without removing it.
func (s *Stack) Peek() (Value, error) {
	if s.top == 0 {
		return Value{}, fmt.Errorf("stack underflow")
	}

	return s.values[s.top-1], nil
}

// Size returns the current number of elements in the stack.
func (s *Stack) Size() int {
	return s.top
}

// Reset clears the stack.
func (s *Stack) Reset() {
	s.top = 0
}

func (s *Stack) Get(offset int) (Value, error) {
	if offset < 0 || offset >= s.top {
		return Value{}, fmt.Errorf("stack get: index %d out of bounds", offset)
	}
	return s.values[offset], nil
}

// String returns a string representation of the stack for debugging.
func (s *Stack) String() string {
	if s.top == 0 {
		return "[]"
	}

	result := "["
	for i := 0; i < s.top; i++ {
		result += s.values[i].String()
		if i < s.top-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}
