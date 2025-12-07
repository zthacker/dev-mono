package mock

import (
	"context"
	"fmt"
	"time"

	"example_automation/pkg/backend"
)

// MockGroundStationService simulates ground station operations
type MockGroundStationService struct{}

// NewMockGroundStationService creates a new mock ground station service
func NewMockGroundStationService() *MockGroundStationService {
	return &MockGroundStationService{}
}

// Track simulates satellite tracking
func (m *MockGroundStationService) Track(ctx context.Context, satID string, duration int) (bool, error) {
	fmt.Printf("[MOCK GS] Tracking satellite %s for %d seconds\n", satID, duration)
	time.Sleep(50 * time.Millisecond)
	return true, nil
}

// Schedule simulates pass scheduling
func (m *MockGroundStationService) Schedule(ctx context.Context, passID string, config map[string]interface{}) (string, error) {
	fmt.Printf("[MOCK GS] Scheduling pass %s with config: %v\n", passID, config)
	time.Sleep(50 * time.Millisecond)
	scheduleID := fmt.Sprintf("SCH-%s-%d", passID, time.Now().Unix())
	return scheduleID, nil
}

// Ensure it implements the interface
var _ backend.GroundStationService = (*MockGroundStationService)(nil)
