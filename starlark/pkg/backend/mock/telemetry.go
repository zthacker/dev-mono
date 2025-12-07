package mock

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"example_automation/pkg/backend"
)

// MockTelemetryService provides mock telemetry data
type MockTelemetryService struct {
	values map[string]float64
	mu     sync.RWMutex
}

// NewMockTelemetryService creates a new mock telemetry service
func NewMockTelemetryService() *MockTelemetryService {
	return &MockTelemetryService{
		values: map[string]float64{
			"battery_level":   85.5,
			"temperature":     23.2,
			"solar_voltage":   28.4,
			"signal_strength": 78.3,
		},
	}
}

// GetTelemetry retrieves a telemetry value with simulated variance
func (m *MockTelemetryService) GetTelemetry(ctx context.Context, mnemonic string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Simulate some variance
	if value, exists := m.values[mnemonic]; exists {
		// Add random variance ±5%
		variance := (rand.Float64() - 0.5) * 0.1 * value
		return value + variance, nil
	}

	return 0, fmt.Errorf("telemetry mnemonic %s not found", mnemonic)
}

// GetTelemetryBatch retrieves multiple telemetry values
func (m *MockTelemetryService) GetTelemetryBatch(ctx context.Context, mnemonics []string) (map[string]float64, error) {
	result := make(map[string]float64)
	for _, mnemonic := range mnemonics {
		value, err := m.GetTelemetry(ctx, mnemonic)
		if err != nil {
			return nil, err
		}
		result[mnemonic] = value
	}
	return result, nil
}

// Ensure it implements the interface
var _ backend.TelemetryService = (*MockTelemetryService)(nil)
