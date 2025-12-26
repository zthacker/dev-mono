package zone

// Uncomment when implementing:
// import (
// 	"realm-server/internal/entity"
// 	"realm-server/pkg/math"
// )

// =============================================================================
// DYNAMIC SHARDING
// =============================================================================
//
// WoW's sharding dynamically splits zones based on population.
// When too many players are in one area:
// 1. Server creates a new shard (instance of the zone)
// 2. New players joining the area go to the less populated shard
// 3. Players can be grouped to stay on the same shard (party/raid)
//
// This is different from static sharding (your satellite partitions)
// because it's dynamic and players can move between shards.

// ShardManager handles dynamic sharding for a zone.
type ShardManager struct {
	ZoneID      uint32
	MaxPlayers  int // Per-shard player limit

	// Active shards for this zone
	// shards map[uint32]*Shard

	// Shard assignment
	// playerShards map[entity.EntityID]uint32 // Which shard each player is on
}

// Shard is one instance of a zone.
type Shard struct {
	ID         uint32
	ZoneID     uint32
	Zone       *Zone          // The actual zone instance
	Population int            // Current player count
	ServerID   uint32         // Which world server runs this shard
}

// TODO: Implement ShardManager:
//
// func NewShardManager(zoneID uint32, maxPlayersPerShard int) *ShardManager
//
// func (sm *ShardManager) AssignShard(player *entity.Player, preferredShard uint32) *Shard
//   Logic:
//   1. If player has party members, put on same shard as party leader
//   2. If preferredShard specified and has room, use it
//   3. Otherwise, find shard with lowest population
//   4. If all shards are full, create new shard
//
// func (sm *ShardManager) CreateShard() *Shard
//   - Allocate new shard ID
//   - Initialize zone instance
//   - Register with world server
//
// func (sm *ShardManager) DestroyShard(shardID uint32)
//   - Called when shard is empty
//   - Cleanup resources
//
// func (sm *ShardManager) TransferBetweenShards(player *entity.Player, fromShard, toShard uint32) error
//   - Remove from source shard
//   - Add to target shard
//   - Update player's visible entities (they see different people now)

// =============================================================================
// CROSS-SHARD VISIBILITY
// =============================================================================

// Some things should be visible across all shards:
// - World bosses
// - Major city NPCs
// - Guild/group members (optionally phase them together)

// TODO: Implement cross-shard entities:
//
// func (sm *ShardManager) BroadcastToAllShards(msg interface{})
//   - Send update to all shards of this zone
//
// func (sm *ShardManager) SyncWorldBoss(boss *entity.Mob)
//   - Keep boss state synchronized across shards
//   - Only one shard "owns" the boss for combat

// =============================================================================
// SHARD BALANCING
// =============================================================================

// TODO: Implement automatic shard balancing:
//
// func (sm *ShardManager) Rebalance()
//   Called periodically to:
//   1. Merge underpopulated shards
//   2. Split overpopulated shards
//   3. Move AFK players to lower-pop shards
//
// Constraints:
// - Don't split players who are in combat
// - Don't split players who are grouped
// - Provide seamless transition (ideally invisible to player)

// =============================================================================
// REGION-BASED SHARDING
// =============================================================================

// For very large zones, you can shard by region instead of population.
// Each shard handles a geographic sub-region of the zone.

// RegionShard handles a geographic portion of a zone.
type RegionShard struct {
	ID     uint32
	ZoneID uint32
	// Region math.WorldBounds // This shard's area - uncomment when implementing

	// Border handling
	// Entities near borders need to be visible to adjacent shards
	BorderWidth float32
}

// TODO: Implement border handling:
//
// func (rs *RegionShard) IsInBorder(pos math.Vec3) bool
//   - True if position is within BorderWidth of region edge
//
// func (rs *RegionShard) GetBorderEntities() []entity.Entity
//   - Entities that should be replicated to adjacent shards
//
// func (rs *RegionShard) ReceiveBorderUpdate(fromShard uint32, entities []entity.Entity)
//   - Update "ghost" entities from neighboring shard
