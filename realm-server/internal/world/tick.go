package world

import (
	"log"
	"time"
)

// TickLoop manages the game simulation loop.
// This is the heartbeat of the server - everything happens in ticks.
//
// Tick rate: 20 Hz (50ms) is common for MMOs
// - Fast enough for responsive combat
// - Slow enough to not overload the server
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

// NewTickLoop creates a new tick loop with the given tick rate.
// tickRate is ticks per second (e.g., 20 for 50ms ticks).
func NewTickLoop(tickRate int) *TickLoop {
	tickDelta := time.Second / time.Duration(tickRate)
	return &TickLoop{
		tickRate:     tickRate,
		tickDelta:    tickDelta,
		ticker:       time.NewTicker(tickDelta),
		lastTickTime: time.Now(),
	}
}

// SetOnTick sets the callback function called each tick.
// tick is the current tick number (starts at 0).
// dt is seconds since last tick (typically 0.05 for 20 tick/sec).
func (t *TickLoop) SetOnTick(fn func(tick uint64, dt float64)) {
	t.onTick = fn
}

// Run starts the tick loop. Blocks until stopChan is closed.
func (t *TickLoop) Run(stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			t.ticker.Stop()
			return
		case <-t.ticker.C:
			t.processTick()
		}
	}
}

// processTick handles a single tick.
func (t *TickLoop) processTick() {
	startTime := time.Now()

	// Calculate actual dt (may vary slightly from tickDelta)
	dt := startTime.Sub(t.lastTickTime).Seconds()
	t.lastTickTime = startTime

	// Call the tick handler
	if t.onTick != nil {
		t.onTick(t.currentTick, dt)
	}

	t.currentTick++

	// Track timing stats
	elapsed := time.Since(startTime)
	if elapsed > t.maxTickTime {
		t.maxTickTime = elapsed
	}
	if elapsed > t.tickDelta {
		t.slowTicks++
		log.Printf("Tick %d took %v (budget: %v)", t.currentTick-1, elapsed, t.tickDelta)
	}
}

// Stats returns current tick loop statistics.
func (t *TickLoop) Stats() (currentTick uint64, maxTickTime time.Duration, slowTicks uint64) {
	return t.currentTick, t.maxTickTime, t.slowTicks
}

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

// Example onTick handler for the Server:
//
// func (s *Server) onTick(tick uint64, dt float64) {
// 	// Phase 1: Input
// 	s.processQueuedActions()
// 	s.processNATSMessages()
//
// 	// Phase 2: Simulate
// 	s.entityMgr.UpdateAll(dt)
// 	s.zoneMgr.UpdateAll(dt)
//
// 	// Phase 3: Synchronize
// 	s.broadcastUpdates()
//
// 	// Periodic tasks (not every tick)
// 	if tick%20 == 0 { // Every second
// 		s.updatePlayerLatencies()
// 	}
// 	if tick%200 == 0 { // Every 10 seconds
// 		s.persistHotState()
// 	}
// }
