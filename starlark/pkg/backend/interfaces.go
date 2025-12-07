package backend

import (
	"context"
	"time"
)

// TelemetryService handles telemetry retrieval
type TelemetryService interface {
	// GetTelemetry retrieves a telemetry value by mnemonic
	GetTelemetry(ctx context.Context, mnemonic string) (float64, error)

	// GetTelemetryBatch retrieves multiple telemetry values
	GetTelemetryBatch(ctx context.Context, mnemonics []string) (map[string]float64, error)
}

// CommandService handles command sending
type CommandService interface {
	// SendCommand sends a command with parameters
	SendCommand(ctx context.Context, command string, params map[string]interface{}) error

	// GetCommandStatus checks the status of a previously sent command
	GetCommandStatus(ctx context.Context, commandID string) (CommandStatus, error)
}

// CommandStatus represents the execution status of a command
type CommandStatus struct {
	ID        string
	Status    string // "pending", "executing", "completed", "failed"
	Message   string
	Timestamp time.Time
}

// GroundStationService handles ground station operations
type GroundStationService interface {
	// Track initiates tracking of a satellite
	Track(ctx context.Context, satID string, duration int) (bool, error)

	// Schedule creates a pass schedule
	Schedule(ctx context.Context, passID string, config map[string]interface{}) (string, error)
}

// DataProcessingService handles data operations
type DataProcessingService interface {
	// Process processes raw data
	Process(ctx context.Context, rawData []byte) ([]byte, error)

	// Validate validates data against a schema
	Validate(ctx context.Context, schema string, data interface{}) (bool, error)
}
