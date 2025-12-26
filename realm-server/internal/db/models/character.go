package models

import (
	"time"
)

// =============================================================================
// CHARACTER DATA MODEL
// =============================================================================
//
// This is what gets persisted to MySQL.
// Hot runtime state goes to NATS KV.
//
// WoW character data is split across many tables:
// - characters (core data)
// - character_inventory
// - character_spells
// - character_quests
// - character_reputation
// - character_skills
// - etc.

// Character is the core character record.
type Character struct {
	// Identity
	ID        uint64 `db:"id"`         // Primary key
	AccountID uint32 `db:"account_id"` // Owner account
	Name      string `db:"name"`       // Unique per realm

	// Appearance
	Race      uint8 `db:"race"`
	Class     uint8 `db:"class"`
	Gender    uint8 `db:"gender"`
	Skin      uint8 `db:"skin"`
	Face      uint8 `db:"face"`
	HairStyle uint8 `db:"hair_style"`
	HairColor uint8 `db:"hair_color"`

	// Position (where they logged out)
	MapID       uint32  `db:"map_id"`
	ZoneID      uint32  `db:"zone_id"`
	PositionX   float32 `db:"position_x"`
	PositionY   float32 `db:"position_y"`
	PositionZ   float32 `db:"position_z"`
	Orientation float32 `db:"orientation"`

	// Progression
	Level      uint8  `db:"level"`
	Experience uint32 `db:"xp"`
	Money      uint32 `db:"money"` // Copper

	// Stats
	Health    int32 `db:"health"`
	MaxHealth int32 `db:"max_health"`
	Mana      int32 `db:"mana"`
	MaxMana   int32 `db:"max_mana"`

	// Timestamps
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	LastLogin time.Time `db:"last_login"`
	TotalTime uint32    `db:"total_time"` // Seconds played

	// Flags
	AtLogin   uint32 `db:"at_login"` // Flags for things to do at login (rename, customize, etc.)
	IsDeleted bool   `db:"is_deleted"`
}

// CharacterInventory represents an item in a character's inventory.
type CharacterInventory struct {
	CharacterID uint64 `db:"character_id"`
	Bag         uint8  `db:"bag"`   // 0 = backpack, 1-4 = equipped bags
	Slot        uint8  `db:"slot"`  // Slot within bag
	ItemID      uint32 `db:"item_id"` // Reference to item template
	StackCount  uint16 `db:"stack_count"`
	// Durability, enchants, gems, etc. would go here
}

// CharacterSpell represents a learned spell.
type CharacterSpell struct {
	CharacterID uint64 `db:"character_id"`
	SpellID     uint32 `db:"spell_id"`
	Active      bool   `db:"active"` // Some spells can be "unlearned" temporarily
}

// CharacterQuest tracks quest progress.
type CharacterQuest struct {
	CharacterID uint64 `db:"character_id"`
	QuestID     uint32 `db:"quest_id"`
	Status      uint8  `db:"status"` // 0=incomplete, 1=complete, 2=failed
	// Objective progress would be additional columns or JSON
}

// =============================================================================
// TEMPLATE DATA (read-only, loaded at startup)
// =============================================================================

// ItemTemplate defines an item type (loaded from DB, cached in memory).
type ItemTemplate struct {
	ID          uint32 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Class       uint8  `db:"class"`    // Weapon, armor, consumable, etc.
	SubClass    uint8  `db:"subclass"` // Sword, axe, cloth, plate, etc.
	Quality     uint8  `db:"quality"`  // Common, uncommon, rare, epic, legendary
	Level       uint8  `db:"level"`
	RequiredLevel uint8 `db:"required_level"`
	// Stats, on-use effects, etc.
}

// SpellTemplate defines a spell/ability.
type SpellTemplate struct {
	ID           uint32  `db:"id"`
	Name         string  `db:"name"`
	Description  string  `db:"description"`
	SchoolMask   uint8   `db:"school_mask"` // Fire, frost, nature, etc.
	ManaCost     uint32  `db:"mana_cost"`
	CastTime     uint32  `db:"cast_time"`   // Milliseconds
	Cooldown     uint32  `db:"cooldown"`    // Milliseconds
	Range        float32 `db:"range"`
	// Effects would be in a separate table (spell_effects)
}

// NPCTemplate defines an NPC type.
type NPCTemplate struct {
	ID           uint32  `db:"id"`
	Name         string  `db:"name"`
	SubName      string  `db:"subname"` // "<Innkeeper>", "<Quest>"
	Level        uint8   `db:"level"`
	Health       int32   `db:"health"`
	Mana         int32   `db:"mana"`
	Faction      uint32  `db:"faction"`
	NPCFlags     uint32  `db:"npc_flags"` // Vendor, quest giver, etc.
	// AI script, loot table, etc.
}
