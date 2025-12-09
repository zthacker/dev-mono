# JEDIL - Hot-Reloadable Intermediate Language for Astrodynamics

**JEDIL** (Jedi Intermediate Language) is a specialized IL/VM for scientific computing. It is named after the my cat (Jedi) who recently passed. Write algorithms in `.jedil` scripts and reload them at runtime without recompiling your code.

## Why JEDIL?

### The Problem
In astrodynamics and scientific computing, you often need to:
- Experiment with different algorithms (MOID calculations, spline interpolation, orbit propagation)
- Tune parameters and test variations
- Fix bugs in production systems
- Share algorithms across C, Python, and other languages

## Quick Start

### 1. Build the library

```bash
cd jedil
./pkg/ffi/build_ffi.sh
```

This creates `libjedil.so`.

### 2. Write a `.jedil` script

```jedil
// vec_add.jedil
return vec3(1, 2, 3) + vec3(4, 5, 6)
```

### 3. Load and execute from C

```c
#include "pkg/ffi/jedil.h"

int main() {
    // Compile .jedil file to bytecode
    JedilProgram prog = jedil_compile_file("vec_add.jedil");

    // Execute and get result
    double x, y, z;
    jedil_execute_vec3(prog, NULL, 0, &x, &y, &z);

    printf("Result: (%g, %g, %g)\n", x, y, z);  // (5, 7, 9)

    jedil_free_program(prog);
    return 0;
}
```

**Compile once:**
```bash
gcc -o myapp main.c -L. -ljedil -Wl,-rpath,.
```

**Now you can edit `vec_add.jedil` and the changes take effect immediately** - no C recompilation needed!

## Language Guide

### Types
- `float` - 64-bit floating point
- `vec3` - 3D vector (3 floats)

### Operators
```jedil
// Scalar math
a + b
a - b
a * b
a / b

// Vector math
v1 + v2         // Vector addition
v1 - v2         // Vector subtraction
v * scalar      // Scaling
cross(v1, v2)   // Cross product
dot(v1, v2)     // Dot product (returns float)
mag(v)          // Magnitude (returns float)
```

### Built-in Functions
- `vec3(x, y, z)` - Create 3D vector
- `cross(v1, v2)` - Cross product
- `dot(v1, v2)` - Dot product
- `mag(v)` - Vector magnitude

### Examples

**Cross Product:**
```jedil
// i × j = k
return cross(vec3(1, 0, 0), vec3(0, 1, 0))
// Result: vec3(0, 0, 1)
```

**Vector Magnitude:**
```jedil
// Pythagorean triple: 3-4-5
return mag(vec3(3, 4, 0))
// Result: 5.0
```

### Program Lifecycle

```c
// Compile from .jedil file
JedilProgram jedil_compile_file(const char* filepath);

// Compile from source string
JedilProgram jedil_compile_source(const char* source);

// Free program
void jedil_free_program(JedilProgram program);
```

### Execution

```c
// Execute and return vec3
JedilError jedil_execute_vec3(
    JedilProgram program,
    const void* input_data,
    size_t input_len,
    double* result_x,
    double* result_y,
    double* result_z
);

// Execute and return float
JedilError jedil_execute_float(
    JedilProgram program,
    const void* input_data,
    size_t input_len,
    double* result
);
```

### Error Handling

```c
const char* jedil_get_last_error();
```

Error codes:
- `JEDIL_OK = 0`
- `JEDIL_ERROR_EXECUTION_FAILED = 3`
- `JEDIL_ERROR_STACK_UNDERFLOW = 4`
- `JEDIL_ERROR_TYPE_MISMATCH = 5`

## Project Structure

```
jedil/
├── pkg/
│   ├── bytecode/       # OpCode definitions
│   ├── vm/             # Virtual machine + SIMD ops
│   ├── types/          # Vec3, Vec3Batch (SoA layout)
│   ├── compiler/       # Lexer, parser, code generator
│   └── ffi/            # C bindings (CGO)
├── cmd/
│   ├── test/           # VM tests
│   └── test_compiler/  # Compiler tests
├── examples/
│   ├── vec_add.jedil         # Simple example
│   ├── moid_helper.jedil     # MOID calculation
│   └── c_hotreload/          # Hot-reload demo
└── README.md
```

## Use Cases

### Astrodynamics
- **Hermite splines** - Tune interpolation parameters without C recompilation
- **Orbit propagation** - Experiment with integration methods
- **Kepler solvers** - Hot-swap between different numerical methods

### General Scientific Computing
- Physics simulations with runtime-tunable parameters
- Machine learning inference with hot-swappable models
- Signal processing with live filter updates
- Robotics path planning with dynamic obstacle avoidance

## Running Examples

```bash
# Test the compiler
go run cmd/test_compiler/main.go

# Test VM directly
go run cmd/test/main.go

# Test C FFI
gcc -o test examples/c_test/test.c -L. -ljedil -Wl,-rpath,.
./test

# Test hot-reload capability
gcc -o demo examples/c_hotreload/demo_hotreload.c -L. -ljedil -Wl,-rpath,.
./demo
```

## Benchmarks

```bash
# VM benchmarks (SIMD comparison)
cd pkg/vm
go test -bench=.

# Results:
# BenchmarkVMScalar-8    7,569,165 ns/op  (157.8 ns per 4-vec-add)
# BenchmarkVMBatch-8    11,191,740 ns/op  (106.7 ns per 4-vec-add)
# Speedup: 1.48×
```

## Roadmap

- [x] Milestone 1: Minimal stack VM
- [x] Milestone 2: Vector types (Vec3)
- [x] Milestone 3: SIMD batch processing
- [x] Milestone 4: C FFI
- [x] Milestone 5: Text-to-bytecode compiler (**HOT-RELOAD!**)
- [ ] Milestone 6: Variables and functions
- [ ] Milestone 7: JIT compilation (3-4× speedup)
- [ ] Milestone 8: Python bindings
- [ ] Milestone 9: Hermite spline reference implementation

## Architecture

### Two-Tier Design

1. **Compile-time (Go):** `.jedil` → Bytecode
2. **Runtime (C):** Execute bytecode

This avoids Go's garbage collection in performance-critical paths while leveraging Go's SIMD auto-vectorization.

### VM Registry Pattern

JEDIL uses a VM registry to safely pass Go objects to C without violating CGO pointer rules:

```go
var vmRegistry = make(map[int]*vm.VM)

func registerVM(v *vm.VM) unsafe.Pointer {
    handle := nextVMHandle++
    vmRegistry[handle] = v
    return unsafe.Pointer(uintptr(handle))
}
```

C receives an integer handle, not a Go pointer.

## Contributing

This is an experimental project for astrodynamics research. Contributions welcome!

Areas of interest:
- JIT compilation backend
- More vector operations (quaternions, matrices)
- Optimization passes
- Python/Rust bindings

## License

MIT License - see LICENSE file

## Author

Built for the Mission Planning and astrodynamics community.

---

**"Change your algorithms at the speed of thought, not the speed of compilation."**
