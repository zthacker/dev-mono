package mock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"example_automation/pkg/backend"
)

// MockDataProcessingService simulates data processing operations
type MockDataProcessingService struct{}

// NewMockDataProcessingService creates a new mock data processing service
func NewMockDataProcessingService() *MockDataProcessingService {
	return &MockDataProcessingService{}
}

// Process simulates data processing
func (m *MockDataProcessingService) Process(ctx context.Context, rawData []byte) ([]byte, error) {
	fmt.Printf("[MOCK DATA] Processing %d bytes of data\n", len(rawData))
	time.Sleep(50 * time.Millisecond)

	// Simple mock processing: convert to uppercase
	processed := []byte(strings.ToUpper(string(rawData)))
	return processed, nil
}

// Validate simulates schema validation
func (m *MockDataProcessingService) Validate(ctx context.Context, schema string, data interface{}) (bool, error) {
	fmt.Printf("[MOCK DATA] Validating data against schema: %s\n", schema)
	time.Sleep(50 * time.Millisecond)

	// Mock validation always passes
	return true, nil
}

// Ensure it implements the interface
var _ backend.DataProcessingService = (*MockDataProcessingService)(nil)
