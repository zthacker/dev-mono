package interest

// Uncomment when implementing:
// import (
// 	"realm-server/internal/entity"
// 	"realm-server/internal/net"
// 	"realm-server/pkg/math"
// )

// =============================================================================
// MESSAGE BROADCASTING
// =============================================================================
//
// When something happens in the world, only nearby players need to know.
// Broadcasting routes messages based on interest/visibility.
//
// Message types with different broadcast patterns:
// - Movement: Everyone who can see the mover
// - Combat: Everyone who can see attacker OR target
// - Chat: Depends on channel (say=nearby, yell=zone, guild=all members)
// - Spell effects: Everyone who can see the effect location
// - Loot: Only the looter (maybe group)

// Broadcaster handles routing messages to interested players.
type Broadcaster struct {
	// aoi *AOIManager
	// zoneMgr *zone.Manager
}

// TODO: Implement Broadcaster:
//
// func NewBroadcaster(aoi *AOIManager) *Broadcaster

// =============================================================================
// BROADCAST PATTERNS
// =============================================================================

// TODO: Implement broadcast methods:
//
// func (b *Broadcaster) ToNearby(pos math.Vec3, radius float32, pkt *net.Packet, exclude entity.EntityID)
//   - Send to all players within radius of pos
//   - Optionally exclude one player (usually the source)
//   - Uses zone grid for efficient lookup
//
// func (b *Broadcaster) ToVisible(source entity.Entity, pkt *net.Packet, includeSelf bool)
//   - Send to all players who can see source
//   - Uses AOI visibility sets
//
// func (b *Broadcaster) ToZone(zoneID uint32, pkt *net.Packet)
//   - Send to all players in zone
//   - For zone-wide announcements
//
// func (b *Broadcaster) ToAll(pkt *net.Packet)
//   - Send to all connected players
//   - Server announcements, maintenance warnings
//
// func (b *Broadcaster) ToPlayer(playerID entity.EntityID, pkt *net.Packet)
//   - Send to specific player
//
// func (b *Broadcaster) ToGroup(groupID uint32, pkt *net.Packet)
//   - Send to all party/raid members
//
// func (b *Broadcaster) ToGuild(guildID uint32, pkt *net.Packet)
//   - Send to all online guild members

// =============================================================================
// ENTITY UPDATES
// =============================================================================

// UpdateType categorizes what changed about an entity.
type UpdateType uint8

const (
	UpdateCreate   UpdateType = 1 // Entity entered visibility
	UpdateDestroy  UpdateType = 2 // Entity left visibility
	UpdateMovement UpdateType = 3 // Position/orientation changed
	UpdateValues   UpdateType = 4 // Stats/health/state changed
)

// TODO: Implement entity update broadcasting:
//
// func (b *Broadcaster) BroadcastEntityCreate(e entity.Entity)
//   - Build SMSG_UPDATE_OBJECT packet with full entity data
//   - Send to all players who can see it
//
// func (b *Broadcaster) BroadcastEntityDestroy(entityID entity.EntityID, lastPos math.Vec3)
//   - Build SMSG_DESTROY_OBJECT packet
//   - Send to players who could see it
//
// func (b *Broadcaster) BroadcastMovement(e entity.Entity, oldPos math.Vec3)
//   - Build SMSG_MOVE_UPDATE packet
//   - Send to visible players
//   - Handle case where entity entered/left visibility of some players
//
// func (b *Broadcaster) BroadcastStatChange(e entity.Entity, changedFields []uint32)
//   - Build partial update with only changed fields
//   - Send to visible players

// =============================================================================
// EFFICIENT BATCHING
// =============================================================================

// TODO: Implement packet batching to reduce syscalls:
//
// type BatchedBroadcast struct {
//     packets  map[entity.EntityID][]*net.Packet
// }
//
// func (b *Broadcaster) BeginBatch() *BatchedBroadcast
// func (bb *BatchedBroadcast) Queue(playerID entity.EntityID, pkt *net.Packet)
// func (bb *BatchedBroadcast) Flush()
//   - Combine queued packets per player
//   - Send in batches to reduce overhead
