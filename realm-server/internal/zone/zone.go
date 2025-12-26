package zone

import (
	"realm-server/internal/entity"
	"realm-server/pkg/math"
)

// Zone represents a game area (like a WoW zone: Elwynn Forest, Stormwind, etc.)
// A zone contains:
// - Static geometry/collision (loaded from files)
// - Spawn points for NPCs/mobs
// - Area triggers (for quests, events, zone transitions)
// - Currently active entities
//
// Zones can be sharded - multiple instances of the same zone running
// on different servers to handle population.
type Zone struct {
	ID     uint32
	Name   string
	Bounds math.WorldBounds

	// Spatial partitioning for efficient queries
	// grid *Grid

	// Entities currently in this zone
	// entities map[entity.EntityID]entity.Entity

	// Static data (loaded from DB/files)
	// spawns   []SpawnPoint
	// triggers []AreaTrigger

	// Shard info (if this is a shard of a larger zone)
	ShardID     uint32
	ShardBounds math.WorldBounds // Sub-region if sharded
}

// SpawnPoint defines where NPCs/mobs appear.
type SpawnPoint struct {
	ID         uint32
	TemplateID uint32          // NPC/Mob template
	Position   math.Vec3
	Orientation float32
	RespawnTime uint32         // Seconds
	WanderRadius float32       // How far from spawn point

	// State
	SpawnedEntityID entity.EntityID // Currently spawned entity (0 if dead)
	RespawnAt       int64           // Unix timestamp when to respawn
}

// AreaTrigger defines a region that fires events when entered.
type AreaTrigger struct {
	ID       uint32
	Bounds   math.WorldBounds
	TriggerType TriggerType
	TargetID uint32 // Quest ID, zone ID, script ID, etc.
}

type TriggerType uint8

const (
	TriggerQuest       TriggerType = 1 // Quest objective
	TriggerZoneChange  TriggerType = 2 // Teleport to another zone
	TriggerScript      TriggerType = 3 // Run custom script
	TriggerSanctuary   TriggerType = 4 // No PvP area
	TriggerInstance    TriggerType = 5 // Dungeon/raid entrance
)

// TODO: Implement Zone:
//
// func NewZone(id uint32, name string, bounds math.WorldBounds) *Zone
//
// func (z *Zone) LoadData(loader ZoneDataLoader) error
//   - Load spawn points from DB
//   - Load area triggers
//   - Initialize grid
//
// func (z *Zone) AddEntity(e entity.Entity)
//   - Add to entities map
//   - Add to spatial grid
//
// func (z *Zone) RemoveEntity(id entity.EntityID)
//   - Remove from entities
//   - Remove from grid
//
// func (z *Zone) GetEntity(id entity.EntityID) (entity.Entity, bool)
//
// func (z *Zone) Update(dt float64)
//   - Check spawn timers
//   - Respawn dead mobs
//   - Process zone events

// =============================================================================
// SPATIAL QUERIES
// =============================================================================

// TODO: Implement spatial queries (delegate to Grid):
//
// func (z *Zone) GetEntitiesInRange(pos math.Vec3, radius float32) []entity.Entity
//
// func (z *Zone) GetPlayersInRange(pos math.Vec3, radius float32) []*entity.Player
//
// func (z *Zone) GetEntitiesInBounds(bounds math.WorldBounds) []entity.Entity
//
// func (z *Zone) CheckAreaTriggers(e entity.Entity) []AreaTrigger
//   - Check if entity is inside any triggers
//   - Return newly entered triggers

// =============================================================================
// ZONE TRANSITIONS
// =============================================================================

// TODO: Implement zone transitions:
//
// func (z *Zone) TransferEntity(e entity.Entity, targetZoneID uint32, position math.Vec3) error
//   - Remove from this zone
//   - Notify zone manager
//   - If target zone is on different server, serialize and send via NATS
