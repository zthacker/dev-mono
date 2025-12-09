#include <stdio.h>
#include <unistd.h>
#include "../../pkg/ffi/jedil.h"

int main() {
    printf("=== JEDIL Hot-Reload Demo ===\n\n");
    printf("This demo shows you can change .jedil files without recompiling C!\n\n");

    // Load the initial algorithm
    printf("Loading examples/vec_add.jedil...\n");
    JedilProgram prog = jedil_compile_file("examples/vec_add.jedil");

    if (prog == NULL) {
        printf("Failed to load: %s\n", jedil_get_last_error());
        return 1;
    }

    // Execute it
    double x, y, z;
    jedil_execute_vec3(prog, NULL, 0, &x, &y, &z);
    printf("Result: (%g, %g, %g)\n\n", x, y, z);

    // Free and reload
    jedil_free_program(prog);

    printf("Now loading examples/moid_helper.jedil (different algorithm)...\n");
    prog = jedil_compile_file("examples/moid_helper.jedil");

    if (prog == NULL) {
        printf("Failed to load: %s\n", jedil_get_last_error());
        return 1;
    }

    // Execute the new algorithm
    jedil_execute_vec3(prog, NULL, 0, &x, &y, &z);
    printf("Relative velocity: (%g, %g, %g) km/s\n\n", x, y, z);

    printf("Success! Algorithm changed without recompiling C code!\n");
    printf("\nKey Point: You can edit .jedil files and reload them at runtime.\n");
    printf("Perfect for:\n");
    printf("  - Experimenting with different MOID algorithms\n");
    printf("  - Tuning Hermite spline interpolation\n");
    printf("  - Testing new astrodynamics calculations\n");
    printf("  - Hot-fixing bugs in production\n");

    jedil_free_program(prog);
    return 0;
}
