package entity

// NPC represents a non-player character (vendors, quest givers, etc.)
// NPCs are typically spawned from database templates.
type NPC struct {
	*BaseEntity

	// Template ID references the NPC definition in the database
	// Multiple NPC instances can share the same template
	TemplateID uint32

	// Instance-specific data
	Stats       Stats
	SpawnPoint  Transform // Where to respawn
	RespawnTime uint32    // Seconds until respawn

	// AI state
	AIState    AIState
	TargetID   EntityID // Current target (if in combat)
	HomeRadius float32  // How far from spawn before leashing
}

type AIState uint8

const (
	AIStateIdle       AIState = 0
	AIStatePatrolling AIState = 1
	AIStateChasing    AIState = 2
	AIStateCombat     AIState = 3
	AIStateFleeing    AIState = 4
	AIStateReturning  AIState = 5 // Leashing back to spawn
	AIStateDead       AIState = 6
)

// TODO: Implement NPC:
//
// func NewNPC(id EntityID, templateID uint32) *NPC
//   - Create BaseEntity
//   - Load stats from template
//
// func (n *NPC) Update(dt float64)
//   - Run AI state machine
//   - Process combat if in combat
//   - Handle respawn timer if dead
//
// func (n *NPC) LoadTemplate(template *db.NPCTemplate)
//   - Copy base stats, faction, model info, etc.

// =============================================================================
// NPC AI
// =============================================================================

// TODO: Implement AI behavior:
//
// func (n *NPC) Think(dt float64)
//   - State machine for AI behavior
//   - Called each tick when not in combat
//
// func (n *NPC) OnAggro(attacker Entity)
//   - Enter combat state
//   - Alert nearby friendly NPCs (social aggro)
//
// func (n *NPC) OnDeath(killer Entity)
//   - Drop loot
//   - Give XP/quest credit
//   - Start respawn timer
//
// func (n *NPC) Leash()
//   - Return to spawn point
//   - Reset health
//   - Clear threat table

// =============================================================================
// MOB (HOSTILE NPC)
// =============================================================================

// Mob extends NPC with combat-specific behavior.
type Mob struct {
	*NPC

	// Threat table - maps EntityID to threat value
	// Higher threat = more likely to be attacked
	// TODO: ThreatTable map[EntityID]float32

	// Combat timers
	LastAttack   int64 // Unix timestamp of last auto-attack
	SpecialTimer int64 // Timer for special abilities

	// Loot
	LootTableID uint32
	LootLocked  bool     // True while being looted
	LootedBy    EntityID // Who has loot rights
}

// TODO: Implement Mob:
//
// func NewMob(id EntityID, templateID uint32) *Mob
//
// func (m *Mob) AddThreat(source EntityID, amount float32)
//   - Add to threat table
//   - May cause target switch
//
// func (m *Mob) GetTopThreat() EntityID
//   - Return entity with highest threat
//
// func (m *Mob) ClearThreat()
//   - Called on leash/death
//
// func (m *Mob) DropLoot() []Item
//   - Roll loot table
//   - Assign loot rights
