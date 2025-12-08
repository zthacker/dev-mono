#include <stdio.h>
#include "../../pkg/ffi/jedil.h"

int main() {
    printf("=== JEDIL C FFI Test ===\n\n");

    // Test 1: Direct vector addition (no VM)
    printf("Test 1: Direct Vec3 Addition (no VM overhead)\n");
    double ax = 1.0, ay = 2.0, az = 3.0;
    double bx = 4.0, by = 5.0, bz = 6.0;
    double rx, ry, rz;

    jedil_vec3_add(ax, ay, az, bx, by, bz, &rx, &ry, &rz);
    printf("  (%g, %g, %g) + (%g, %g, %g) = (%g, %g, %g)\n",
           ax, ay, az, bx, by, bz, rx, ry, rz);

    if (rx == 5.0 && ry == 7.0 && rz == 9.0) {
        printf("  ✓ PASS\n\n");
    } else {
        printf("  ✗ FAIL\n\n");
        return 1;
    }

    // Test 2: Batch SIMD addition
    printf("Test 2: Batch SIMD Addition (4 vectors at once)\n");
    double a_xs[4] = {1.0, 2.0, 3.0, 4.0};
    double a_ys[4] = {1.0, 2.0, 3.0, 4.0};
    double a_zs[4] = {1.0, 2.0, 3.0, 4.0};

    double b_xs[4] = {10.0, 20.0, 30.0, 40.0};
    double b_ys[4] = {10.0, 20.0, 30.0, 40.0};
    double b_zs[4] = {10.0, 20.0, 30.0, 40.0};

    double result_xs[4], result_ys[4], result_zs[4];

    jedil_batch_add(a_xs, a_ys, a_zs, b_xs, b_ys, b_zs, result_xs, result_ys, result_zs);

    printf("  Batch A[0] + B[0] = (%g, %g, %g)\n", result_xs[0], result_ys[0], result_zs[0]);
    printf("  Batch A[3] + B[3] = (%g, %g, %g)\n", result_xs[3], result_ys[3], result_zs[3]);

    if (result_xs[0] == 11.0 && result_xs[3] == 44.0) {
        printf("PASS - SIMD works\n\n");
    } else {
        printf("FAIL womp\n\n");
        return 1;
    }

    printf("=== All Tests Passed! ===\n");

    return 0;
}
