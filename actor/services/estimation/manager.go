package estimation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// ActorManager manages the lifecycle of individual EKF Actors.
// Responsibility: Manage the map of satellites and ensure thread safety.
// It does NOT know about NATS sharding - it just processes whatever satellite ID it receives.
type ActorManager struct {
	actors    map[int]*EKFActor
	actorLock sync.RWMutex

	// NATS KV stores for persistence
	kvConfig jetstream.KeyValue
	kvState  jetstream.KeyValue

	// Algorithm factory (allows swapping EKF implementations)
	algorithmFactory func() StateEstimator
}

// NewActorManager creates a new ActorManager.
func NewActorManager() *ActorManager {
	return &ActorManager{
		actors: make(map[int]*EKFActor),
		// Default to ExtendedKalmanFilter
		algorithmFactory: func() StateEstimator {
			return NewExtendedKalmanFilter()
		},
	}
}

// SetAlgorithmFactory allows injecting a different state estimator algorithm.
func (am *ActorManager) SetAlgorithmFactory(factory func() StateEstimator) {
	am.algorithmFactory = factory
}

// Start initializes the ActorManager and connects to NATS KV stores.
func (am *ActorManager) Start(nc *nats.Conn, js jetstream.JetStream) error {
	var err error

	// Create or get KV bucket for configuration
	am.kvConfig, err = js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "ses_config",
		Description: "State Estimation Service Configuration",
	})
	if err != nil {
		return fmt.Errorf("failed to create/get config KV: %w", err)
	}

	// Create or get KV bucket for state persistence
	am.kvState, err = js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "ses_state",
		Description: "Satellite State Persistence",
	})
	if err != nil {
		return fmt.Errorf("failed to create/get state KV: %w", err)
	}

	log.Println("ActorManager initialized with NATS KV stores")
	return nil
}

// ProcessMeasurement is the main entry point for processing a measurement.
// This is called by worker threads, so it must be thread-safe.
func (am *ActorManager) ProcessMeasurement(satID int, measurement *GPSMeasurement) (*EKFState, error) {
	// 1. Get or initialize the actor (handles map-level locking)
	actor := am.getOrInitActor(satID)

	// 2. Lock the specific actor (ensures only one thread processes this satellite at a time)
	actor.mu.Lock()
	defer actor.mu.Unlock()

	// 3. Calculate time delta for prediction step
	dt := time.Since(actor.state.LastUpdate).Seconds()

	// 4. Predict step (propagate state forward)
	if dt > 0 {
		if err := actor.algorithm.Predict(actor.state, dt); err != nil {
			return nil, fmt.Errorf("predict step failed: %w", err)
		}
	}

	// 5. Update step (incorporate measurement)
	if err := actor.algorithm.Update(actor.state, measurement); err != nil {
		return nil, fmt.Errorf("update step failed: %w", err)
	}

	// 6. Update timestamp
	actor.state.LastUpdate = measurement.Timestamp

	// Return a copy of the state (caller shouldn't hold reference to internal state)
	return actor.state, nil
}

// getOrInitActor uses double-checked locking to efficiently get or create actors.
func (am *ActorManager) getOrInitActor(satID int) *EKFActor {
	// Fast path: Read lock
	am.actorLock.RLock()
	actor, exists := am.actors[satID]
	am.actorLock.RUnlock()
	if exists {
		return actor
	}

	// Slow path: Write lock
	am.actorLock.Lock()
	defer am.actorLock.Unlock()

	// Double-check: another goroutine might have created it
	if actor, exists = am.actors[satID]; exists {
		return actor
	}

	// Create new actor
	newActor := NewEKFActor(satID, am.algorithmFactory())

	// Try to load previous state from KV
	if am.kvState != nil {
		if err := am.loadActorState(newActor); err != nil {
			log.Printf("Could not load state for satellite %d: %v (starting fresh)", satID, err)
		}
	}

	am.actors[satID] = newActor
	log.Printf("Initialized new actor for satellite %d", satID)
	return newActor
}

// loadActorState attempts to restore actor state from NATS KV.
func (am *ActorManager) loadActorState(actor *EKFActor) error {
	key := fmt.Sprintf("sat_%d", actor.satID)
	entry, err := am.kvState.Get(context.Background(), key)
	if err != nil {
		return err
	}

	var savedState struct {
		State      []float64   `json:"state"`
		Covariance [][]float64 `json:"covariance"`
		LastUpdate int64       `json:"last_update"`
	}

	if err := json.Unmarshal(entry.Value(), &savedState); err != nil {
		return err
	}

	// Restore state vector
	for i, v := range savedState.State {
		actor.state.State.SetVec(i, v)
	}

	// Restore covariance matrix
	rows, cols := actor.state.Covariance.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			actor.state.Covariance.Set(i, j, savedState.Covariance[i][j])
		}
	}

	actor.state.LastUpdate = time.Unix(0, savedState.LastUpdate)
	log.Printf("Loaded state for satellite %d from KV (age: %v)", actor.satID, time.Since(actor.state.LastUpdate))
	return nil
}

// Shutdown persists all actor states to NATS KV.
func (am *ActorManager) Shutdown(ctx context.Context) error {
	am.actorLock.RLock()
	defer am.actorLock.RUnlock()

	log.Printf("Persisting state for %d satellites...", len(am.actors))

	for satID, actor := range am.actors {
		actor.mu.Lock()

		// Serialize state
		rows, cols := actor.state.Covariance.Dims()
		covData := make([][]float64, rows)
		for i := 0; i < rows; i++ {
			covData[i] = make([]float64, cols)
			for j := 0; j < cols; j++ {
				covData[i][j] = actor.state.Covariance.At(i, j)
			}
		}

		savedState := struct {
			State      []float64   `json:"state"`
			Covariance [][]float64 `json:"covariance"`
			LastUpdate int64       `json:"last_update"`
		}{
			State:      actor.state.State.RawVector().Data,
			Covariance: covData,
			LastUpdate: actor.state.LastUpdate.UnixNano(),
		}

		data, err := json.Marshal(savedState)
		actor.mu.Unlock()

		if err != nil {
			log.Printf("Failed to serialize state for satellite %d: %v", satID, err)
			continue
		}

		// Save to KV
		key := fmt.Sprintf("sat_%d", satID)
		if _, err := am.kvState.Put(ctx, key, data); err != nil {
			log.Printf("Failed to persist state for satellite %d: %v", satID, err)
		}
	}

	log.Println("State persistence complete")
	return nil
}

// GetStats returns statistics about the actor manager.
func (am *ActorManager) GetStats() map[string]interface{} {
	am.actorLock.RLock()
	defer am.actorLock.RUnlock()

	return map[string]interface{}{
		"active_satellites": len(am.actors),
	}
}
