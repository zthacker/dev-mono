# Dev Mono Repo - Root Makefile
# Orchestrates builds across all language ecosystems

.PHONY: help build clean clean-all rebuild test run format
.DEFAULT_GOAL := help

# Directory paths
MONO_ROOT := $(shell pwd)
C_DIR := $(MONO_ROOT)/c
# GO_DIR := $(MONO_ROOT)/go
# IAC_DIR := $(MONO_ROOT)/iac

# Color output (disable with NO_COLOR=1)
ifndef NO_COLOR
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m
else
CYAN :=
GREEN :=
YELLOW :=
NC :=
endif

help:
	@echo "$(CYAN)Dev Mono Repo - Build System$(NC)"
	@echo "=============================="
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  make build            - Build all projects"
	@echo "  make clean            - Clean build artifacts (specific BUILD_TYPE)"
	@echo "  make clean-all        - Clean all build artifacts"
	@echo "  make rebuild          - Clean and rebuild all"
	@echo "  make test             - Run all tests"
	@echo "  make format           - Format all code"
	@echo "  make help             - Show this help"
	@echo ""
	@echo "$(GREEN)Language-specific targets:$(NC)"
	@echo "  make build-c          - Build C++ projects"
	@echo "  make clean-c          - Clean C++ artifacts"
	@echo "  make test-c           - Run C++ tests"
	@echo "  make run-c-event-app  - Run C++ event-app"
	@echo "  make format-c         - Format C++ code"
	@echo ""
	@echo "$(GREEN)Variables:$(NC)"
	@echo "  BUILD_TYPE=<type>     - Set build type (Debug, Release, etc.)"
	@echo "  ARGS='<args>'         - Pass arguments to executables"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make build                           - Build all (Debug)"
	@echo "  make build BUILD_TYPE=Release        - Build all (Release)"
	@echo "  make run-c-event-app                 - Run event-app (Debug)"
	@echo "  make run-c-event-app BUILD_TYPE=Release ARGS='--help'"
	@echo ""
	@echo "$(YELLOW)For more details, see language-specific Makefiles:$(NC)"
	@echo "  C++: $(C_DIR)/Makefile"

build: build-c
	@echo "$(GREEN)All builds complete!$(NC)"

build-c:
	@echo "$(CYAN)Building C++ projects...$(NC)"
ifdef BUILD_TYPE
	@cd $(C_DIR) && $(MAKE) build BUILD_TYPE=$(BUILD_TYPE)
else
	@cd $(C_DIR) && $(MAKE) build
endif

clean: clean-c

clean-c:
	@echo "$(CYAN)Cleaning C++ artifacts...$(NC)"
ifdef BUILD_TYPE
	@cd $(C_DIR) && $(MAKE) clean BUILD_TYPE=$(BUILD_TYPE)
else
	@cd $(C_DIR) && $(MAKE) clean
endif

clean-all: clean-all-c
	@echo "$(GREEN)All clean complete!$(NC)"

clean-all-c:
	@echo "$(CYAN)Cleaning all C++ artifacts...$(NC)"
	@cd $(C_DIR) && $(MAKE) clean-all

rebuild: clean build

test: test-c

test-c:
	@echo "$(CYAN)Running C++ tests...$(NC)"
ifdef BUILD_TYPE
	@cd $(C_DIR) && $(MAKE) test BUILD_TYPE=$(BUILD_TYPE)
else
	@cd $(C_DIR) && $(MAKE) test
endif

format: format-c

format-c:
	@echo "$(CYAN)Formatting C++ code...$(NC)"
	@cd $(C_DIR) && $(MAKE) format

# C++ specific run targets
run-c-event-app:
ifdef BUILD_TYPE
	@cd $(C_DIR) && $(MAKE) run-event-app BUILD_TYPE=$(BUILD_TYPE) ARGS=$(ARGS)
else
	@cd $(C_DIR) && $(MAKE) run-event-app ARGS=$(ARGS)
endif
