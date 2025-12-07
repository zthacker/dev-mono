package registry

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// InjectIntrospection adds help() and dir() functions to globals
func InjectIntrospection(registry *Registry, globals starlark.StringDict) {
	globals["help"] = starlark.NewBuiltin("help", makeHelpFunc(registry))
	globals["dir"] = starlark.NewBuiltin("dir", makeDirFunc(registry))
}

// makeHelpFunc creates the help() builtin
func makeHelpFunc(registry *Registry) func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// help() with no args - show all modules
		if len(args) == 0 {
			return starlark.String(formatAllModulesHelp(registry)), nil
		}

		// help(module) or help(module.function)
		var obj starlark.Value
		if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &obj); err != nil {
			return nil, err
		}

		// Check if it's a module (Struct)
		if module, ok := obj.(*starlarkstruct.Struct); ok {
			moduleName := module.Constructor().String()
			// Remove quotes from module name
			moduleName = strings.Trim(moduleName, "\"")
			return starlark.String(formatModuleHelp(registry, moduleName)), nil
		}

		// Check if it's a function (Builtin)
		if fn, ok := obj.(*starlark.Builtin); ok {
			return starlark.String(formatFunctionHelp(registry, fn.Name())), nil
		}

		return starlark.String("No help available for this object"), nil
	}
}

// makeDirFunc creates the dir() builtin
func makeDirFunc(registry *Registry) func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// dir() with no args - show all available modules
		if len(args) == 0 {
			modules := registry.All()
			items := make([]starlark.Value, len(modules))
			for i, mod := range modules {
				items[i] = starlark.String(mod.Metadata().Name)
			}
			return starlark.NewList(items), nil
		}

		// dir(module) - show all functions in module
		var obj starlark.Value
		if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &obj); err != nil {
			return nil, err
		}

		if module, ok := obj.(*starlarkstruct.Struct); ok {
			var items []starlark.Value
			for _, name := range module.AttrNames() {
				items = append(items, starlark.String(name))
			}
			return starlark.NewList(items), nil
		}

		return starlark.NewList(nil), nil
	}
}

// formatAllModulesHelp formats help for all modules
func formatAllModulesHelp(registry *Registry) string {
	var sb strings.Builder
	sb.WriteString("Available Modules:\n\n")

	modules := registry.All()
	for _, mod := range modules {
		meta := mod.Metadata()
		sb.WriteString(fmt.Sprintf("  %s - %s (v%s)\n", meta.Name, meta.Description, meta.Version))
	}

	sb.WriteString("\nUse help(module_name) for detailed information.\n")
	return sb.String()
}

// formatModuleHelp formats help for a specific module
func formatModuleHelp(registry *Registry, moduleName string) string {
	module, err := registry.Get(moduleName)
	if err != nil {
		return fmt.Sprintf("Module %s not found", moduleName)
	}

	meta := module.Metadata()
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Module: %s (v%s)\n", meta.Name, meta.Version))
	sb.WriteString(fmt.Sprintf("%s\n\n", meta.Description))
	sb.WriteString("Functions:\n\n")

	for _, fn := range meta.Functions {
		sb.WriteString(fmt.Sprintf("  %s.%s(", meta.Name, fn.Name))

		// Parameter list
		params := make([]string, len(fn.Parameters))
		for i, p := range fn.Parameters {
			if p.Required {
				params[i] = p.Name
			} else {
				params[i] = fmt.Sprintf("%s=%s", p.Name, p.Default)
			}
		}
		sb.WriteString(strings.Join(params, ", "))
		sb.WriteString(fmt.Sprintf(") -> %s\n", fn.ReturnType.Name))
		sb.WriteString(fmt.Sprintf("    %s\n\n", fn.Description))
	}

	if len(meta.Examples) > 0 {
		sb.WriteString("Examples:\n\n")
		for _, ex := range meta.Examples {
			sb.WriteString(fmt.Sprintf("  %s:\n", ex.Title))
			if ex.Description != "" {
				sb.WriteString(fmt.Sprintf("    %s\n", ex.Description))
			}
			// Indent code examples
			codeLines := strings.Split(ex.Code, "\n")
			for _, line := range codeLines {
				sb.WriteString(fmt.Sprintf("    %s\n", line))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// formatFunctionHelp formats help for a specific function
func formatFunctionHelp(registry *Registry, functionName string) string {
	// Search all modules for this function
	modules := registry.All()
	for _, mod := range modules {
		meta := mod.Metadata()
		for _, fn := range meta.Functions {
			if fn.Name == functionName {
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("%s.%s()\n\n", meta.Name, fn.Name))
				sb.WriteString(fmt.Sprintf("%s\n\n", fn.Description))

				sb.WriteString("Parameters:\n")
				for _, p := range fn.Parameters {
					req := ""
					if p.Required {
						req = " (required)"
					} else {
						req = fmt.Sprintf(" (optional, default: %s)", p.Default)
					}
					sb.WriteString(fmt.Sprintf("  %s (%s)%s\n    %s\n", p.Name, p.Type.Name, req, p.Description))
				}

				sb.WriteString(fmt.Sprintf("\nReturns: %s\n", fn.ReturnType.Name))

				if len(fn.Examples) > 0 {
					sb.WriteString("\nExamples:\n")
					for _, ex := range fn.Examples {
						sb.WriteString(fmt.Sprintf("  %s\n", ex.Code))
					}
				}

				return sb.String()
			}
		}
	}

	return fmt.Sprintf("Function %s not found", functionName)
}
