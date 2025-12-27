package net

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// Packet represents a network message.
// Wire format:
//
//	[2 bytes] size (big-endian, includes opcode)
//	[2 bytes] opcode (little-endian)
//	[N bytes] payload
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
	opcode Opcode
	buf    []byte
}

// =============================================================================
// PACKET READER
// =============================================================================

func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{
		data: data,
		pos:  0,
	}
}

func (r *PacketReader) ReadUint8() (uint8, error) {
	if err := r.checkBounds(1); err != nil {
		return 0, err
	}

	val := r.data[r.pos]

	r.pos++

	return val, nil
}

func (r *PacketReader) ReadUint16() (uint16, error) {
	if err := r.checkBounds(2); err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint16(r.data[r.pos:])
	r.pos += 2

	return val, nil
}

func (r *PacketReader) ReadUint32() (uint32, error) {
	if err := r.checkBounds(4); err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4

	return val, nil
}

func (r *PacketReader) ReadUint64() (uint64, error) {
	if err := r.checkBounds(8); err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint64(r.data[r.pos:])
	r.pos += 8

	return val, nil
}

func (r *PacketReader) ReadInt32() (int32, error) {
	if err := r.checkBounds(4); err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4

	return int32(val), nil
}

func (r *PacketReader) ReadFloat32() (float32, error) {
	if err := r.checkBounds(4); err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4

	return math.Float32frombits(val), nil
}

func (r *PacketReader) ReadString() (string, error) {
	// null term check
	for i := r.pos; i < len(r.data); i++ {
		if r.data[i] == 0x00 {
			str := string(r.data[r.pos:i])
			r.pos = i + 1 // move past null term
			return str, nil
		}
	}

	return "", errors.New("string not null-terminated")
}

func (r *PacketReader) ReadBytes(n int) ([]byte, error) {
	if err := r.checkBounds(n); err != nil {
		return nil, err
	}

	val := r.data[r.pos : r.pos+n]
	r.pos += n
	return val, nil
}

func (r *PacketReader) Remaining() int {
	return len(r.data) - r.pos
}

func (r *PacketReader) checkBounds(bounds int) error {
	if r.pos+bounds > len(r.data) {
		return errors.New("read past end of packet")
	}

	return nil
}

// =============================================================================
// PACKET WRITER
// =============================================================================

func NewPacketWriter(opcode Opcode) *PacketWriter {
	return &PacketWriter{
		opcode: opcode,
		buf:    make([]byte, 0, 64),
	}
}

func (w *PacketWriter) WriteUint8(v uint8) {
	w.buf = append(w.buf, v)
}
func (w *PacketWriter) WriteUint16(v uint16) {
	w.buf = append(w.buf, 0, 0) // grow by 2
	binary.LittleEndian.PutUint16(w.buf[len(w.buf)-2:], v)
}

func (w *PacketWriter) WriteUint32(v uint32) {
	w.buf = append(w.buf, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(w.buf[len(w.buf)-4:], v)
}

func (w *PacketWriter) WriteUint64(v uint64) {
	w.buf = append(w.buf, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.LittleEndian.PutUint64(w.buf[len(w.buf)-8:], v)
}

func (w *PacketWriter) WriteInt32(v int32) {
	w.buf = append(w.buf, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(w.buf[len(w.buf)-4:], uint32(v))
}

func (w *PacketWriter) WriteFloat32(v float32) {
	w.buf = append(w.buf, 0, 0, 0, 0)
	binary.LittleEndian.PutUint32(w.buf[len(w.buf)-4:], math.Float32bits(v))
}

func (w *PacketWriter) WriteString(s string) {
	w.buf = append(w.buf, s...)
	w.buf = append(w.buf, 0x00)
}

func (w *PacketWriter) WriteBytes(b []byte) {
	w.buf = append(w.buf, b...)
}

func (w *PacketWriter) Finish() *Packet {
	return &Packet{
		Opcode:  w.opcode,
		Payload: w.buf,
	}
}

// =============================================================================
// WIRE FORMAT
// =============================================================================

func ReadPacketFromConn(r io.Reader) (*Packet, error) {
	// using ReadFull to handle TCP partials
	//Read 2-byte size header
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint16(header)
	if size < 2 {
		return nil, errors.New("packet too small")
	}

	//Read opcode + payload
	body := make([]byte, size)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}

	opcode := Opcode(binary.LittleEndian.Uint16(body[:2]))
	payload := body[2:]

	return &Packet{
		Opcode:  opcode,
		Payload: payload,
	}, nil

}

func WritePacketToConn(w io.Writer, pkt *Packet) error {
	// Size = opcode (2) + payload length
	size := uint16(2 + len(pkt.Payload))

	// Build header: size (big-endian) + opcode (little-endian)
	header := make([]byte, 4)
	binary.BigEndian.PutUint16(header[:2], size)
	binary.LittleEndian.PutUint16(header[2:], uint16(pkt.Opcode))

	// Write header
	if _, err := w.Write(header); err != nil {
		return err
	}

	// Write payload
	if _, err := w.Write(pkt.Payload); err != nil {
		return err
	}

	return nil
}

// TODO:
// - Packet encryption (after auth handshake)
// - Compression for large packets
// - Packet pooling to reduce allocations
