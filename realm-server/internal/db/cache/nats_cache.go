package cache

// Uncomment when implementing:
// import (
// 	"context"
// 	"encoding/json"
// 	"time"

// 	"github.com/nats-io/nats.go/jetstream"

// 	"realm-server/internal/entity"
// 	"realm-server/pkg/math"
// )

// =============================================================================
// NATS KV CACHE LAYER
// =============================================================================
//
// NATS KV stores "hot" state that:
// 1. Changes frequently (position, health)
// 2. Needs fast access (sub-millisecond)
// 3. Can be reconstructed from MySQL if lost
//
// This is similar to your ses_state bucket in the satellite system.
//
// Key design:
// - MySQL is source of truth for character data
// - NATS KV is fast cache for runtime state
// - On crash: reload from MySQL, lose some recent changes (acceptable)

// Cache wraps NATS KV for game state caching.
type Cache struct {
	// Different buckets for different data types
	// Uncomment when implementing:
	// kvPlayers jetstream.KeyValue // Hot player state
	// kvSessions jetstream.KeyValue // Session data
	// kvWorld   jetstream.KeyValue // World state (spawns, events)
}

// PlayerState is the hot state cached in NATS KV.
// This is what changes every tick.
type PlayerState struct {
	EntityID    uint64  `json:"entity_id"`
	ZoneID      uint32  `json:"zone_id"`
	PositionX   float32 `json:"position_x"` // Flattened Vec3 for now
	PositionY   float32 `json:"position_y"`
	PositionZ   float32 `json:"position_z"`
	Orientation float32 `json:"orientation"`
	Health      int32   `json:"health"`
	Mana        int32   `json:"mana"`
	MoveFlags   uint32  `json:"move_flags"`
	LastUpdate  int64   `json:"last_update"` // Unix nano
}

// TODO: Implement Cache:
//
// func NewCache(js jetstream.JetStream) (*Cache, error)
//   - Create or get KV buckets:
//     - "realm_players" - TTL 1 hour (players get removed on logout anyway)
//     - "realm_sessions" - TTL 30 minutes
//     - "realm_world" - No TTL (persists until explicit delete)
//
// func (c *Cache) Close() error

// =============================================================================
// PLAYER STATE OPERATIONS
// =============================================================================

// TODO: Implement player state caching:
//
// func (c *Cache) GetPlayerState(ctx context.Context, entityID entity.EntityID) (*PlayerState, error)
//   - Key: fmt.Sprintf("player:%d", entityID)
//   - Return nil, nil if not found (not an error)
//
// func (c *Cache) SetPlayerState(ctx context.Context, state *PlayerState) error
//   - Serialize to JSON
//   - Put to KV
//
// func (c *Cache) DeletePlayerState(ctx context.Context, entityID entity.EntityID) error
//   - Called on logout/disconnect
//
// func (c *Cache) UpdatePlayerPosition(ctx context.Context, entityID entity.EntityID, pos math.Vec3, orientation float32) error
//   - Read-modify-write
//   - Or use NATS KV's optimistic locking (revision)

// =============================================================================
// SESSION STATE
// =============================================================================

// SessionState tracks authentication state.
type SessionState struct {
	AccountID   uint32 `json:"account_id"`
	CharacterID uint64 `json:"character_id"`
	ServerID    uint32 `json:"server_id"` // Which world server
	ConnectedAt int64  `json:"connected_at"`
}

// TODO: Implement session caching:
//
// func (c *Cache) GetSession(ctx context.Context, accountID uint32) (*SessionState, error)
//
// func (c *Cache) SetSession(ctx context.Context, state *SessionState) error
//
// func (c *Cache) DeleteSession(ctx context.Context, accountID uint32) error

// =============================================================================
// WORLD STATE
// =============================================================================

// World state includes things that persist across server restarts
// but change during gameplay.

// TODO: Implement world state:
//
// func (c *Cache) GetSpawnState(ctx context.Context, spawnID uint32) (*SpawnState, error)
//   - Is the mob alive? When does it respawn?
//
// func (c *Cache) SetSpawnState(ctx context.Context, state *SpawnState) error

// type SpawnState struct {
//     SpawnID       uint32 `json:"spawn_id"`
//     IsAlive       bool   `json:"is_alive"`
//     RespawnAt     int64  `json:"respawn_at"` // Unix timestamp
//     CurrentHealth int32  `json:"current_health"`
// }

// =============================================================================
// BULK OPERATIONS
// =============================================================================

// TODO: Implement bulk operations for efficiency:
//
// func (c *Cache) GetPlayersInZone(ctx context.Context, zoneID uint32) ([]*PlayerState, error)
//   - NATS KV doesn't have great query support
//   - Options:
//     a) Iterate all keys with prefix (slow)
//     b) Maintain a separate index (zone:123 -> [player IDs])
//     c) Use NATS JetStream for queries instead
//
// func (c *Cache) SaveAllPlayerStates(ctx context.Context, states []*PlayerState) error
//   - Batch write for periodic saves

// =============================================================================
// CONSISTENCY CONSIDERATIONS
// =============================================================================

// NATS KV provides:
// - Optimistic concurrency via revision numbers
// - Watch for change notifications
// - TTL for automatic expiration
//
// It does NOT provide:
// - Transactions across multiple keys
// - Complex queries
//
// For complex operations, use MySQL transactions.
// NATS KV is for hot, frequently-updated single-key data.
