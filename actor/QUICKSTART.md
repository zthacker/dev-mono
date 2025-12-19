# Quick Start Guide

Get the State Estimation Service up and running in 5 minutes.

## Prerequisites

Install NATS Server with JetStream:
```bash
# macOS
brew install nats-server

# Or download from https://nats.io/download/
```

## Step 1: Start NATS Server

```bash
nats-server -js
```

You should see:
```
[1] 2024/12/19 ... [INF] Starting nats-server
[1] 2024/12/19 ... [INF] JetStream enabled
```

## Step 2: Run the Service

In a new terminal:
```bash
cd /Users/zachthacker/zach_repos/dev-mono/actor
go run main.go
```

You should see:
```
Starting State Estimation Service...
Connecting to NATS at nats://localhost:4222
Stream ready: TELEMETRY
NATS infrastructure ready
Starting ActorManager...
ActorManager initialized with NATS KV stores
Starting 8 workers
Subscribing to telemetry.identified.partition.1.> (consumer: ses-worker-partition-1)
Service started successfully. Listening on telemetry.identified.partition.1.>
```

## Step 3: Publish Test Data

In another terminal:
```bash
cd /Users/zachthacker/zach_repos/dev-mono/actor
go run examples/publisher/main.go
```

You should see measurements being published and the service processing them.

## Step 4: Monitor State Output

Subscribe to the output stream to see estimated states:
```bash
nats sub "telemetry.state.*"
```

You'll see JSON messages with updated satellite states:
```json
{
  "satellite_id": 1,
  "timestamp": 1734639240,
  "state": [7000123.45, 456789.12, 234567.89, 7456.78, 1234.56, 890.12],
  "covariance": [...]
}
```

## Step 5: Run Multiple Shards (Optional)

To run multiple service instances handling different satellite ranges:

**Terminal 1 - Shard 1 (even satellites):**
```bash
# Modify main.go line 33 to:
cfg.InputSubject = "telemetry.identified.partition.1.>"
go run main.go
```

**Terminal 2 - Shard 2 (odd satellites):**
```bash
# Modify main.go line 33 to:
cfg.InputSubject = "telemetry.identified.partition.2.>"
go run main.go
```

## Verify Everything Works

Check NATS streams:
```bash
nats stream ls
nats stream info TELEMETRY
```

Check KV buckets:
```bash
nats kv ls
nats kv get ses_state sat_1
```

Check consumer status:
```bash
nats consumer ls TELEMETRY
```

## What's Happening?

1. **Service** listens to `telemetry.identified.partition.1.>` for measurements
2. **Worker threads** (8 concurrent) process incoming messages
3. **ActorManager** routes each satellite to its own EKFActor
4. **EKFActor** runs Predict + Update steps
5. **Updated state** is published to `telemetry.state.{sat_id}`
6. **On shutdown** (Ctrl+C), states are saved to NATS KV

## Next Steps

Now you're ready to implement the actual EKF algorithm!

### Replace the placeholder orbital dynamics in `services/estimation/algos.go`:

1. **Predict step (line 40)**: Replace constant-velocity model with proper orbital propagation
   - Use Keplerian propagation or SGP4
   - Add J2 perturbations
   - Consider atmospheric drag and solar radiation pressure

2. **Tune noise matrices**:
   - `Q` (process noise): Depends on propagation uncertainty
   - `R` (measurement noise): From GPS receiver specs

3. **Add state initialization**:
   - Use TLE data or other orbit determination
   - Initialize covariance based on uncertainty

### Example improvements:

```go
// In Predict():
// 1. Use RK4 integration with orbital dynamics
// 2. Compute state transition matrix Φ via STM propagation
// 3. Include perturbations (J2, drag, SRP)

// In Update():
// 1. Handle different measurement types (range, range-rate, angles)
// 2. Add measurement validation (chi-square test)
// 3. Consider extended Kalman filter linearization for nonlinear measurements
```

## Troubleshooting

**Service won't start:**
- Check NATS is running: `nats-server --version`
- Check port 4222 is available: `lsof -i :4222`

**No measurements processed:**
- Verify stream exists: `nats stream info TELEMETRY`
- Check subject routing matches publisher and consumer
- Review service logs for errors

**States not persisting:**
- Ensure graceful shutdown (Ctrl+C, not kill -9)
- Check KV bucket: `nats kv ls`
- Review shutdown logs

## Development Tips

1. **Enable debug logging** in main.go:
   ```go
   log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
   ```

2. **Monitor performance**:
   ```bash
   # Watch processing rate
   watch -n 1 'nats stream info TELEMETRY'
   ```

3. **Test with realistic data**: Replace random GPS data with TLE-derived positions

4. **Add unit tests**: Test EKF math independently of NATS infrastructure
