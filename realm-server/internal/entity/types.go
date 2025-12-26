package entity

import (
	"time"

	"realm-server/pkg/math"
)

// =============================================================================
// ENTITY IDENTIFICATION
// =============================================================================
// EntityID is a globally unique identifier for any entity.
// High bits encode type, low bits encode instance ID.
// Format: [8 bits type][56 bits id]
//
// - Single uint64 is fast to compare, hash, and transmit
// - Type bits let you identify entity kind without a lookup
// - 56 bits = 72 quadrillion unique IDs per type
// - This means that if the ID uses the full uint64, there's going to be a problem, but that shouldn't happen
// - Visuals
// - Type (1) shifted: 0x0100000000000000
// - ID Masked:        0x00FFFFFFFFFFFFFF
// - Result:           0x01FFFFFFFFFFFFFF

type EntityID uint64

type EntityType uint8

const (
	EntityTypePlayer EntityType = 1
	EntityTypeNPC    EntityType = 2
	EntityTypeMob    EntityType = 3
	EntityTypeObject EntityType = 4
	EntityTypeItem   EntityType = 5
)

// - NewEntityID(typ EntityType, id uint64) EntityID - pack type and id
func NewEntityID(typ EntityType, id uint64) EntityID {
	// Shift type to high 8 bits, mask id to low 56 bits, combine with OR
	return EntityID(uint64(typ)<<56 | (id & 0x00FFFFFFFFFFFFFF))
}

// - (e EntityID) Type() EntityType - extract type bits
func (e EntityID) Type() EntityType {
	// Shift right 56 bits to get the type in low 8 bits
	return EntityType(e >> 56)
}

// - (e EntityID) ID() uint64 - extract id bits
func (e EntityID) ID() uint64 {
	// Mask ff the high 8 bits, keep low 56
	return uint64(e) & 0x00FFFFFFFFFFFFFF
}

// - (e EntityID) IsPlayer() bool
func (e EntityID) IsPlayer() bool {
	return e.Type() == EntityTypePlayer
}

// =============================================================================
// POSITION & MOVEMENT
// =============================================================================

// Transform represents an entity's position and orientation in the world.
type Transform struct {
	Position    math.Vec3
	Orientation float32 // Radians, 0 = North, clockwise
}

// MovementState tracks current movement for interpolation and validation.
// This is updated frequently (every client movement packet).
type MovementState struct {
	Transform
	Velocity     math.Vec3
	MoveFlags    MoveFlags
	JumpVelocity float32   // Vertical velocity when jumping
	FallTime     float32   // How long falling (for fall damage)
	LastUpdate   time.Time // For server-side interpolation
}

// MoveFlags are bitflags for movement state.
// These are synced between client and server.
type MoveFlags uint32

const (
	MoveFlagNone     MoveFlags = 0
	MoveFlagForward  MoveFlags = 1 << 0
	MoveFlagBackward MoveFlags = 1 << 1
	MoveFlagLeft     MoveFlags = 1 << 2
	MoveFlagRight    MoveFlags = 1 << 3
	MoveFlagJumping  MoveFlags = 1 << 4
	MoveFlagFalling  MoveFlags = 1 << 5
	MoveFlagSwimming MoveFlags = 1 << 6
	MoveFlagFlying   MoveFlags = 1 << 7
	MoveFlagMounted  MoveFlags = 1 << 8
	MoveFlagRooted   MoveFlags = 1 << 9  // Cannot move (spell effect)
	MoveFlagStunned  MoveFlags = 1 << 10 // Cannot act
	MoveFlagDead     MoveFlags = 1 << 11
	MoveFlagGhost    MoveFlags = 1 << 12
	MoveFlagInCombat MoveFlags = 1 << 13
)

// - (f MoveFlags) Has(flag MoveFlags) bool
func (f MoveFlags) Has(flag MoveFlags) bool {
	return f&flag == f
}

// - (f *MoveFlags) Set(flag MoveFlags)
// Set turns on the specified bits in the MoveFlags.
func (f *MoveFlags) Set(flag MoveFlags) {
	*f |= flag
}

// - (f *MoveFlags) Clear(flag MoveFlags)
// Clear turns off the specified bits in the MoveFlags.
func (f *MoveFlags) Clear(flag MoveFlags) {
	*f &^= flag
}

// =============================================================================
// ENTITY STATS
// =============================================================================

// Stats represents entity attributes.
// These are the "derived" stats after gear/buffs are applied.
type Stats struct {
	Health    int32
	MaxHealth int32
	Mana      int32
	MaxMana   int32
	Level     uint8

	// Primary stats
	Strength  int16
	Agility   int16
	Intellect int16
	Stamina   int16
	Spirit    int16
}

// TODO: Consider adding:
// - BaseStats vs BonusStats separation
// - Stat modifiers (buffs/debuffs)
// - Derived combat stats (attack power, spell power, crit, etc.)

// =============================================================================
// ENTITY INTERFACE
// =============================================================================

// Entity is the base interface all world objects implement.
// This allows the world server to treat all entities uniformly.
type Entity interface {
	ID() EntityID
	Type() EntityType
	Transform() Transform
	SetTransform(t Transform)
	Movement() *MovementState
	ZoneID() uint32
	SetZoneID(id uint32)

	// Update is called each server tick.
	// dt is seconds since last tick (typically 0.05 for 20 tick/sec).
	Update(dt float64)
}
