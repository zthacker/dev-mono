package satellite

import (
	"context"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"example_automation/pkg/backend"
	"example_automation/pkg/registry"
)

// SatelliteModule implements the satellite operations module
type SatelliteModule struct {
	telemetryService backend.TelemetryService
	commandService   backend.CommandService
}

// NewSatelliteModule creates a new satellite module with backend services
func NewSatelliteModule(
	telemetryService backend.TelemetryService,
	commandService backend.CommandService,
) *SatelliteModule {
	return &SatelliteModule{
		telemetryService: telemetryService,
		commandService:   commandService,
	}
}

// Metadata returns the module metadata
func (m *SatelliteModule) Metadata() registry.ModuleMetadata {
	return registry.ModuleMetadata{
		Name:        "satellite",
		Description: "Satellite telemetry and command operations",
		Version:     "1.0.0",
		Category:    "satellite",
		Author:      "Mission Control Team",
		Functions: []registry.FunctionMetadata{
			{
				Name:        "getTLM",
				Description: "Retrieve telemetry value for a given mnemonic",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "mnemonic",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Telemetry point identifier (e.g., 'battery_level', 'temperature')",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "float"},
				Examples: []registry.Example{
					{
						Title:  "Get battery level",
						Code:   `battery = satellite.getTLM("battery_level")`,
						Output: "85.5",
					},
					{
						Title:  "Check temperature",
						Code:   `temp = satellite.getTLM("temperature")`,
						Output: "23.2",
					},
				},
			},
			{
				Name:        "sendCMD",
				Description: "Send a command to the satellite with parameters",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "command",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Command name to execute",
						Required:    true,
					},
					{
						Name: "params",
						Type: registry.TypeMetadata{
							Name:   "dict",
							IsDict: true,
						},
						Description: "Command parameters as key-value pairs",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "None"},
				Examples: []registry.Example{
					{
						Title: "Set transmitter configuration",
						Code:  `satellite.sendCMD("SET_TRANSMITTER", {"power": "high", "freq": 2200})`,
					},
					{
						Title: "Power on subsystem",
						Code:  `satellite.sendCMD("POWER_ON", {"subsystem": "payload"})`,
					},
				},
			},
		},
		Examples: []registry.Example{
			{
				Title:       "Pre-pass safety check",
				Description: "Check battery level before starting transmitter",
				Code: `battery = satellite.getTLM("battery_level")
if battery < 20:
    print("ABORT: Low battery")
else:
    satellite.sendCMD("SET_TRANSMITTER", {"power": "high"})`,
			},
		},
	}
}

// Build constructs the Starlark module value
func (m *SatelliteModule) Build() starlark.Value {
	members := starlark.StringDict{
		"getTLM":  starlark.NewBuiltin("getTLM", m.getTLM),
		"sendCMD": starlark.NewBuiltin("sendCMD", m.sendCMD),
	}
	return starlarkstruct.FromStringDict(starlark.String("satellite"), members)
}

// getTLM is the Starlark function implementation
func (m *SatelliteModule) getTLM(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var mnemonic string
	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 1, &mnemonic); err != nil {
		return nil, err
	}

	// Call backend service
	ctx := context.Background()
	value, err := m.telemetryService.GetTelemetry(ctx, mnemonic)
	if err != nil {
		return nil, err
	}

	return starlark.Float(value), nil
}

// sendCMD is the Starlark function implementation
func (m *SatelliteModule) sendCMD(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var command string
	var params *starlark.Dict

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 2, &command, &params); err != nil {
		return nil, err
	}

	// Convert Starlark dict to Go map
	paramMap := make(map[string]interface{})
	for _, key := range params.Keys() {
		value, _, _ := params.Get(key)
		paramMap[key.String()] = convertStarlarkValue(value)
	}

	// Call backend service
	ctx := context.Background()
	if err := m.commandService.SendCommand(ctx, command, paramMap); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

// convertStarlarkValue converts a Starlark value to a Go interface{}
func convertStarlarkValue(v starlark.Value) interface{} {
	switch v := v.(type) {
	case starlark.String:
		return v.GoString()
	case starlark.Int:
		i, _ := v.Int64()
		return i
	case starlark.Float:
		return float64(v)
	case starlark.Bool:
		return bool(v)
	default:
		return v.String()
	}
}

// Ensure it implements the Module interface
var _ registry.Module = (*SatelliteModule)(nil)
