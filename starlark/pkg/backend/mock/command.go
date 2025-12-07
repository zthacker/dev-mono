package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"example_automation/pkg/backend"
)

// CommandRecord tracks commands sent
type CommandRecord struct {
	Timestamp time.Time
	Command   string
	Params    map[string]interface{}
}

// MockCommandService simulates command execution
type MockCommandService struct {
	commandLog []CommandRecord
	mu         sync.Mutex
}

// NewMockCommandService creates a new mock command service
func NewMockCommandService() *MockCommandService {
	return &MockCommandService{
		commandLog: make([]CommandRecord, 0),
	}
}

// SendCommand logs and simulates command execution
func (m *MockCommandService) SendCommand(ctx context.Context, command string, params map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	record := CommandRecord{
		Timestamp: time.Now(),
		Command:   command,
		Params:    params,
	}
	m.commandLog = append(m.commandLog, record)

	fmt.Printf("[MOCK CMD] Executing %s with params: %v\n", command, params)

	// Simulate command execution delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// GetCommandStatus returns mock status
func (m *MockCommandService) GetCommandStatus(ctx context.Context, commandID string) (backend.CommandStatus, error) {
	return backend.CommandStatus{
		ID:        commandID,
		Status:    "completed",
		Message:   "Command executed successfully",
		Timestamp: time.Now(),
	}, nil
}

// GetCommandLog returns all commands sent (for testing/debugging)
func (m *MockCommandService) GetCommandLog() []CommandRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	logCopy := make([]CommandRecord, len(m.commandLog))
	copy(logCopy, m.commandLog)
	return logCopy
}

// Ensure it implements the interface
var _ backend.CommandService = (*MockCommandService)(nil)
