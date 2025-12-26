package world

import (
	"context"
	"net"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"realm-server/internal/entity"
	"realm-server/internal/zone"
)

// Server is the main world server orchestrator.
// This is analogous to your Service in the satellite system.
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
	// sessionMgr *SessionManager

	// Game loop
	tickRate   int     // Ticks per second (typically 20)
	tickDelta  float64 // Seconds per tick

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
	ServerID   uint32   // Unique ID for this server instance
	ZoneIDs    []uint32 // Which zones this server handles

	// Game loop
	TickRate int // Default 20 (50ms per tick)

	// Limits
	MaxPlayers      int
	MaxEntitiesZone int
}

// TODO: Implement Server:
//
// func NewServer(cfg Config) *Server
//
// func (s *Server) Run(ctx context.Context) error
//   Main startup sequence:
//   1. Connect to NATS
//   2. Initialize entity manager
//   3. Initialize zone manager, load zone data
//   4. Start TCP listener
//   5. Start game loop
//   6. Start accepting connections
//   7. Block until context cancelled
//
// func (s *Server) Stop() error
//   Graceful shutdown:
//   1. Stop accepting new connections
//   2. Notify all players of shutdown
//   3. Save all player data
//   4. Persist entity states to NATS KV
//   5. Close NATS connection

// =============================================================================
// CONNECTION HANDLING
// =============================================================================

// TODO: Implement connection acceptance:
//
// func (s *Server) acceptLoop()
//   - Accept TCP connections
//   - Create Session for each
//   - Spawn goroutine to handle session
//   - Respect MaxPlayers limit
//
// func (s *Server) handleSession(session *net.Session)
//   - Run session packet loop
//   - Cleanup on disconnect

// =============================================================================
// GAME LOOP
// =============================================================================

// TODO: Implement the game tick loop:
//
// func (s *Server) gameLoop()
//   - Use time.Ticker for consistent tick rate
//   - Each tick:
//     1. Process pending messages from NATS (cross-shard)
//     2. Update all entities (entityMgr.UpdateAll)
//     3. Process combat timers
//     4. Update zone states (respawns, events)
//     5. Broadcast state updates to clients
//
// Tick timing matters:
// - 20 ticks/sec = 50ms per tick (WoW standard)
// - If tick takes >50ms, you're falling behind
// - Log warnings when ticks take too long

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
