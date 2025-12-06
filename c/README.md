# C/C++ Build System - Quick Reference

## Build Commands

### Basic Build Operations

| Command | Description |
|---------|-------------|
| `make build` | Build all C++ projects in Debug mode (default) |
| `make build BUILD_TYPE=Release` | Build all C++ projects in Release mode |
| `make clean` | Remove build artifacts for current BUILD_TYPE |
| `make clean-all` | Remove all build directories (Debug, Release, etc.) |
| `make rebuild` | Clean and rebuild (same as `make clean && make build`) |
| `make configure` | Force CMake reconfiguration for current BUILD_TYPE |

### Build Types

Set the build type using `BUILD_TYPE=<type>`:

| Build Type | Description |
|------------|-------------|
| `Debug` | Debug build with symbols and assertions (default) |
| `Release` | Optimized release build with no debug symbols |
| `RelWithDebInfo` | Release build with debug information |
| `MinSizeRel` | Release build optimized for minimum size |

**Example:**
```bash
make build BUILD_TYPE=Release
```

### Running Applications

| Command | Description |
|---------|-------------|
| `make run-event-app` | Run event-app in Debug mode |
| `make run-event-app BUILD_TYPE=Release` | Run event-app in Release mode |
| `make run-event-app ARGS='--help'` | Run event-app with arguments |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all tests for current BUILD_TYPE |
| `make test BUILD_TYPE=Release` | Run tests for Release build |

### Code Formatting

| Command | Description |
|---------|-------------|
| `make format` | Format all C++ source files with clang-format |
| `make format-check` | Check formatting without modifying files |

### Help

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands and options |

## Build Directory Structure

The build system creates separate directories for each build type:

```
c/
├── build-Debug/          # Debug build artifacts
├── build-Release/        # Release build artifacts
├── build-RelWithDebInfo/ # RelWithDebInfo build artifacts (if used)
└── build-MinSizeRel/     # MinSizeRel build artifacts (if used)
```

This allows you to switch between build types without reconfiguring or rebuilding.

## Common Workflows

### Debug Workflow
```bash
# Build and run in Debug mode
make build
make run-event-app

# Make changes to code...

# Rebuild and test
make rebuild
make test
```

### Release Workflow
```bash
# Build optimized release version
make build BUILD_TYPE=Release

# Run release build
make run-event-app BUILD_TYPE=Release

# Clean only release build
make clean BUILD_TYPE=Release
```

### Clean Start
```bash
# Remove all build artifacts and rebuild
make clean-all
make build
```

### Format Code Before Commit
```bash
# Format all source files
make format

# Or check formatting first
make format-check
```

## Examples

```bash
# Build Debug version with parallel compilation
make build

# Build Release version
make build BUILD_TYPE=Release

# Run Debug version
make run-event-app

# Run Release version with arguments
make run-event-app BUILD_TYPE=Release ARGS='--verbose'

# Clean Debug build
make clean

# Clean Release build
make clean BUILD_TYPE=Release

# Clean everything
make clean-all

# Format all code
make format

# Run tests in Debug mode
make test
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BUILD_TYPE` | `Debug` | Build configuration type |
| `ARGS` | (empty) | Arguments to pass to executables |
| `NO_COLOR` | (unset) | Set to `1` to disable color output |

**Example:**
```bash
# Disable color output
NO_COLOR=1 make build
```

## Troubleshooting

### CMake version error
If you see an error about CMake version, ensure you have CMake 3.20 or higher:
```bash
cmake --version
```

### Clean build issues
If you encounter build errors, try a clean rebuild:
```bash
make clean-all
make build
```

### Multiple build types
You can maintain multiple build types simultaneously:
```bash
make build BUILD_TYPE=Debug
make build BUILD_TYPE=Release
make build BUILD_TYPE=RelWithDebInfo

# All three build directories now exist
ls -l build-*
```
