# C/C++ Build System
# Orchestrates CMake-based builds for all C++ projects

.PHONY: help build clean clean-all rebuild configure test run run-event-app format format-check
.DEFAULT_GOAL := help

# Dynamic path resolution - NO hardcoded paths
C_ROOT := $(shell pwd)
BUILD_TYPE ?= Debug
BUILD_DIR := $(C_ROOT)/build-$(BUILD_TYPE)

# CMake configuration
CMAKE := cmake
CMAKE_GENERATOR := Unix Makefiles

# Detect number of CPU cores for parallel builds
NPROC := $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

# Find all source files for formatting (exclude build and CMake-generated files)
SOURCES := $(shell find $(C_ROOT) -type f \( -name "*.cpp" -o -name "*.h" \) 2>/dev/null | grep -v "/build" | grep -v "/CMakeFiles")

# Available applications
APPS := event-app

# Color output (disable with NO_COLOR=1)
ifndef NO_COLOR
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m
else
CYAN :=
GREEN :=
YELLOW :=
RED :=
NC :=
endif

help:
	@echo "$(CYAN)C/C++ Build System$(NC)"
	@echo "=================="
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  make build              - Build all C++ projects (default: Debug)"
	@echo "  make clean              - Remove build artifacts for BUILD_TYPE"
	@echo "  make clean-all          - Remove all build directories"
	@echo "  make rebuild            - Clean and build"
	@echo "  make configure          - Force CMake reconfiguration"
	@echo "  make test               - Run all tests"
	@echo "  make run-event-app      - Run event-app"
	@echo "  make format             - Format all source files with clang-format"
	@echo "  make format-check       - Check formatting without changes"
	@echo "  make help               - Show this help message"
	@echo ""
	@echo "$(GREEN)Build types (set with BUILD_TYPE):$(NC)"
	@echo "  Debug           - Debug build with symbols (default)"
	@echo "  Release         - Optimized release build"
	@echo "  RelWithDebInfo  - Release with debug info"
	@echo "  MinSizeRel      - Minimum size release"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make build                              - Build in Debug mode"
	@echo "  make build BUILD_TYPE=Release           - Build in Release mode"
	@echo "  make run-event-app                      - Run event-app (Debug)"
	@echo "  make run-event-app BUILD_TYPE=Release ARGS='--help'"
	@echo "  make clean                              - Clean Debug build"
	@echo "  make clean BUILD_TYPE=Release           - Clean Release build"
	@echo "  make clean-all                          - Clean all builds"
	@echo ""
	@echo "$(YELLOW)Current configuration:$(NC)"
	@echo "  BUILD_TYPE: $(BUILD_TYPE)"
	@echo "  BUILD_DIR:  $(BUILD_DIR)"
	@echo "  NPROC:      $(NPROC)"
	@echo "  Available apps: $(APPS)"

# Validate build type
validate-build-type:
	@case "$(BUILD_TYPE)" in \
		Debug|Release|RelWithDebInfo|MinSizeRel) ;; \
		*) echo "$(RED)Invalid BUILD_TYPE: $(BUILD_TYPE)$(NC)"; \
		   echo "Valid options: Debug, Release, RelWithDebInfo, MinSizeRel"; \
		   exit 1 ;; \
	esac

# Build all C++ projects
build: validate-build-type $(BUILD_DIR)
	@echo "$(CYAN)Building C++ projects in $(BUILD_TYPE) mode...$(NC)"
	@cd $(BUILD_DIR) && $(CMAKE) --build . --config $(BUILD_TYPE) -j$(NPROC)
	@echo "$(GREEN)Build complete!$(NC) Executables in: $(BUILD_DIR)"

# Configure CMake (runs if build dir doesn't exist)
$(BUILD_DIR): validate-build-type
	@echo "$(CYAN)Configuring CMake for $(BUILD_TYPE) build...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && $(CMAKE) -S $(C_ROOT) -B . \
		-DCMAKE_BUILD_TYPE=$(BUILD_TYPE) \
		-G "$(CMAKE_GENERATOR)"
	@echo "$(GREEN)CMake configuration complete!$(NC)"

