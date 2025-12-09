package main

/*
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"jedil/pkg/bytecode"
	"jedil/pkg/compiler"
	"jedil/pkg/types"
	"jedil/pkg/vm"
	"os"
	"unsafe"
)

// ============================================================================
// VM Registry (to avoid passing Go pointers to C)
// ============================================================================

var vmRegistry = make(map[int]*vm.VM)
var nextVMHandle = 1

func registerVM(v *vm.VM) unsafe.Pointer {
	handle := nextVMHandle
	nextVMHandle++
	vmRegistry[handle] = v
	return unsafe.Pointer(uintptr(handle))
}

func getVM(handle unsafe.Pointer) *vm.VM {
	h := int(uintptr(handle))
	return vmRegistry[h]
}

func unregisterVM(handle unsafe.Pointer) {
	h := int(uintptr(handle))
	delete(vmRegistry, h)
}

// ============================================================================
// Error Handling
// ============================================================================

var lastError string

//export jedil_get_last_error
func jedil_get_last_error() *C.char {
	return C.CString(lastError)
}

func setError(err error) {
	if err != nil {
		lastError = err.Error()
	} else {
		lastError = ""
	}
}

// ============================================================================
// Program Lifecycle
// ============================================================================

//export jedil_create_program
func jedil_create_program(bytecode_data *C.uint8_t, length C.size_t) unsafe.Pointer {
	// Convert C bytecode to Go bytecode
	cBytes := C.GoBytes(unsafe.Pointer(bytecode_data), C.int(length))

	// Parse bytecode (simple format: each instruction is 9 bytes: 1 opcode + 8 arg)
	var instructions []bytecode.Instruction
	for i := 0; i < len(cBytes); i += 9 {
		if i+9 > len(cBytes) {
			break
		}
		op := bytecode.OpCode(cBytes[i])
		// Decode float64 argument (little-endian)
		arg := binary.LittleEndian.Uint64(cBytes[i+1 : i+9])
		argFloat := *(*float64)(unsafe.Pointer(&arg))

		instructions = append(instructions, bytecode.Instruction{
			Op:   op,
			Args: argFloat,
		})
	}

	// Create VM
	v := vm.New(instructions)

	// Register and return handle
	return registerVM(v)
}

//export jedil_compile_file
func jedil_compile_file(filepath *C.char) unsafe.Pointer {
	// Convert C string to Go string
	path := C.GoString(filepath)

	// Read the file
	sourceBytes, err := os.ReadFile(path)
	if err != nil {
		setError(fmt.Errorf("failed to read file %s: %v", path, err))
		return nil
	}

	// Compile the source
	source := string(sourceBytes)
	instructions, err := compiler.CompileSource(source)
	if err != nil {
		setError(fmt.Errorf("compilation failed: %v", err))
		return nil
	}

	// Create VM
	v := vm.New(instructions)
	setError(nil)
	return registerVM(v)
}

//export jedil_compile_source
func jedil_compile_source(sourceStr *C.char) unsafe.Pointer {
	// Convert C string to Go string
	source := C.GoString(sourceStr)

	// Compile the source
	instructions, err := compiler.CompileSource(source)
	if err != nil {
		setError(fmt.Errorf("compilation failed: %v", err))
		return nil
	}

	// Create VM
	v := vm.New(instructions)
	setError(nil)
	return registerVM(v)
}

//export jedil_free_program
func jedil_free_program(program unsafe.Pointer) {
	unregisterVM(program)
}

// ============================================================================
// Execution Functions
// ============================================================================

//export jedil_execute_vec3
func jedil_execute_vec3(program unsafe.Pointer, input_data unsafe.Pointer, input_len C.size_t, result_x *C.double, result_y *C.double, result_z *C.double) C.int {
	v := getVM(program)

	// Execute
	err := v.Run()
	if err != nil {
		setError(err)
		return 3 // JEDIL_ERROR_EXECUTION_FAILED
	}

	// Get result
	val, err := v.GetResult()
	if err != nil {
		setError(err)
		return 4 // JEDIL_ERROR_STACK_UNDERFLOW
	}

	if !val.IsVec3() {
		setError(fmt.Errorf("expected Vec3, got %s", val.String()))
		return 5 // JEDIL_ERROR_TYPE_MISMATCH
	}

	vec := val.AsVec3()
	*result_x = C.double(vec.X)
	*result_y = C.double(vec.Y)
	*result_z = C.double(vec.Z)

	setError(nil)
	return 0 // JEDIL_OK
}

//export jedil_execute_float
func jedil_execute_float(program unsafe.Pointer, input_data unsafe.Pointer, input_len C.size_t, result *C.double) C.int {
	v := getVM(program)

	err := v.Run()
	if err != nil {
		setError(err)
		return 3
	}

	val, err := v.GetResult()
	if err != nil {
		setError(err)
		return 4
	}

	if !val.IsFloat() {
		setError(fmt.Errorf("expected float, got %s", val.String()))
		return 5
	}

	*result = C.double(val.AsFloat())
	setError(nil)
	return 0
}

// ============================================================================
// Convenience Functions (No VM, Direct Calls)
// ============================================================================

//export jedil_vec3_add
func jedil_vec3_add(ax, ay, az C.double, bx, by, bz C.double, result_x, result_y, result_z *C.double) {
	*result_x = ax + bx
	*result_y = ay + by
	*result_z = az + bz
}

//export jedil_batch_add
func jedil_batch_add(
	a_xs, a_ys, a_zs *C.double,
	b_xs, b_ys, b_zs *C.double,
	result_xs, result_ys, result_zs *C.double,
) {
	// Convert C arrays to Go
	var goA types.Vec3Batch
	var goB types.Vec3Batch

	// Copy input data (4 elements each)
	for i := 0; i < 4; i++ {
		goA.Xs[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(a_xs)) + uintptr(i)*8)))
		goA.Ys[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(a_ys)) + uintptr(i)*8)))
		goA.Zs[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(a_zs)) + uintptr(i)*8)))

		goB.Xs[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(b_xs)) + uintptr(i)*8)))
		goB.Ys[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(b_ys)) + uintptr(i)*8)))
		goB.Zs[i] = float64(*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(b_zs)) + uintptr(i)*8)))
	}

	// Perform addition (SIMD!)
	goResult := goA.Add(goB)

	// Copy results back
	for i := 0; i < 4; i++ {
		*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(result_xs)) + uintptr(i)*8)) = C.double(goResult.Xs[i])
		*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(result_ys)) + uintptr(i)*8)) = C.double(goResult.Ys[i])
		*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(result_zs)) + uintptr(i)*8)) = C.double(goResult.Zs[i])
	}
}

// ============================================================================
// Main (required for building shared library)
// ============================================================================

func main() {
	// Empty main required for Go build
}
