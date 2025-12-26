package interest

// Uncomment when implementing:
// import (
// 	"realm-server/internal/entity"
// 	"realm-server/pkg/math"
// )

// =============================================================================
// AREA OF INTEREST (AOI)
// =============================================================================
//
// AOI determines what each player can see. This is critical for:
// 1. Reducing network bandwidth (don't send updates for distant entities)
// 2. Reducing client load (don't render distant entities)
// 3. Hiding invisible/stealthed entities
// 4. Phasing (different players see different world states)
//
// WoW's view distance is ~100-200 yards depending on settings.
// That's the radius within which you receive entity updates.

// AOIManager tracks what each player can see.
type AOIManager struct {
	// View distance for updates
	ViewDistance float32

	// Per-player visibility sets
	// visibility map[entity.EntityID]*VisibilitySet
}

// VisibilitySet tracks which entities a player can currently see.
type VisibilitySet struct {
	PlayerID uint64 // entity.EntityID

	// Entities currently visible to this player
	// visible map[entity.EntityID]struct{}

	// Position when visibility was last calculated
	// lastPosX, lastPosY, lastPosZ float32
}

// TODO: Implement AOIManager:
//
// func NewAOIManager(viewDistance float32) *AOIManager
//
// func (m *AOIManager) RegisterPlayer(player *entity.Player)
//   - Create visibility set for player
//
// func (m *AOIManager) UnregisterPlayer(playerID entity.EntityID)
//   - Remove visibility set
//   - Notify was-visible entities? (optional)
//
// func (m *AOIManager) UpdateVisibility(player *entity.Player, nearbyEntities []entity.Entity) *VisibilityUpdate
//   Core algorithm:
//   1. Get current visible set
//   2. Compute new visible set from nearbyEntities + distance check + visibility rules
//   3. Diff: entered = new - current, left = current - new
//   4. Update stored visible set
//   5. Return VisibilityUpdate with entered/left lists

// VisibilityUpdate contains changes to a player's visible entities.
// Note: Use interface{} to avoid import cycles - cast to entity types when implementing.
type VisibilityUpdate struct {
	PlayerID uint64        // entity.EntityID
	Entered  []interface{} // []entity.Entity - Newly visible
	Left     []uint64      // []entity.EntityID - No longer visible
	Updated  []interface{} // []entity.Entity - Still visible but changed
}

// =============================================================================
// VISIBILITY RULES
// =============================================================================

// Not everything in range is visible. Additional rules:

// TODO: Implement visibility checks:
//
// func (m *AOIManager) CanSee(viewer *entity.Player, target entity.Entity) bool
//   Check in order:
//   1. Distance check (within ViewDistance)
//   2. Dead check (can't see if viewer is dead, unless ghost)
//   3. Stealth check (rogues/druids)
//   4. Invisibility check (mages, potions)
//   5. Phase check (quest state)
//   6. Faction check? (depends on game design)
//
// func (m *AOIManager) GetStealthDetectionRange(viewer *entity.Player, stealthLevel int) float32
//   - Higher level stealth = must be closer to detect
//   - Stealth detection abilities increase range
//
// func (m *AOIManager) GetPhase(player *entity.Player) uint32
//   - Return player's current phase based on quest state
//   - Players in different phases can't see each other

// =============================================================================
// PHASING
// =============================================================================

// Phasing allows the same location to look different based on quest progress.
// Example: Before a quest, a village is intact. After, it's destroyed.
// Players see the version matching their quest state.

// Phase represents a world state variation.
type Phase struct {
	ID        uint32
	Name      string
	// Entities only visible in this phase
	// Terrain modifications (if supported)
}

// TODO: Implement phasing:
//
// func (m *AOIManager) SetPlayerPhase(playerID entity.EntityID, phaseID uint32)
//
// func (m *AOIManager) GetEntityPhase(entityID entity.EntityID) uint32
//   - Return 0 for entities visible in all phases
//
// func (m *AOIManager) PhasesMatch(playerPhase, entityPhase uint32) bool
//   - Phase 0 matches all phases
//   - Otherwise must be equal
