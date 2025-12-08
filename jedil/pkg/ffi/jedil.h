#ifndef JEDIL_H
#define JEDIL_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// ============================================================================
// TYPES
// ============================================================================

// Opaque handle to a JEDIL program
typedef void* JedilProgram;

// 3D Vector (matches Go's types.Vec3)
typedef struct {
    double x;
    double y;
    double z;
} JedilVec3;

// Batch of 4 vectors (matches Go's types.Vec3Batch)
typedef struct {
    double xs[4];
    double ys[4];
    double zs[4];
} JedilVec3Batch;

// Error code
typedef enum {
    JEDIL_OK = 0,
    JEDIL_ERROR_NULL_POINTER = 1,
    JEDIL_ERROR_INVALID_BYTECODE = 2,
    JEDIL_ERROR_EXECUTION_FAILED = 3,
    JEDIL_ERROR_STACK_UNDERFLOW = 4,
    JEDIL_ERROR_TYPE_MISMATCH = 5,
} JedilError;

// ============================================================================
// PROGRAM LIFECYCLE
// ============================================================================

// Create a program from bytecode array
// Returns: Opaque program handle (NULL on error)
JedilProgram jedil_create_program(const uint8_t* bytecode, size_t len);

// Free a program
void jedil_free_program(JedilProgram program);

// ============================================================================
// EXECUTION - VECTOR OPERATIONS
// ============================================================================

// Execute program that returns a Vec3
// Inputs: arbitrary data buffer (interpreted by bytecode)
// Output: result vector
JedilError jedil_execute_vec3(
    JedilProgram program,
    const void* input_data,
    size_t input_len,
    double* result_x,
    double* result_y,
    double* result_z
);

// Execute program that returns a float
JedilError jedil_execute_float(
    JedilProgram program,
    const void* input_data,
    size_t input_len,
    double* result
);

// Execute program that returns a batch
JedilError jedil_execute_batch(
    JedilProgram program,
    const void* input_data,
    size_t input_len,
    JedilVec3Batch* result
);

// ============================================================================
// CONVENIENCE - PRE-BUILT OPERATIONS
// ============================================================================

// Vector addition (no bytecode needed, direct call)
void jedil_vec3_add(
    double ax, double ay, double az,
    double bx, double by, double bz,
    double* result_x,
    double* result_y,
    double* result_z
);
// Batch vector addition (SIMD)
void jedil_batch_add(
    const double* a_xs, const double* a_ys, const double* a_zs,
    const double* b_xs, const double* b_ys, const double* b_zs,
    double* result_xs, double* result_ys, double* result_zs
);
// ============================================================================
// ERROR HANDLING
// ============================================================================

// Get last error message (thread-safe)
const char* jedil_get_last_error();

#ifdef __cplusplus
}
#endif

#endif // JEDIL_H