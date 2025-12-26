package zone

// Uncomment when implementing:
// import (
// 	"context"

// 	"github.com/nats-io/nats.go/jetstream"
// )

// Manager handles all zones on this world server.
// Responsible for:
// - Loading zone data on startup
// - Managing zone lifecycle
// - Routing entities between zones
// - Coordinating with other servers for zone handoffs
type Manager struct {
	// zones map[uint32]*Zone

	// Shard managers for each zone (if using dynamic sharding)
	// shardMgrs map[uint32]*ShardManager

	// NATS for cross-server zone communication
	// js jetstream.JetStream
}

// TODO: Implement Manager:
//
// func NewManager() *Manager
//
// func (m *Manager) Start(js jetstream.JetStream, zoneIDs []uint32) error
//   1. Connect to NATS
//   2. For each zoneID:
//      - Load zone data from DB/files
//      - Initialize grid
//      - Load spawn points
//      - Spawn initial NPCs
//   3. Subscribe to zone-related NATS subjects
//
// func (m *Manager) Shutdown(ctx context.Context) error
//   - Save zone states
//   - Cleanup resources
//
// func (m *Manager) GetZone(id uint32) (*Zone, bool)
//
// func (m *Manager) UpdateAll(dt float64)
//   - Called each tick
//   - Update each zone

// =============================================================================
// ZONE LOADING
// =============================================================================

// ZoneDataLoader provides zone data from database/files.
type ZoneDataLoader interface {
	LoadZoneInfo(zoneID uint32) (*ZoneInfo, error)
	LoadSpawnPoints(zoneID uint32) ([]SpawnPoint, error)
	LoadAreaTriggers(zoneID uint32) ([]AreaTrigger, error)
}

type ZoneInfo struct {
	ID     uint32
	Name   string
	MinX   float32
	MinY   float32
	MinZ   float32
	MaxX   float32
	MaxY   float32
	MaxZ   float32
}

// TODO: Implement zone loading:
//
// func (m *Manager) LoadZone(loader ZoneDataLoader, zoneID uint32) (*Zone, error)
//   1. Get zone info
//   2. Create zone with bounds
//   3. Load spawn points
//   4. Load area triggers
//   5. Initialize grid
//   6. Spawn NPCs

// =============================================================================
// ZONE TRANSFERS
// =============================================================================

// TODO: Implement zone transfers:
//
// func (m *Manager) TransferEntity(entityID entity.EntityID, fromZone, toZone uint32, pos math.Vec3) error
//   If toZone is local:
//   1. Remove from fromZone
//   2. Add to toZone at pos
//   3. Update entity's zone ID
//
//   If toZone is on different server:
//   1. Serialize entity state
//   2. Publish transfer request to NATS
//   3. Remove from local zone after confirmation

// =============================================================================
// NATS INTEGRATION
// =============================================================================

// Zone-related NATS subjects:
// - zone.{zoneID}.enter      - Entity entering zone
// - zone.{zoneID}.leave      - Entity leaving zone
// - zone.{zoneID}.update     - Entity position/state update (for ghosts)
// - zone.{zoneID}.spawn      - NPC spawned
// - zone.{zoneID}.event      - Zone-wide event (boss spawn, weather change)

// TODO: Implement NATS handlers:
//
// func (m *Manager) handleZoneEnter(msg *nats.Msg)
//   - Deserialize entity
//   - Add to appropriate zone
//   - Notify nearby players
//
// func (m *Manager) handleZoneLeave(msg *nats.Msg)
//
// func (m *Manager) publishZoneEvent(zoneID uint32, event ZoneEvent)
