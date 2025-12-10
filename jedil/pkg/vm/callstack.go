package vm

import "fmt"

const CALL_STACK_MAX = 64

type CallFrame struct {
	returnAddress int // IP to return to after
	basePointer int // Stack offset where functions params begin
	localCount int // Number of stack slots used (params + locals)
}

type CallStack struct {
	frames [CALL_STACK_MAX]CallFrame
	top int
}

func NewCallStack() *CallStack {
	return &CallStack{
		top: 0,
	}
}

func (c *CallStack) Push(frame CallFrame) error {
	if c.top >= CALL_STACK_MAX {
		return fmt.Errorf("call stack overflow")
	}

	c.frames[c.top] = frame
	c.top++

	return nil
}

func (c *CallStack) Pop() (CallFrame, error) {
	if c.top == 0 {
		return CallFrame{}, fmt.Errorf("stack underflow")
	}

	c.top--

	return c.frames[c.top], nil
}

func (c *CallStack) Peek() (CallFrame, error) {
	if c.top == 0 {
		return CallFrame{}, fmt.Errorf("stack underflow")
	}

	return c.frames[c.top-1], nil
}

