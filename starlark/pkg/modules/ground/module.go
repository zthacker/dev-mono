package ground

import (
	"context"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"example_automation/pkg/backend"
	"example_automation/pkg/registry"
)

// GroundModule implements ground station operations
type GroundModule struct {
	gsService backend.GroundStationService
}

// NewGroundModule creates a new ground module
func NewGroundModule(gsService backend.GroundStationService) *GroundModule {
	return &GroundModule{
		gsService: gsService,
	}
}

// Metadata returns the module metadata
func (m *GroundModule) Metadata() registry.ModuleMetadata {
	return registry.ModuleMetadata{
		Name:        "ground",
		Description: "Ground station operations and scheduling",
		Version:     "1.0.0",
		Category:    "ground",
		Author:      "Mission Control Team",
		Functions: []registry.FunctionMetadata{
			{
				Name:        "track",
				Description: "Initiate tracking of a satellite",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "sat_id",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Satellite identifier",
						Required:    true,
					},
					{
						Name:        "duration",
						Type:        registry.TypeMetadata{Name: "int"},
						Description: "Tracking duration in seconds",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "bool"},
				Examples: []registry.Example{
					{
						Title: "Track satellite for 300 seconds",
						Code:  `success = ground.track("SAT-001", 300)`,
					},
				},
			},
			{
				Name:        "schedule",
				Description: "Schedule a satellite pass",
				Parameters: []registry.ParameterMetadata{
					{
						Name:        "pass_id",
						Type:        registry.TypeMetadata{Name: "string"},
						Description: "Pass identifier",
						Required:    true,
					},
					{
						Name:        "config",
						Type:        registry.TypeMetadata{Name: "dict", IsDict: true},
						Description: "Pass configuration parameters",
						Required:    true,
					},
				},
				ReturnType: registry.TypeMetadata{Name: "string"},
				Examples: []registry.Example{
					{
						Title: "Schedule a pass",
						Code:  `schedule_id = ground.schedule("PASS-123", {"start_time": "2025-12-07T10:00:00Z", "duration": 600})`,
					},
				},
			},
		},
		Examples: []registry.Example{
			{
				Title:       "Track and schedule",
				Description: "Track a satellite and schedule its next pass",
				Code: `# Start tracking
success = ground.track("SAT-001", 300)
if success:
    print("Tracking initiated")
    # Schedule next pass
    schedule_id = ground.schedule("PASS-456", {"duration": 600})
    print("Scheduled: " + schedule_id)`,
			},
		},
	}
}

// Build constructs the Starlark module value
func (m *GroundModule) Build() starlark.Value {
	members := starlark.StringDict{
		"track":    starlark.NewBuiltin("track", m.track),
		"schedule": starlark.NewBuiltin("schedule", m.schedule),
	}
	return starlarkstruct.FromStringDict(starlark.String("ground"), members)
}

// track is the Starlark function implementation
func (m *GroundModule) track(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var satID string
	var duration int

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 2, &satID, &duration); err != nil {
		return nil, err
	}

	ctx := context.Background()
	success, err := m.gsService.Track(ctx, satID, duration)
	if err != nil {
		return nil, err
	}

	return starlark.Bool(success), nil
}

// schedule is the Starlark function implementation
func (m *GroundModule) schedule(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var passID string
	var config *starlark.Dict

	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 2, &passID, &config); err != nil {
		return nil, err
	}

	// Convert Starlark dict to Go map
	configMap := make(map[string]interface{})
	for _, key := range config.Keys() {
		value, _, _ := config.Get(key)
		configMap[key.String()] = convertStarlarkValue(value)
	}

	ctx := context.Background()
	scheduleID, err := m.gsService.Schedule(ctx, passID, configMap)
	if err != nil {
		return nil, err
	}

	return starlark.String(scheduleID), nil
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
var _ registry.Module = (*GroundModule)(nil)
