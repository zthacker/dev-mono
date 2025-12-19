# State Estimation Service - Actor Model Architecture

A distributed satellite state estimation service using the Actor Model pattern with NATS for message passing and sharding.

## Architecture Overview

The system follows a clean separation of concerns with three main layers:

### 1. Service Layer ([service.go](services/estimation/service.go))
**Responsibility**: Listen to NATS and manage worker threads.

- Connects to NATS JetStream
- Manages a pool of worker goroutines
- Handles message acknowledgment and backpressure
- Does NOT know about individual satellites
- **Sharding happens here** via `InputSubject` configuration

```go
cfg := estimation.DefaultServiceConfig()
cfg.InputSubject = "telemetry.identified.partition.1.>"  // Shard 1
cfg.ConsumerName = "ses-worker-partition-1"
cfg.WorkerCount = 8
```

### 2. ActorManager Layer ([manager.go](services/estimation/manager.go))
**Responsibility**: Manage the map of satellites and ensure thread safety.

- Maintains a map of `satellite_id -> EKFActor`
- Double-checked locking for efficient actor creation
- Handles state persistence to NATS KV
- Does NOT know about NATS sharding - just processes whatever satellite ID it receives
- Thread-safe map operations

### 3. EKFActor Layer ([actor.go](services/estimation/actor.go))
**Responsibility**: Execute the EKF math for ONE specific satellite.

- One actor per satellite
- Protected by its own mutex
- Ensures mathematical consistency (only one measurement processed at a time per satellite)
- Stores state vector, covariance, and last update time

## Data Flow

```
NATS Message → Service Worker → ActorManager.ProcessMeasurement()
                                      ↓
                                 getOrInitActor(satID)
                                      ↓
                                 Lock EKFActor
                                      ↓
                                 EKF.Predict(dt)
                                      ↓
                                 EKF.Update(measurement)
                                      ↓
                                 Unlock & Return State
                                      ↓
                                 Publish to telemetry.state.*
```

## Sharding Strategy

Sharding is configured purely at the Service level via NATS subject routing:

```
Shard 1: telemetry.identified.partition.1.> → Satellites 0-1999
Shard 2: telemetry.identified.partition.2.> → Satellites 2000-3999
Shard 3: telemetry.identified.partition.3.> → Satellites 4000-5999
```

The publisher routes messages to the correct partition:
```go
partition := (satelliteID % numPartitions) + 1
subject := fmt.Sprintf("telemetry.identified.partition.%d.%d", partition, satelliteID)
```

## State Persistence

States are automatically saved to NATS KV on graceful shutdown:
- Bucket: `ses_state`
- Key format: `sat_{satellite_id}`
- TTL: 7 days
- On startup, actors are hydrated from KV if available

## EKF Implementation

The Extended Kalman Filter is implemented in [algos.go](services/estimation/algos.go):

### State Vector (6D in ECI coordinates)
```
x = [x, y, z, vx, vy, vz]ᵀ
```

### Predict Step
```go
// Propagate state forward using dynamics model
// Currently: constant velocity (placeholder)
// TODO: Implement orbital propagation (J2, drag, SRP)
x_k+1 = F * x_k
P_k+1 = F * P_k * F^T + Q
```

### Update Step
```go
// GPS measurement model: H = [I_3x3 | 0_3x3]
// Innovation: y = z - H*x
// Kalman Gain: K = P*H^T*(H*P*H^T + R)^-1
// State update: x = x + K*y
// Covariance update: P = (I - K*H)*P
```

### Next Steps for EKF
1. Replace constant-velocity model with proper orbital dynamics (Keplerian/SGP4)
2. Implement J2 perturbations
3. Add atmospheric drag model
4. Add solar radiation pressure
5. Tune Q and R matrices based on real data
6. Consider using Runge-Kutta integration for prediction

## Running the Service

### Prerequisites
- NATS Server with JetStream enabled
- Go 1.21+

### Start NATS
```bash
nats-server -js
```

### Run a single instance
```bash
cd actor
go run main.go
```

### Run multiple shards (different terminals)
```bash
# Terminal 1 - Shard 1
go run main.go -shard=1

# Terminal 2 - Shard 2
go run main.go -shard=2

# Terminal 3 - Shard 3
go run main.go -shard=3
```

## File Structure

```
actor/
├── main.go                          # Entry point, NATS setup
└── services/
    └── estimation/
        ├── service.go               # Service layer (NATS + workers)
        ├── manager.go               # ActorManager (satellite map)
        ├── actor.go                 # EKFActor (per-satellite)
        ├── algos.go                 # EKF algorithm implementation
        └── types.go                 # Shared types and interfaces
```

## Configuration

```go
type ServiceConfig struct {
    NATSAddr       string        // NATS connection
    InputSubject   string        // Sharding: which partition to consume
    ConsumerName   string        // JetStream consumer name
    WorkerCount    int           // Concurrent processors
    MessageTimeout time.Duration // Processing timeout
    MaxRetries     int           // NATS reconnect retries
    RetryDelay     time.Duration // Reconnect delay
}
```

## Monitoring

Service exposes statistics via `service.Stats()`:
```go
stats := svc.Stats()
log.Printf("Received: %d, Processed: %d, Failed: %d",
    stats.Received(), stats.Processed(), stats.Failed())
```

ActorManager exposes satellite count:
```go
stats := actorManager.GetStats()
log.Printf("Active satellites: %d", stats["active_satellites"])
```

## Thread Safety

- **Service**: Worker pool handles concurrent message processing
- **ActorManager**: RWMutex protects satellite map
- **EKFActor**: Each actor has its own mutex for state updates

This ensures:
1. Multiple workers can process different satellites in parallel
2. Multiple measurements for the same satellite are serialized
3. No race conditions on state or covariance matrices

## Future Enhancements

1. **Metrics/Observability**: Add Prometheus metrics
2. **Health Checks**: Expose HTTP health endpoint
3. **Dynamic Sharding**: Support shard rebalancing
4. **State Queries**: HTTP/gRPC API to query satellite states
5. **Algorithm Swapping**: Hot-reload different EKF implementations
6. **Batch Processing**: Process multiple measurements per satellite in batch
7. **Prediction Service**: Allow clients to request propagated states at future times
