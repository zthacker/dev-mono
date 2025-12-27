package net

import (
	"errors"
	"log"
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
	AccountID     uint32
	AccountName   string
	Permissions   uint32 // Bitflags for GM commands, etc.
	Authenticated bool

	// Current player (nil until CMSG_PLAYER_LOGIN)
	// Note: Use interface{} to avoid circular import. Cast to *entity.Player when using.
	Player interface{}

	// Timing
	LastActivity time.Time
	Latency      time.Duration // Measured from ping/pong

	// Rate limiting
	packetCount uint32
	lastReset   time.Time
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn:         conn,
		sendChan:     make(chan *Packet, 256),
		LastActivity: time.Now(),
	}
}

func (s *Session) Run() {
	// start send thread
	go s.sendLoop()

	for {
		pkt, err := ReadPacketFromConn(s.conn)
		if err != nil {
			break
		}
		s.LastActivity = time.Now()
		DispatchPacket(s, pkt)
	}

	s.Close()
}

func (s *Session) Send(pkt *Packet) {
	s.closeMu.Lock()
	defer s.closeMu.Unlock()

	if s.closed {
		return
	}

	s.sendChan <- pkt
}

func (s *Session) sendLoop() {
	for pkt := range s.sendChan {
		WritePacketToConn(s.conn, pkt)
	}
}

func (s *Session) Close() {
	s.closeMu.Lock()
	defer s.closeMu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	s.conn.Close()
	close(s.sendChan)

	if s.Player != nil {
		s.Player = nil
	}
}

func (s *Session) IsAlive() bool {
	return !s.closed && time.Since(s.LastActivity) < 30*time.Second
}

// =============================================================================
// PACKET HANDLING
// =============================================================================

// PacketHandler is the signature for opcode handlers.
type PacketHandler func(s *Session, pkt *Packet) error

//HandlerRegistry maps opcodes to handlers.

var handlers = map[Opcode]PacketHandler{
	CMSG_PING:           handlePing,
	CMSG_MOVE_HEARTBEAT: handleMoveHeartbeat,
}

func handlePing(s *Session, pkt *Packet) error {
	reader := NewPacketReader(pkt.Payload)
	seq, err := reader.ReadUint32()
	if err != nil {
		return err
	}

	writer := NewPacketWriter(SMSG_PONG)
	writer.WriteUint32(seq)
	s.Send(writer.Finish())

	return nil
}

func handleMoveHeartbeat(s *Session, pkt *Packet) error {
	s.LastActivity = time.Now()

	return nil
}

func DispatchPacket(s *Session, pkt *Packet) error {
	handler, ok := handlers[pkt.Opcode]
	if !ok {
		log.Printf("unknown opcode: %d", pkt.Opcode)
		return errors.New("unknown opcode")
	}

	return handler(s, pkt)
}

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
