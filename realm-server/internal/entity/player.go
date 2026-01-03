package entity

import (
	"realm-server/internal/net"
)

// Player represents a connected player character.
// Each player has:
// - A network session for sending/receiving packets
// - Character data loaded from database
// - Inventory, quests, social connections, etc.
type Player struct {
	*BaseEntity

	// Network session - nil if disconnected
	Session *net.Session

	// Character data (loaded from DB)
	Name       string
	Race       uint8
	Class      uint8
	Gender     uint8
	Faction    uint8 // 0 = Alliance, 1 = Horde (example)
	Stats      Stats
	Experience uint32
	Money      uint32

	// TODO: Add these as you need them:
	// Inventory   *Inventory
	// Equipment   *Equipment
	// SpellBook   *SpellBook
	// QuestLog    *QuestLog
	// Guild       *GuildMember
	// Group       *GroupMember
	// Talents     *TalentSpec
	// Cooldowns   *CooldownManager
}

func NewPlayer(id EntityID, session *net.Session) *Player {
	// - Create BaseEntity
	// - Attach session
	// - Initialize empty character data (will be loaded from DB)

	baseEntity := NewBaseEntity(id)

	return &Player{
		BaseEntity: baseEntity,
		Session:    session,
	}
}

// func (p *Player) Update(dt float64)
//   - Process any pending actions
//   - Update combat timers
//   - Regenerate health/mana
//   - Check buff expirations
//
// func (p *Player) SendPacket(pkt *net.Packet)
//   - Send packet to client via session
//   - Handle nil session gracefully
//
// func (p *Player) Disconnect(reason string)
//   - Save character to DB
//   - Remove from world
//   - Close session
//
// func (p *Player) LoadFromDB(charData *db.CharacterData) error
//   - Populate player fields from database record
//
// func (p *Player) SaveToDB() error
//   - Persist current state to database

// =============================================================================
// PLAYER-SPECIFIC MOVEMENT
// =============================================================================

// Movement speed constants (yards per second, WoW-style)
const (
	BaseWalkSpeed   float32 = 2.5
	BaseRunSpeed    float32 = 7.0
	BaseSwimSpeed   float32 = 4.7
	BaseFlySpeed    float32 = 7.0
	MountedSpeed100 float32 = 14.0 // 100% mount
	MountedSpeed150 float32 = 21.0 // 150% mount (epic ground)
	MountedSpeed310 float32 = 21.7 // 310% flying
)

// TODO: Implement movement speed calculation:
//
// func (p *Player) GetMovementSpeed() float32
//   - Base speed for current movement type
//   - Apply mount speed if mounted
//   - Apply speed buffs/debuffs
//   - Apply slowing effects (dazed, snared, etc.)

// =============================================================================
// PLAYER VISIBILITY
// =============================================================================

// TODO: Implement visibility tracking:
//
// func (p *Player) GetVisibleEntities() []EntityID
//   - Return entities this player can currently see
//   - Used for AOI updates
//
// func (p *Player) CanSee(other Entity) bool
//   - Distance check
//   - Stealth/invisibility check
//   - Phasing check (quest state)
