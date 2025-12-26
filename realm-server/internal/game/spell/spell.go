package spell

import (
	"realm-server/internal/entity"
	"realm-server/pkg/math"
)

// =============================================================================
// SPELL SYSTEM
// =============================================================================
//
// Spells in WoW are complex. Each spell can have multiple effects:
// - Direct damage/healing
// - Apply aura (buff/debuff)
// - Summon creature
// - Teleport
// - Create item
// - etc.
//
// Spell casting flow:
// 1. Player sends CMSG_CAST_SPELL
// 2. Validate: known spell, resources, cooldown, target
// 3. Start cast (if has cast time)
// 4. Broadcast SMSG_SPELL_START
// 5. On cast complete: apply effects
// 6. Broadcast SMSG_SPELL_GO

// SpellManager handles spell casting.
type SpellManager struct {
	// Spell templates loaded from DB
	// templates map[uint32]*SpellTemplate

	// Active casts (for interrupting)
	// activeCasts map[entity.EntityID]*ActiveCast
}

// SpellTemplate defines a spell's properties.
type SpellTemplate struct {
	ID          uint32
	Name        string
	SchoolMask  uint8   // Fire, Frost, Nature, etc.
	ManaCost    uint32
	CastTime    uint32  // Milliseconds (0 = instant)
	Cooldown    uint32  // Milliseconds
	Range       float32 // Max range (0 = self only)

	// Effects (simplified - real WoW has up to 3 effects per spell)
	// Effects []SpellEffect
}

// ActiveCast tracks a spell being cast.
type ActiveCast struct {
	Caster   entity.EntityID
	SpellID  uint32
	TargetID entity.EntityID
	TargetPos math.Vec3 // For ground-target spells
	StartTime int64
	EndTime   int64
}

// TODO: Implement SpellManager:
//
// func NewSpellManager() *SpellManager
//
// func (sm *SpellManager) LoadTemplates(templates map[uint32]*SpellTemplate)
//
// func (sm *SpellManager) CanCast(caster *entity.Player, spellID uint32, target entity.Entity) error
//   Checks:
//   - Player knows spell
//   - Not on cooldown
//   - Has mana/resources
//   - Target is valid
//   - In range
//   - Line of sight (optional)
//   - Not silenced/interrupted
//
// func (sm *SpellManager) StartCast(caster *entity.Player, spellID uint32, target entity.Entity) (*ActiveCast, error)
//   - Create active cast
//   - Schedule completion
//   - Return cast info (for SMSG_SPELL_START)
//
// func (sm *SpellManager) CompleteCast(cast *ActiveCast) error
//   - Apply spell effects
//   - Consume mana
//   - Start cooldown
//
// func (sm *SpellManager) InterruptCast(casterID entity.EntityID, reason uint8)
//   - Cancel active cast
//   - Partial cooldown for interrupted spells

// =============================================================================
// SPELL EFFECTS
// =============================================================================

// EffectType identifies what a spell effect does.
type EffectType uint8

const (
	EffectDamage       EffectType = 1
	EffectHeal         EffectType = 2
	EffectApplyAura    EffectType = 3
	EffectTeleport     EffectType = 4
	EffectSummon       EffectType = 5
	EffectEnergize     EffectType = 6  // Restore mana/energy
	EffectDispel       EffectType = 7
	EffectKnockback    EffectType = 8
	EffectCreateItem   EffectType = 9
	// Many more in real WoW...
)

// SpellEffect represents one effect of a spell.
type SpellEffect struct {
	Type     EffectType
	Value    int32   // Base value (damage, heal amount, etc.)
	Radius   float32 // For AoE
	AuraID   uint32  // For apply aura effect
}

// TODO: Implement effect handlers:
//
// func (sm *SpellManager) ApplyEffect(caster, target entity.Entity, effect *SpellEffect)
//   switch effect.Type {
//   case EffectDamage:
//       // Calculate and apply damage
//   case EffectHeal:
//       // Calculate and apply healing
//   case EffectApplyAura:
//       // Apply buff/debuff
//   // etc.
//   }

// =============================================================================
// AURAS (Buffs/Debuffs)
// =============================================================================

// Aura represents an active buff or debuff on an entity.
type Aura struct {
	ID        uint32
	SpellID   uint32          // Spell that applied this
	CasterID  entity.EntityID // Who applied it
	Duration  int32           // Remaining duration (ms), -1 = permanent
	Stacks    uint8           // Stack count

	// Periodic effects
	TickPeriod int32 // For DoTs/HoTs (ms per tick)
	NextTick   int64 // When next tick fires
}

// AuraType categorizes auras.
type AuraType uint8

const (
	AuraPeriodicDamage  AuraType = 1  // DoT
	AuraPeriodicHeal    AuraType = 2  // HoT
	AuraModStat         AuraType = 3  // Stat buff/debuff
	AuraModSpeed        AuraType = 4  // Speed buff/debuff
	AuraModDamage       AuraType = 5  // Damage increase/decrease
	AuraRoot            AuraType = 6  // Cannot move
	AuraStun            AuraType = 7  // Cannot act
	AuraSilence         AuraType = 8  // Cannot cast
	AuraInvisibility    AuraType = 9
	AuraStealth         AuraType = 10
	// Many more...
)

// TODO: Implement aura management:
//
// type AuraManager struct {
//     // Per-entity aura lists
//     auras map[entity.EntityID][]*Aura
// }
//
// func (am *AuraManager) AddAura(target entity.EntityID, aura *Aura) error
//   - Check for existing aura (refresh or stack)
//   - Add to list
//   - Apply initial effects
//
// func (am *AuraManager) RemoveAura(target entity.EntityID, auraID uint32)
//
// func (am *AuraManager) Update(dt float64)
//   - Tick all periodic auras
//   - Remove expired auras
