package entity

// Uncomment when implementing:
// import (
// 	"context"
// 	"github.com/nats-io/nats.go/jetstream"
// )

// Manager handles the lifecycle of all entities in a shard.
// Similar to your ActorManager in the satellite system - one manager per shard.
//
// Key responsibilities:
// - Create/destroy entities
// - Provide fast lookup by EntityID
// - Coordinate entity updates each tick
// - Persist entity state to NATS KV on shutdown
type Manager struct {
	// Entity storage - consider sync.Map for better concurrent access
	// or sharded maps if you have many entities
	// entities map[EntityID]Entity
	// mu       sync.RWMutex

	// NATS KV for hot state (position, health, etc.)
	// kvState jetstream.KeyValue

	// ID generator for new entities
	// nextPlayerID uint64
	// nextNPCID    uint64
}

// TODO: Implement Manager:
//
// func NewManager() *Manager
//
// func (m *Manager) Start(js jetstream.JetStream) error
//   - Create/connect to NATS KV buckets
//   - Load any cached state
//
// func (m *Manager) Shutdown(ctx context.Context) error
//   - Persist all entity states to KV
//   - Similar to your ActorManager.Shutdown
//
// func (m *Manager) Get(id EntityID) (Entity, bool)
//   - Fast lookup by ID
//
// func (m *Manager) Add(e Entity) error
//   - Add entity to world
//   - Notify zone manager
//
// func (m *Manager) Remove(id EntityID) error
//   - Remove from world
//   - Cleanup references
//
// func (m *Manager) UpdateAll(dt float64)
//   - Called each server tick
//   - Iterate all entities and call Update(dt)
//   - Consider parallelizing with worker pool

// =============================================================================
// ENTITY ID GENERATION
// =============================================================================

// TODO: Implement ID generation:
//
// func (m *Manager) NextPlayerID() EntityID
//   - Atomically increment counter
//   - Pack with EntityTypePlayer
//
// func (m *Manager) NextNPCID() EntityID
//
// Note: For a distributed system, you'd use:
// - Snowflake IDs (timestamp + node + sequence)
// - Or NATS KV atomic counter
// - Or pre-allocated ID ranges per shard

// =============================================================================
// SPATIAL QUERIES
// =============================================================================

// TODO: Implement spatial queries (or delegate to Zone):
//
// func (m *Manager) GetEntitiesInRange(pos math.Vec3, radius float32) []Entity
//   - Used for AOI, spell targeting, etc.
//
// func (m *Manager) GetPlayersInRange(pos math.Vec3, radius float32) []*Player
//   - Specifically players (for broadcasts)
//
// func (m *Manager) GetEntitiesInZone(zoneID uint32) []Entity

// =============================================================================
// PERSISTENCE
// =============================================================================

// TODO: Implement state persistence:
//
// func (m *Manager) SaveEntityState(e Entity) error
//   - Serialize to JSON/protobuf
//   - Write to NATS KV
//   - Key format: "entity:{type}:{id}"
//
// func (m *Manager) LoadEntityState(id EntityID) (Entity, error)
//   - Read from KV
//   - Deserialize and reconstruct
//
// Note: Player characters also go to MySQL for durability.
// NATS KV is for hot state that can be reconstructed.
