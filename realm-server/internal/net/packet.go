package net

// Uncomment when implementing:
// import (
// 	"encoding/binary"
// 	"io"
// )

// Packet represents a network message.
// Wire format:
//   [2 bytes] size (big-endian, includes opcode)
//   [2 bytes] opcode (little-endian)
//   [N bytes] payload
//
// WoW uses a similar format but with encryption after auth.
// The size header allows reading exact packet boundaries on TCP.
type Packet struct {
	Opcode  Opcode
	Payload []byte
}

// PacketReader wraps a byte slice for reading structured data.
// All reads are little-endian (x86 native).
type PacketReader struct {
	data []byte
	pos  int
}

// PacketWriter builds a packet payload.
type PacketWriter struct {
	buf []byte
}

// =============================================================================
// PACKET READER
// =============================================================================

// TODO: Implement PacketReader:
//
// func NewPacketReader(data []byte) *PacketReader
//
// func (r *PacketReader) ReadUint8() (uint8, error)
// func (r *PacketReader) ReadUint16() (uint16, error)
// func (r *PacketReader) ReadUint32() (uint32, error)
// func (r *PacketReader) ReadUint64() (uint64, error)
// func (r *PacketReader) ReadInt32() (int32, error)
// func (r *PacketReader) ReadFloat32() (float32, error)
//
// func (r *PacketReader) ReadString() (string, error)
//   - Read null-terminated string
//
// func (r *PacketReader) ReadBytes(n int) ([]byte, error)
//
// func (r *PacketReader) Remaining() int
//   - Bytes left to read

// =============================================================================
// PACKET WRITER
// =============================================================================

// TODO: Implement PacketWriter:
//
// func NewPacketWriter(opcode Opcode) *PacketWriter
//   - Pre-allocate reasonable buffer
//
// func (w *PacketWriter) WriteUint8(v uint8)
// func (w *PacketWriter) WriteUint16(v uint16)
// func (w *PacketWriter) WriteUint32(v uint32)
// func (w *PacketWriter) WriteUint64(v uint64)
// func (w *PacketWriter) WriteInt32(v int32)
// func (w *PacketWriter) WriteFloat32(v float32)
//
// func (w *PacketWriter) WriteString(s string)
//   - Write null-terminated
//
// func (w *PacketWriter) WriteBytes(b []byte)
//
// func (w *PacketWriter) Finish() *Packet
//   - Return completed packet

// =============================================================================
// WIRE FORMAT
// =============================================================================

// TODO: Implement wire encoding/decoding:
//
// func ReadPacketFromConn(r io.Reader) (*Packet, error)
//   - Read 2-byte size header
//   - Read opcode + payload
//   - Handle partial reads (TCP)
//
// func WritePacketToConn(w io.Writer, pkt *Packet) error
//   - Write size header
//   - Write opcode
//   - Write payload
//
// Note: For production, you'd want:
// - Packet encryption (after auth handshake)
// - Compression for large packets
// - Packet pooling to reduce allocations
