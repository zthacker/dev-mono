package data

import (
	"context"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"example_automation/pkg/backend"
	"example_automation/pkg/registry"
)

// DataModule implements data processing operations
type DataModule struct {
	dataService backend.DataProcessingService
}

// NewDataModule creates a new data module
func NewDataModule(dataService backend.DataProcessingService) *DataModule {
	return &DataModule{
		dataService: dataService,
	}
}

// Metadata returns the module metadata
func (m *DataModule) Metadata() registry.ModuleMetadata {
	return registry.ModuleMetadata{
		Name:        "data",
		Description: "Data processing and validation operations",
		Version:     "1.0.0",
		Category:    "data",
		Author:      "Mission Control Team",
		Functions: []registry.FunctionMetadata{
			{
				Name:        "process",
				Description: "Process raw data and return processed result",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "raw_data",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Raw data to process",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "string"},
				Examples: []registry.Example{
					{
						Title: "Process telemetry data",
						Code:  `processed = data.process("raw telemetry data")`,
					},
				},
			},
			{
				Name:        "validate",
				Description: "Validate data against a schema",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "schema",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Schema name or definition",
						Required:    true,
					},
					{
						Name:        "data_value",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Data to validate",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "bool"},
				Examples: []registry.Example{
					{
						Title: "Validate telemetry format",
						Code:  `is_valid = data.validate("telemetry_schema", telemetry_data)`,
					},
				},
			},
		},
		Examples: []registry.Example{
			{
				Title:       "Process and validate data",
				Description: "Process raw data and validate the result",
				Code: `raw = "sensor data"
processed = data.process(raw)
is_valid = data.validate("sensor_schema", processed)
if is_valid:
    print("Data validated successfully")`,
			},
		},
	}
}

// Build constructs the Starlark module value
func (m *DataModule) Build() starlark.Value {
	members := starlark.StringDict{
		"process":  starlark.NewBuiltin("process", m.process),
		"validate": starlark.NewBuiltin("validate", m.validate),
	}
	return starlarkstruct.FromStringDict(starlark.String("data"), members)
}

// process is the Starlark function implementation
func (m *DataModule) process(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var rawData string

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 1, &rawData); err != nil {
		return nil, err
	}

	ctx := context.Background()
	processed, err := m.dataService.Process(ctx, []byte(rawData))
	if err != nil {
		return nil, err
	}

	return starlark.String(string(processed)), nil
}

// validate is the Starlark function implementation
func (m *DataModule) validate(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var schema string
	var dataValue string

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 2, &schema, &dataValue); err != nil {
		return nil, err
	}

	ctx := context.Background()
	isValid, err := m.dataService.Validate(ctx, schema, dataValue)
	if err != nil {
		return nil, err
	}

	return starlark.Bool(isValid), nil
}

// Ensure it implements the Module interface
var _ registry.Module = (*DataModule)(nil)
