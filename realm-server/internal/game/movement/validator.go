package movement

// Uncomment when implementing:
// import (
// 	"time"

// 	"realm-server/internal/entity"
// 	"realm-server/pkg/math"
// )

// =============================================================================
// MOVEMENT VALIDATION (Anti-Cheat)
// =============================================================================
//
// WoW uses client-authoritative movement for responsiveness.
// The server validates but doesn't simulate movement.
//
// Validation catches:
// - Speed hacks (moving faster than allowed)
// - Teleport hacks (instant position changes)
// - Fly hacks (flying without ability)
// - Wall climbing (passing through collision)
//
// Tradeoffs:
// - Too strict: Legitimate players get false positives (lag, rubberbanding)
// - Too lenient: Cheaters exploit the game
//
// Typical approach: Allow some slack, flag suspicious patterns, review/ban.

// Validator checks movement packets for cheating.
type Validator struct {
	// Tolerance for network latency
	// LatencyTolerance time.Duration // e.g., 500ms

	// Speed tolerance (multiplier over expected speed)
	SpeedTolerance float32 // e.g., 1.5 = allow 50% over max speed
}

// MovementCheck represents results of validation.
type MovementCheck struct {
	Valid      bool
	Suspicious bool   // Valid but questionable
	Reason     string // If invalid or suspicious
	Corrected  bool   // If position was corrected
	// NewPosition math.Vec3 // Corrected position if applicable
}

// TODO: Implement Validator:
//
// func NewValidator() *Validator
//
// func (v *Validator) ValidateMovement(player *entity.Player, newPos math.Vec3, newFlags entity.MoveFlags, clientTime uint32) *MovementCheck
//   Steps:
//   1. Calculate time delta (server time vs client time, accounting for latency)
//   2. Calculate distance from last position
//   3. Determine max allowed speed for current state
//   4. Check if distance / time <= max speed * tolerance
//   5. Check special cases (teleport, flying, swimming)
//   6. Return result

// =============================================================================
// SPEED CALCULATION
// =============================================================================

// Base speeds (units per second)
const (
	SpeedWalk    float32 = 2.5
	SpeedRun     float32 = 7.0
	SpeedSwim    float32 = 4.7
	SpeedFly     float32 = 7.0
	SpeedFalling float32 = 60.0 // Terminal velocity
)

// TODO: Implement speed calculation:
//
// func (v *Validator) GetMaxSpeed(player *entity.Player) float32
//   Base speed depends on:
//   - Movement type (walk, run, swim, fly)
//   - Mounted? What mount speed?
//   - Speed buffs (sprint, aspect of the cheetah, etc.)
//   - Speed debuffs (dazed, slowed, etc.)
//
//   return baseSpeed * mountMultiplier * buffMultiplier * debuffMultiplier

// =============================================================================
// SPECIAL VALIDATIONS
// =============================================================================

// TODO: Implement special checks:
//
// func (v *Validator) ValidateFlight(player *entity.Player, pos math.Vec3) bool
//   - Player must have flying ability or be on flying mount
//   - Must be in flyable zone
//   - Must not be in combat (optional rule)
//
// func (v *Validator) ValidateJump(player *entity.Player, jumpVelocity float32) bool
//   - Jump velocity should match physics
//   - Can't jump while rooted/stunned
//
// func (v *Validator) ValidateTeleport(player *entity.Player, oldPos, newPos math.Vec3) bool
//   - Large instant position change
//   - Only valid if: mage blink, warlock teleport, hearthstone, etc.
//   - Check if player cast such ability recently

// =============================================================================
// COLLISION (Optional - complex to implement)
// =============================================================================

// Full collision detection requires:
// - Loading map geometry (terrain heightmaps, building collision)
// - Pathfinding data
// - This is complex and often skipped in private servers
//
// Simpler approach: Trust client for now, ban obvious exploiters manually

// TODO: If implementing collision:
//
// func (v *Validator) CheckCollision(from, to math.Vec3) bool
//   - Raycast from old to new position
//   - Check for terrain/wall intersection
//
// func (v *Validator) GetGroundHeight(pos math.Vec3) float32
//   - Sample terrain heightmap
//   - Check if player is floating (fly hack without flying)
