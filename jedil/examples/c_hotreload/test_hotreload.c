#include <stdio.h>
#include "../../pkg/ffi/jedil.h"

int main() {
    printf("=== JEDIL Hot-Reload Test ===\n\n");

    // Test 1: Compile from source string
    printf("Test 1: Compile from source string\n");
    JedilProgram prog1 = jedil_compile_source("return vec3(1, 2, 3) + vec3(4, 5, 6)");

    if (prog1 == NULL) {
        printf("  ❌ FAIL: Compilation failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  ✓ Compilation succeeded\n");

    // Execute it
    double rx, ry, rz;
    int result = jedil_execute_vec3(prog1, NULL, 0, &rx, &ry, &rz);

    if (result != JEDIL_OK) {
        printf("  ❌ FAIL: Execution failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  Result: (%g, %g, %g)\n", rx, ry, rz);

    if (rx == 5.0 && ry == 7.0 && rz == 9.0) {
        printf("  ✓ PASS\n\n");
    } else {
        printf("  ❌ FAIL: incorrect result\n\n");
        return 1;
    }

    jedil_free_program(prog1);

    // Test 2: Compile from file
    printf("Test 2: Compile from .jedil file\n");
    JedilProgram prog2 = jedil_compile_file("examples/vec_add.jedil");

    if (prog2 == NULL) {
        printf("  ❌ FAIL: File compilation failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  ✓ File compilation succeeded\n");

    result = jedil_execute_vec3(prog2, NULL, 0, &rx, &ry, &rz);

    if (result != JEDIL_OK) {
        printf("  ❌ FAIL: Execution failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  Result: (%g, %g, %g)\n", rx, ry, rz);

    if (rx == 5.0 && ry == 7.0 && rz == 9.0) {
        printf("  ✓ PASS\n\n");
    } else {
        printf("  ❌ FAIL: incorrect result\n\n");
        return 1;
    }

    jedil_free_program(prog2);

    // Test 3: Cross product from source
    printf("Test 3: Cross product (i × j = k)\n");
    JedilProgram prog3 = jedil_compile_source(
        "return cross(vec3(1, 0, 0), vec3(0, 1, 0))"
    );

    if (prog3 == NULL) {
        printf("  ❌ FAIL: Compilation failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    result = jedil_execute_vec3(prog3, NULL, 0, &rx, &ry, &rz);

    if (result != JEDIL_OK) {
        printf("  ❌ FAIL: Execution failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  Result: (%g, %g, %g)\n", rx, ry, rz);

    if (rx == 0.0 && ry == 0.0 && rz == 1.0) {
        printf("  ✓ PASS\n\n");
    } else {
        printf("  ❌ FAIL: incorrect result\n\n");
        return 1;
    }

    jedil_free_program(prog3);

    // Test 4: Dot product (returns scalar)
    printf("Test 4: Dot product\n");
    JedilProgram prog4 = jedil_compile_source(
        "return dot(vec3(1, 2, 3), vec3(4, 5, 6))"
    );

    if (prog4 == NULL) {
        printf("  ❌ FAIL: Compilation failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    double scalar_result;
    result = jedil_execute_float(prog4, NULL, 0, &scalar_result);

    if (result != JEDIL_OK) {
        printf("  ❌ FAIL: Execution failed: %s\n\n", jedil_get_last_error());
        return 1;
    }

    printf("  Result: %g\n", scalar_result);

    // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
    if (scalar_result == 32.0) {
        printf("  ✓ PASS\n\n");
    } else {
        printf("  ❌ FAIL: expected 32.0\n\n");
        return 1;
    }

    jedil_free_program(prog4);

    printf("=== All Hot-Reload Tests Passed! ===\n");
    printf("\n🎉 SUCCESS: You can now change .jedil files without recompiling C code!\n");

    return 0;
}
