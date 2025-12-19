package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"actor/services/estimation"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting State Estimation Service...")

	// 1. Setup NATS infrastructure (streams)
	if err := setupNATSInfrastructure(); err != nil {
		log.Fatalf("Failed to setup NATS infrastructure: %v", err)
	}

	// 2. Configure the service
	cfg := estimation.DefaultServiceConfig()
	cfg.NATSAddr = "nats://localhost:4222"

	// IMPORTANT: This is where sharding happens!
	// Each service instance listens to a different partition
	// Example for partition 1 (satellites 0-1999):
	cfg.InputSubject = "telemetry.identified.partition.1.>"
	cfg.ConsumerName = "ses-worker-partition-1"
	cfg.WorkerCount = 8 // Adjust based on CPU cores

	// 3. Create and start the service
	svc := estimation.NewService(cfg)

	// 4. Run service in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		if err := svc.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	// 5. Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("Service error: %v", err)
	case sig := <-sigCh:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		cancel()

		// Give it time to shutdown gracefully
		time.Sleep(2 * time.Second)
	}

	log.Println("Service stopped")
}

// setupNATSInfrastructure creates the necessary JetStream streams and KV buckets.
// This is idempotent and can be run multiple times.
func setupNATSInfrastructure() error {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create TELEMETRY stream for incoming measurements
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "TELEMETRY",
		Description: "Satellite telemetry measurements",
		Subjects:    []string{"telemetry.>"},
		Retention:   jetstream.WorkQueuePolicy,
		MaxAge:      24 * time.Hour,
	})
	if err != nil {
		return err
	}
	log.Printf("Stream ready: %s", stream.CachedInfo().Config.Name)

	// Create KV buckets (done automatically by ActorManager, but we can pre-create them)
	_, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "ses_config",
		Description: "State Estimation Service Configuration",
	})
	if err != nil && err != jetstream.ErrBucketExists {
		log.Printf("Warning: Could not create ses_config KV: %v", err)
	}

	_, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "ses_state",
		Description: "Satellite State Persistence",
		TTL:         7 * 24 * time.Hour, // States expire after 7 days
	})
	if err != nil && err != jetstream.ErrBucketExists {
		log.Printf("Warning: Could not create ses_state KV: %v", err)
	}

	log.Println("NATS infrastructure ready")
	return nil
}

// Example: How to run multiple shards
//
// Shard 1: Satellites 0-1999
// cfg.InputSubject = "telemetry.identified.partition.1.>"
// cfg.ConsumerName = "ses-worker-partition-1"
//
// Shard 2: Satellites 2000-3999
// cfg.InputSubject = "telemetry.identified.partition.2.>"
// cfg.ConsumerName = "ses-worker-partition-2"
//
// Shard 3: Satellites 4000-5999
// cfg.InputSubject = "telemetry.identified.partition.3.>"
// cfg.ConsumerName = "ses-worker-partition-3"
//
// And so on...
//
// The publisher is responsible for routing messages to the correct partition
// based on satellite ID modulo number of partitions.
