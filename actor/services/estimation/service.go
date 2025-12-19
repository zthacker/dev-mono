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
	"gonum.org/v1/gonum/mat"
)

// Service is the top-level orchestrator.
// Responsibility: Listen to NATS and manage worker threads.
// It does NOT know about individual satellites - that's the ActorManager's job.
type Service struct {
	cfg ServiceConfig
	nc  *nats.Conn
	js  jetstream.JetStream
	sub jetstream.ConsumeContext

	// The brain: delegates all satellite-specific logic
	actorManager *ActorManager

	// Worker pool
	jobQueue chan *nats.Msg
	workerWg sync.WaitGroup

	// Lifecycle
	shutdownOnce sync.Once
	ctx          context.Context
	cancel       context.CancelFunc

	stats *Stats
}

// NewService creates a new State Estimation Service instance.
func NewService(cfg ServiceConfig) *Service {
	return &Service{
		cfg:          cfg,
		actorManager: NewActorManager(),
		stats:        &Stats{},
	}
}

// Run starts the service and blocks until the context is cancelled.
func (s *Service) Run(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.jobQueue = make(chan *nats.Msg, s.cfg.WorkerCount*2)

	// 1. Connect to NATS
	log.Printf("Connecting to NATS at %s", s.cfg.NATSAddr)
	nc, err := nats.Connect(s.cfg.NATSAddr,
		nats.ReconnectWait(s.cfg.RetryDelay),
		nats.MaxReconnects(s.cfg.MaxRetries),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	s.nc = nc

	// 2. Initialize JetStream
	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}
	s.js = js

	// 3. Start ActorManager (loads state from NATS KV if available)
	log.Println("Starting ActorManager...")
	if err := s.actorManager.Start(s.nc, s.js); err != nil {
		return fmt.Errorf("failed to start actor manager: %w", err)
	}

	// 4. Start worker pool
	log.Printf("Starting %d workers", s.cfg.WorkerCount)
	for i := 0; i < s.cfg.WorkerCount; i++ {
		s.workerWg.Add(1)
		go s.worker(s.ctx, i)
	}

	// 5. Subscribe to JetStream (sharding happens here via InputSubject)
	log.Printf("Subscribing to %s (consumer: %s)", s.cfg.InputSubject, s.cfg.ConsumerName)
	cons, err := js.CreateOrUpdateConsumer(ctx, "TELEMETRY", jetstream.ConsumerConfig{
		Durable:       s.cfg.ConsumerName,
		FilterSubject: s.cfg.InputSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	consumeCtx, err := cons.Consume(s.messageHandler)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}
	s.sub = consumeCtx

	log.Printf("Service started successfully. Listening on %s", s.cfg.InputSubject)

	// Block until context is cancelled
	<-s.ctx.Done()
	s.Stop()
	return nil
}

// Stop gracefully shuts down the service.
func (s *Service) Stop() {
	s.shutdownOnce.Do(func() {
		log.Println("Initiating graceful shutdown...")
		s.cancel()

		// Stop accepting new messages
		if s.sub != nil {
			s.sub.Stop()
		}

		// Drain job queue and wait for workers to finish
		close(s.jobQueue)
		s.workerWg.Wait()

		// Persist all actor states to NATS KV
		log.Println("Persisting actor states...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.actorManager.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during actor manager shutdown: %v", err)
		}

		// Close NATS connection
		if s.nc != nil {
			s.nc.Close()
		}

		log.Println("Shutdown complete")
	})
}

// Stats returns the current service statistics.
func (s *Service) Stats() *Stats {
	return s.stats
}

// messageHandler receives messages from NATS and queues them for workers.
func (s *Service) messageHandler(msg jetstream.Msg) {
	s.stats.received.Add(1)

	// Convert to legacy *nats.Msg for queue compatibility
	legacyMsg := &nats.Msg{
		Subject: msg.Subject(),
		Data:    msg.Data(),
	}

	select {
	case s.jobQueue <- legacyMsg:
		msg.Ack()
	case <-s.ctx.Done():
		msg.Nak()
	default:
		// Backpressure: tell NATS to retry later
		msg.NakWithDelay(5 * time.Second)
	}
}

// worker processes messages from the job queue.
func (s *Service) worker(ctx context.Context, workerID int) {
	defer s.workerWg.Done()
	defer func() {
		if r := recover(); r != nil {
			s.stats.panics.Add(1)
			log.Printf("Worker %d panic: %v", workerID, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-s.jobQueue:
			if !ok {
				return
			}
			s.processMessage(ctx, msg)
		}
	}
}

// processMessage handles a single measurement message.
func (s *Service) processMessage(ctx context.Context, msg *nats.Msg) {
	// 1. Unmarshal the measurement packet
	var packet MeasurementPacket
	if err := json.Unmarshal(msg.Data, &packet); err != nil {
		log.Printf("Failed to unmarshal packet: %v", err)
		s.stats.failed.Add(1)
		return
	}

	// 2. Convert to GPSMeasurement
	measurement := &GPSMeasurement{
		Timestamp:  time.Unix(0, packet.Timestamp),
		Position:   mat.NewVecDense(3, []float64{packet.PosX, packet.PosY, packet.PosZ}),
		Covariance: mat.NewDense(3, 3, []float64{packet.CovXX, 0, 0, 0, packet.CovYY, 0, 0, 0, packet.CovZZ}),
	}

	// 3. Delegate to ActorManager (handles locking and EKF math)
	satID := int(packet.SatelliteID)
	newState, err := s.actorManager.ProcessMeasurement(satID, measurement)
	if err != nil {
		log.Printf("Error processing satellite %d: %v", satID, err)
		s.stats.failed.Add(1)
		return
	}

	// 4. Publish updated state (optional - configure output subject as needed)
	outputSubject := fmt.Sprintf("telemetry.state.%d", satID)
	outputData, _ := json.Marshal(map[string]interface{}{
		"satellite_id": satID,
		"timestamp":    newState.LastUpdate.Unix(),
		"state":        newState.State.RawVector().Data,
		"covariance":   newState.Covariance.RawMatrix().Data,
	})
	s.nc.Publish(outputSubject, outputData)

	s.stats.processed.Add(1)
}
