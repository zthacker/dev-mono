package world

import (
	"time"
)

// TickLoop manages the game simulation loop.
// This is the heartbeat of the server - everything happens in ticks.
//
// Tick rate: 20 Hz (50ms) is common for MMOs
// - Fast enough for responsive combat
// - Slow enough to not overload the server
// - Matches WoW's internal tick rate
//
// Each tick processes:
// 1. Input: Pending player actions, NATS messages
// 2. Simulation: Entity updates, combat, AI
// 3. Output: State broadcasts to clients
type TickLoop struct {
	tickRate  int           // Ticks per second
	tickDelta time.Duration // Duration per tick
	ticker    *time.Ticker

	// Callbacks
	onTick func(tick uint64, dt float64)

	// Stats
	currentTick  uint64
	lastTickTime time.Time
	maxTickTime  time.Duration
	slowTicks    uint64 // Ticks that took longer than tickDelta
}

// TODO: Implement TickLoop:
//
// func NewTickLoop(tickRate int) *TickLoop
//   - Calculate tickDelta
//   - Create ticker
//
// func (t *TickLoop) SetOnTick(fn func(tick uint64, dt float64))
//
// func (t *TickLoop) Run(stopChan <-chan struct{})
//   for {
//       select {
//       case <-stopChan:
//           return
//       case <-t.ticker.C:
//           t.processTick()
//       }
//   }
//
// func (t *TickLoop) processTick()
//   - Record start time
//   - Calculate actual dt since last tick
//   - Call onTick callback
//   - Record end time
//   - Log warning if tick took too long
//   - Update stats

// =============================================================================
// TICK PHASES
// =============================================================================

// The server's onTick should do these phases in order:

// Phase 1: Process Input
// - Drain queued player actions (movement, spells, etc.)
// - Process NATS messages from other shards
// - These are "instantaneous" - they set up state for simulation

// Phase 2: Simulate
// - Update all entity positions (apply velocity)
// - Process combat (damage, healing, effects)
// - Run NPC AI
// - Check triggers (quest objectives, zone events)
// - Process spell effects

// Phase 3: Synchronize
// - Determine what changed this tick
// - Build update packets for each player
// - Send only relevant updates (AOI filtering)

// =============================================================================
// EXAMPLE TICK IMPLEMENTATION
// =============================================================================

// func (s *Server) onTick(tick uint64, dt float64) {
//     // Phase 1: Input
//     s.processQueuedActions()
//     s.processNATSMessages()
//
//     // Phase 2: Simulate
//     s.entityMgr.UpdateAll(dt)
//     s.zoneMgr.UpdateAll(dt)
//
//     // Phase 3: Synchronize
//     s.broadcastUpdates()
//
//     // Periodic tasks (not every tick)
//     if tick % 20 == 0 { // Every second
//         s.updatePlayerLatencies()
//     }
//     if tick % 200 == 0 { // Every 10 seconds
//         s.persistHotState()
//     }
// }
