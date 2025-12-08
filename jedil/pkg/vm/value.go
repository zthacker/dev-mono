package vm

import "fmt"

type ValueType uint8

const (
	TYPE_FLOAT ValueType = iota // floating-point number
	TYPE_BOOL                   // boolean value
	TYPE_NIL                    // nil value
)

type Value struct {
	Type ValueType
	Data float64 // using float64 to represent both numbers and booleans (0.0 = false, 1.0 = true)
}

func NewFloat(f float64) Value {
	return Value{Type: TYPE_FLOAT, Data: f}
}

func NewBool(b bool) Value {
	if b {
		return Value{Type: TYPE_BOOL, Data: 1.0}
	} else {
		return Value{Type: TYPE_BOOL, Data: 0.0}
	}
}

func NewNil() Value {
	return Value{Type: TYPE_NIL, Data: 0.0}
}

// Helpers
func (v Value) IsFloat() bool {
	return v.Type == TYPE_FLOAT
}

func (v Value) IsBool() bool {
	return v.Type == TYPE_BOOL
}

func (v Value) IsNil() bool {
	return v.Type == TYPE_NIL
}

// Conversion helpers
func (v Value) AsFloat() float64 {
	return v.Data
}

func (v Value) AsBool() bool {
	return v.Data != 0.0
}

// String representation for debugging
func (v Value) String() string {
	switch v.Type {
	case TYPE_FLOAT:
		return fmt.Sprintf("%.6f", v.Data)
	case TYPE_BOOL:
		if v.AsBool() {
			return "true"
		}
		return "false"
	case TYPE_NIL:
		return "nil"
	default:
		return "unknown"
	}
}
