package entity

import (
	"sync"
	"time"
)

// BaseEntity provides common functionality for all entities.
// Embed this in Player, NPC, etc. to get the base implementation.
//
// Thread Safety:
// - All position/movement reads and writes go through the mutex
// - The actor model (one goroutine per entity) reduces contention
// - But the world may read positions from any goroutine for AOI queries
type BaseEntity struct {
	id       EntityID
	zoneID   uint32
	movement MovementState
	mu       sync.RWMutex
}

func NewBaseEntity(id EntityID) *BaseEntity {
	return &BaseEntity{
		id: id,
		movement: MovementState{
			LastUpdate: time.Now(),
		},
	}
}

func (e *BaseEntity) ID() EntityID {
	return e.id
}

func (e *BaseEntity) Type() EntityType {
	return e.id.Type()
}

func (e *BaseEntity) Transform() Transform {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.movement.Transform
}

func (e *BaseEntity) SetTransform(t Transform) {
	//   - Lock, set position, update LastUpdate, Unlock
	e.mu.Lock()
	defer e.mu.Unlock()
	e.movement.Transform = t
	e.movement.LastUpdate = time.Now()
}

// Caller must handle locking if modifying
func (e *BaseEntity) Movement() *MovementState {
	//   Returns pointer to movement state
	return &e.movement
}

func (e *BaseEntity) ZoneID() uint32 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.zoneID
}

func (e *BaseEntity) SetZoneID(id uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.zoneID = id
}

// Optional lock helpers:
func (e *BaseEntity) Lock()    { e.mu.Lock() }
func (e *BaseEntity) Unlock()  { e.mu.Unlock() }
func (e *BaseEntity) RLock()   { e.mu.RLock() }
func (e *BaseEntity) RUnlock() { e.mu.RUnlock() }
