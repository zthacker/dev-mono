package system

import (
	"fmt"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"example_automation/pkg/registry"
)

// SystemModule implements system utility functions
type SystemModule struct{}

// NewSystemModule creates a new system module
func NewSystemModule() *SystemModule {
	return &SystemModule{}
}

// Metadata returns the module metadata
func (m *SystemModule) Metadata() registry.ModuleMetadata {
	return registry.ModuleMetadata{
		Name:        "system",
		Description: "System utility functions for script control",
		Version:     "1.0.0",
		Category:    "system",
		Author:      "Mission Control Team",
		Functions: []registry.FunctionMetadata{
			{
				Name:        "wait",
				Description: "Pause script execution for a specified number of seconds",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "seconds",
						Type:        registry.TypeMetadata{Name: "int"},
						Description: "Number of seconds to wait",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "None"},
				Examples: []registry.Example{
					{
						Title: "Wait for 5 seconds",
						Code:  `system.wait(5)`,
					},
					{
						Title: "Wait between operations",
						Code:  `satellite.sendCMD("POWER_ON", {})\nsystem.wait(2)\nprint("Command complete")`,
					},
				},
			},
		},
		Examples: []registry.Example{
			{
				Title:       "Basic wait usage",
				Description: "Pause execution between commands",
				Code: `print("Starting operation")
system.wait(3)
print("Operation complete after 3 seconds")`,
			},
		},
	}
}

// Build constructs the Starlark module value
func (m *SystemModule) Build() starlark.Value {
	members := starlark.StringDict{
		"wait": starlark.NewBuiltin("wait", m.wait),
	}
	return starlarkstruct.FromStringDict(starlark.String("system"), members)
}

// wait is the Starlark function implementation
func (m *SystemModule) wait(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var seconds int
	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 1, &seconds); err != nil {
		return nil, err
	}

	fmt.Printf("[Go-Engine] Sleeping for %d seconds...\n", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)

	return starlark.None, nil
}

// Ensure it implements the Module interface
var _ registry.Module = (*SystemModule)(nil)
