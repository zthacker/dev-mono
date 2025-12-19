package estimation

import (
	"sync"
)

// EKFActor represents a single satellite's state estimation actor.
// Responsibility: Execute the EKF math for ONE specific satellite.
// Each actor is single-threaded (protected by its own mutex) to ensure mathematical consistency.
type EKFActor struct {
	satID int

	// Internal mutex: Ensures only ONE thread updates THIS satellite at a time.
	// This makes the actor thread-safe regardless of NATS parallelism.
	mu sync.Mutex

	// The state estimation algorithm (EKF, UKF, etc.)
	algorithm StateEstimator

	// Current state (position, velocity, covariance)
	state *EKFState
}

// NewEKFActor creates a new actor for a satellite.
// The state is initialized with zeros - you can override this by loading from NATS KV.
func NewEKFActor(satID int, algorithm StateEstimator) *EKFActor {
	return &EKFActor{
		satID:     satID,
		algorithm: algorithm,
		state:     NewEKFState(6), // 6D state: [x, y, z, vx, vy, vz]
	}
}

// GetSatID returns the satellite ID.
func (a *EKFActor) GetSatID() int {
	return a.satID
}

// GetState returns a reference to the current state (caller must hold lock).
func (a *EKFActor) GetState() *EKFState {
	return a.state
}
