package estimation

import (
	"sync/atomic"
	"time"

	"gonum.org/v1/gonum/mat"
)

// ServiceConfig holds configuration for the State Estimation Service.
// The sharding configuration happens here via InputSubject.
type ServiceConfig struct {
	// NATS connection
	NATSAddr string // e.g., "nats://localhost:4222"

	// Sharding: This determines which satellites this service instance handles
	InputSubject string // e.g., "telemetry.identified.partition.1.>" for shard 1
	ConsumerName string // e.g., "ses-worker-partition-1"

	// Worker pool configuration
	WorkerCount int // Number of concurrent message processors

	// Retry and timeout settings
	MessageTimeout time.Duration
	MaxAckPending  int
	MaxRetries     int
	RetryDelay     time.Duration
}

// DefaultServiceConfig returns sensible defaults.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		NATSAddr:       "nats://localhost:4222",
		WorkerCount:    4,
		MessageTimeout: 30 * time.Second,
		MaxAckPending:  1000,
		MaxRetries:     3,
		RetryDelay:     10 * time.Second,
	}
}

// Stats tracks service-level metrics.
type Stats struct {
	received  atomic.Uint64
	processed atomic.Uint64
	failed    atomic.Uint64
	retried   atomic.Uint64
	panics    atomic.Uint64
}

func (s *Stats) Received() uint64  { return s.received.Load() }
func (s *Stats) Processed() uint64 { return s.processed.Load() }
func (s *Stats) Failed() uint64    { return s.failed.Load() }
func (s *Stats) Retried() uint64   { return s.retried.Load() }
func (s *Stats) Panics() uint64    { return s.panics.Load() }

// StateEstimator defines the interface for state estimation algorithms (EKF, UKF, etc.).
type StateEstimator interface {
	// Predict propagates state forward in time using the dynamics model.
	Predict(state *EKFState, dt float64) error

	// Update corrects state based on new GPS measurement.
	Update(state *EKFState, measurement *GPSMeasurement) error

	// Name returns the algorithm name for logging.
	Name() string
}

// EKFState holds the current state vector and covariance for a satellite.
type EKFState struct {
	// State vector: [x, y, z, vx, vy, vz] in ECI coordinates
	State *mat.VecDense

	// Covariance matrix (6x6 for position and velocity)
	Covariance *mat.Dense

	// Last update timestamp
	LastUpdate time.Time
}

// NewEKFState creates an initialized state with given dimensions.
func NewEKFState(stateDim int) *EKFState {
	return &EKFState{
		State:      mat.NewVecDense(stateDim, nil),
		Covariance: mat.NewDense(stateDim, stateDim, nil),
		LastUpdate: time.Now(),
	}
}

// GPSMeasurement represents a GPS observation in ECI coordinates.
type GPSMeasurement struct {
	Timestamp  time.Time
	Position   *mat.VecDense // [x, y, z] in ECI
	Covariance *mat.Dense    // 3x3 measurement covariance
}

// MeasurementPacket is what we expect to unmarshal from NATS messages.
// This is a placeholder - replace with your actual protobuf definition.
type MeasurementPacket struct {
	SatelliteID int32
	Timestamp   int64
	PosX        float64
	PosY        float64
	PosZ        float64
	CovXX       float64
	CovYY       float64
	CovZZ       float64
}
