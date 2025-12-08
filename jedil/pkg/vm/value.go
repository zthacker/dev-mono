package vm

import (
	"fmt"
	"jedil/pkg/types"
	"unsafe"
)

type ValueType uint8

const (
	TYPE_FLOAT ValueType = iota // floating-point number
	TYPE_BOOL                   // boolean value
	TYPE_NIL                    // nil value
	TYPE_VEC3                   // 3D vector
)

// Value represents a runtime value in the VM
// Small values (float, bool) are stored in Data
// Large values (vectors) are stored via Ptr
type Value struct {
	Data float64        // 8 bytes: for scalar types (float, bool)
	Ptr  unsafe.Pointer // 8 bytes: for complex types (e.g., Vec3) that can use a generic pointer
	Type ValueType      // 1 byte: type of the value
	pad  [7]byte        // padding to align to 24 bytes
}

// ============= Constructors =============

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

func NewVec3(v types.Vec3) Value {
	// allocate Vec3 on heap and store pointer
	vec := new(types.Vec3)
	*vec = v
	return Value{Type: TYPE_VEC3, Ptr: unsafe.Pointer(vec)}
}

// ============= Type Checking =============

func (v Value) IsFloat() bool {
	return v.Type == TYPE_FLOAT
}

func (v Value) IsBool() bool {
	return v.Type == TYPE_BOOL
}

func (v Value) IsNil() bool {
	return v.Type == TYPE_NIL
}

func (v Value) IsVec3() bool {
	return v.Type == TYPE_VEC3
}

// ============= Conversions =============

func (v Value) AsFloat() float64 {
	return v.Data
}

func (v Value) AsBool() bool {
	return v.Data != 0.0
}

func (v Value) AsVec3() types.Vec3 {
	if v.Type != TYPE_VEC3 {
		return types.Vec3{} // return zero vector if not Vec3
	}
	return *(*types.Vec3)(v.Ptr)
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
	case TYPE_VEC3:
		vec := v.AsVec3()
		return vec.String()
	default:
		return "unknown"
	}
}
