package world

import (
	"context"
	"log"
	"net"
	"sync/atomic"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"realm-server/internal/entity"
	realm_net "realm-server/internal/net"
	"realm-server/internal/zone"
)

// Server is the main world server orchestrator.
//
// A single world server instance handles:
// - One or more zones (or zone shards)
// - All players currently in those zones
// - NPCs, mobs, objects in those zones
//
// For horizontal scaling, run multiple world servers:
// - Each handles different zones, OR
// - Each handles a shard of the same zone (dynamic sharding)
type Server struct {
	cfg Config

	// NATS for inter-service communication
	nc *nats.Conn
	js jetstream.JetStream

	// TCP listener for client connections
	listener net.Listener

	// Core managers
	entityMgr *entity.Manager
	zoneMgr   *zone.Manager
	//sessionMgr *SessionManager

	// Game loop
	tickLoop *TickLoop

	// Metrics
	activeConnections int64

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// Config holds world server configuration.
type Config struct {
	// Network
	BindAddr string // e.g., ":8085"
	NATSAddr string // e.g., "nats://localhost:4222"

	// This server's identity
	ServerID uint32   // Unique ID for this server instance
	ZoneIDs  []uint32 // Which zones this server handles

	// Game loop
	TickRate int // Default 20 (50ms per tick)

	// Limits
	MaxPlayers      int
	MaxEntitiesZone int
}

func NewServer(cfg Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	listener, err := net.Listen("tcp", s.cfg.BindAddr)
	if err != nil {
		return err
	}
	s.listener = listener
	log.Printf("GameServer listening on: %s", s.cfg.BindAddr)

	// tick rate setup
	tickRate := s.cfg.TickRate
	if tickRate == 0 {
		tickRate = 20 // default to 20 tickets/sec
	}
	s.tickLoop = NewTickLoop(tickRate)
	s.tickLoop.SetOnTick(s.onTick)

	// start game loop in background
	stopChan := make(chan struct{})
	go s.tickLoop.Run(stopChan)

	go s.acceptLoop()
	// block until done
	<-s.ctx.Done()

	// cleanup
	close(stopChan)
	s.listener.Close()

	return nil

}

func (s *Server) Stop() error {
	// 	  Graceful shutdown:
	//   1. Stop accepting new connections
	//   2. Notify all players of shutdown
	//   3. Save all player data
	//   4. Persist entity states to NATS KV
	//   5. Close NATS connection

	return nil

}

func (s *Server) onTick(tick uint64, dt float64) {
	// Phase 1: Input (TODO)
	// s.processNATSMessages()

	// Phase 2: Simulate (TODO)
	// s.entityMgr.UpdateAll(dt)
	// s.zoneMgr.UpdateAll(dt)

	// Phase 3: Synchronize (TODO)
	// s.broadcastUpdates()

	// Periodic tasks
	if tick%200 == 0 { // Every 10 seconds
		log.Printf("Tick %d | Connections: %d | dt=%.3fs",
			tick, atomic.LoadInt64(&s.activeConnections), dt)
	}
}

// =============================================================================
// CONNECTION HANDLING
// =============================================================================

func (s *Server) acceptLoop() {
	for {
		playerConn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Printf("error connecting player: %s", err)
				continue
			}
		}

		playerSession := realm_net.NewSession(playerConn)

		go s.handleSession(playerSession)
	}
}

func (s *Server) handleSession(session *realm_net.Session) {
	atomic.AddInt64(&s.activeConnections, 1)
	defer atomic.AddInt64(&s.activeConnections, -1)

	session.Run()
}

// =============================================================================
// NATS INTEGRATION
// =============================================================================

// TODO: Implement cross-shard communication:
//
// func (s *Server) setupNATSSubscriptions() error
//   Subscribe to:
//   - "world.{serverID}.>" for direct messages
//   - "zone.{zoneID}.>" for each zone we handle
//   - "broadcast.>" for realm-wide broadcasts
//
// func (s *Server) handleCrossShardMessage(msg *nats.Msg)
//   Message types:
//   - Player transfer between shards
//   - Spell effects crossing shard boundaries
//   - Chat/guild messages
//   - Auction/mail notifications
//
// func (s *Server) publishToZone(zoneID uint32, subject string, data []byte)
//
// func (s *Server) transferPlayer(player *entity.Player, targetServerID uint32)
//   - Serialize player state
//   - Publish to target server
//   - Remove from local world
