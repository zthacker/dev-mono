package net

import (
	"net"
	"sync"
	"time"
)

// Session represents a connected client.
// One session per TCP connection, may or may not have a player attached.
//
// Lifecycle:
// 1. TCP accept -> new Session
// 2. Auth handshake -> verified account
// 3. Character select -> Player attached
// 4. In-game play
// 5. Logout or disconnect -> cleanup
type Session struct {
	// Network
	conn     net.Conn
	sendChan chan *Packet // Buffered channel for outgoing packets
	closed   bool
	closeMu  sync.Mutex

	// Authentication
	AccountID   uint32
	AccountName string
	Permissions uint32 // Bitflags for GM commands, etc.
	Authenticated bool

	// Current player (nil until CMSG_PLAYER_LOGIN)
	// Note: Use interface{} to avoid circular import. Cast to *entity.Player when using.
	Player interface{}

	// Timing
	LastActivity time.Time
	Latency      time.Duration // Measured from ping/pong

	// Rate limiting
	// packetCount uint32
	// lastReset   time.Time
}

// TODO: Implement Session:
//
// func NewSession(conn net.Conn) *Session
//   - Initialize send channel (buffer size ~64-256)
//   - Set LastActivity
//   - Start send goroutine
//
// func (s *Session) Run()
//   - Main receive loop
//   - Read packets from conn
//   - Dispatch to handler
//   - Handle disconnect
//
// func (s *Session) Send(pkt *Packet)
//   - Queue packet on sendChan
//   - Non-blocking (drop if full? or disconnect?)
//
// func (s *Session) sendLoop()
//   - Goroutine that writes from sendChan to conn
//   - Batches multiple small packets together
//
// func (s *Session) Close()
//   - Mark closed
//   - Close connection
//   - Cleanup player if attached
//
// func (s *Session) IsAlive() bool
//   - Check if session is still valid

// =============================================================================
// PACKET HANDLING
// =============================================================================

// PacketHandler is the signature for opcode handlers.
type PacketHandler func(s *Session, pkt *Packet) error

// HandlerRegistry maps opcodes to handlers.
// TODO: Implement as a simple map or switch statement:
//
// var handlers = map[Opcode]PacketHandler{
//     CMSG_PING: handlePing,
//     CMSG_MOVE_HEARTBEAT: handleMoveHeartbeat,
//     ...
// }
//
// func DispatchPacket(s *Session, pkt *Packet) error
//   - Look up handler
//   - Call it
//   - Log unknown opcodes

// =============================================================================
// RATE LIMITING
// =============================================================================

// TODO: Implement rate limiting to prevent packet spam:
//
// func (s *Session) CheckRateLimit() bool
//   - Count packets per second
//   - Disconnect if threshold exceeded
//   - Different limits for different opcode categories
//     (movement can be frequent, chat should be limited)

// =============================================================================
// ENCRYPTION (Optional but recommended)
// =============================================================================

// WoW uses RC4 encryption after the auth handshake.
// The session key is derived during SRP6 authentication.
//
// TODO: If you want encryption:
//
// type PacketCrypto struct {
//     sendCipher *rc4.Cipher
//     recvCipher *rc4.Cipher
// }
//
// func (s *Session) EnableEncryption(sessionKey []byte)
// func (s *Session) encryptHeader(header []byte)
// func (s *Session) decryptHeader(header []byte)