# Force CMake reconfiguration
configure: validate-build-type
	@echo "$(CYAN)Force reconfiguring CMake for $(BUILD_TYPE)...$(NC)"
	@rm -rf $(BUILD_DIR)
	@$(MAKE) $(BUILD_DIR)

# Clean build artifacts for specific BUILD_TYPE
clean: validate-build-type
	@if [ -d "$(BUILD_DIR)" ]; then \
		echo "$(YELLOW)Cleaning $(BUILD_TYPE) build artifacts...$(NC)"; \
		rm -rf $(BUILD_DIR); \
		echo "$(GREEN)Clean complete!$(NC)"; \
	else \
		echo "$(YELLOW)Nothing to clean - $(BUILD_DIR) does not exist$(NC)"; \
	fi

# Clean all build directories
clean-all:
	@echo "$(YELLOW)Cleaning all build artifacts...$(NC)"
	@rm -rf $(C_ROOT)/build-*
	@echo "$(GREEN)All builds cleaned!$(NC)"

# Clean and rebuild
rebuild: clean build

# Run tests (requires CMake enable_testing() in CMakeLists.txt)
test: build
	@echo "$(CYAN)Running tests for $(BUILD_TYPE) build...$(NC)"
	@if [ -f $(BUILD_DIR)/CTestTestfile.cmake ]; then \
		cd $(BUILD_DIR) && ctest --output-on-failure -j$(NPROC); \
	else \
		echo "$(YELLOW)No tests configured. Add enable_testing() to CMakeLists.txt$(NC)"; \
	fi

# Run event-app
run-event-app: build
	@if [ -f $(BUILD_DIR)/event-app/event-app ]; then \
		echo "$(GREEN)Running event-app ($(BUILD_TYPE))...$(NC)"; \
		$(BUILD_DIR)/event-app/event-app $(ARGS); \
	else \
		echo "$(RED)Error: event-app not found in $(BUILD_DIR)/event-app/$(NC)"; \
		echo "$(YELLOW)Run 'make build' first.$(NC)"; \
		exit 1; \
	fi

# Format all source files with clang-format
format:
	@echo "$(CYAN)Formatting C++ code with clang-format...$(NC)"
	@if command -v clang-format >/dev/null 2>&1; then \
		if [ -n "$(SOURCES)" ]; then \
			for file in $(SOURCES); do \
				echo "  Formatting $$file"; \
				clang-format -i $$file; \
			done; \
			echo "$(GREEN)Formatted $(words $(SOURCES)) files$(NC)"; \
		else \
			echo "$(YELLOW)No source files found to format$(NC)"; \
		fi; \
	else \
		echo "$(RED)Error: clang-format not found.$(NC)"; \
		echo "$(YELLOW)Install with: sudo apt install clang-format$(NC)"; \
		exit 1; \
	fi

# Check formatting without modifying files
format-check:
	@echo "$(CYAN)Checking code formatting...$(NC)"
	@if command -v clang-format >/dev/null 2>&1; then \
		if [ -n "$(SOURCES)" ]; then \
			FAILED=0; \
			for file in $(SOURCES); do \
				if ! clang-format --dry-run -Werror $$file 2>/dev/null; then \
					echo "$(RED)  ✗ $$file$(NC)"; \
					FAILED=1; \
				else \
					echo "$(GREEN)  ✓ $$file$(NC)"; \
				fi; \
			done; \
			if [ $$FAILED -eq 1 ]; then \
				echo "$(RED)Formatting check failed. Run 'make format' to fix.$(NC)"; \
				exit 1; \
			else \
				echo "$(GREEN)All files properly formatted!$(NC)"; \
			fi; \
		else \
			echo "$(YELLOW)No source files found to check$(NC)"; \
		fi; \
	else \
		echo "$(RED)Error: clang-format not found.$(NC)"; \
		exit 1; \
	fi
