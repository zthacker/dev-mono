<!-- # Session System -->

## Overview

A Session represents one connected client. The Server creates a Session for each TCP connection.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         SESSION                                 │
│                                                                 │
│   ┌─────────────┐                       ┌─────────────┐        │
│   │  Run()      │                       │ sendLoop()  │        │
│   │  (receive)  │                       │  (send)     │        │
│   │             │                       │             │        │
│   │ Read pkts   │                       │ Write pkts  │        │
│   │ from conn   │                       │ to conn     │        │
│   │      ↓      │                       │      ↑      │        │
│   │ Dispatch to │      sendChan        │ Read from   │        │
│   │ handlers    │ ──────────────────→  │ channel     │        │
│   │             │                       │             │        │
│   └─────────────┘                       └─────────────┘        │
│                                                                 │
│         TCP conn (reads)                  TCP conn (writes)     │
└─────────────────────────────────────────────────────────────────┘
```

## Two Goroutines Per Session

**Run() - Receive loop** (main session goroutine):
- Reads packets from TCP connection
- Dispatches to handlers based on opcode
- Blocks on `ReadPacketFromConn()`
- Calls `Close()` on disconnect

**sendLoop() - Send loop** (spawned by Run):
- Reads packets from `sendChan`
- Writes them to TCP connection
- Exits when `sendChan` is closed

## Why Separate Send/Receive?

1. **Non-blocking sends**: Handlers call `session.Send(pkt)` without waiting for TCP write
2. **Buffering**: `sendChan` queues packets if client is slow
3. **No blocking**: Slow clients don't block game logic
4. **Clean shutdown**: Close channel to stop send loop

## Lifecycle

```
1. TCP Accept
      │
      ▼
2. NewSession(conn)
      │
      ▼
3. go session.Run()
      │
      ├──→ spawns sendLoop()
      │
      ▼
4. Read packets in loop
      │
      ├──→ DispatchPacket() for each
      │
      ▼
5. Error or disconnect
      │
      ▼
6. session.Close()
      │
      ├──→ closes conn
      ├──→ closes sendChan (stops sendLoop)
      └──→ cleans up Player
```

## Packet Flow

**Incoming (client → server):**
```
TCP conn → ReadPacketFromConn() → DispatchPacket() → handler function
```

**Outgoing (server → client):**
```
handler calls Send() → sendChan → sendLoop() → WritePacketToConn() → TCP conn
```

## Thread Safety

- `closeMu` mutex protects `closed` flag
- `Send()` checks `closed` before writing to channel (prevents panic)
- `Close()` is idempotent (safe to call multiple times)

## Handler Pattern

```go
var handlers = map[Opcode]PacketHandler{
    CMSG_PING:           handlePing,
    CMSG_MOVE_HEARTBEAT: handleMoveHeartbeat,
}

func handlePing(s *Session, pkt *Packet) error {
    // For simple packets, read directly from payload
    if len(pkt.Payload) < 4 {
        return errors.New("packet too short")
    }
    seq := binary.LittleEndian.Uint32(pkt.Payload[:4])

    // Send response
    writer := NewPacketWriter(SMSG_PONG)
    writer.WriteUint32(seq)
    s.Send(writer.Finish())

    return nil
}
```

## Memory Optimization

For simple handlers (ping, heartbeat), read directly from `pkt.Payload` instead of creating a `PacketReader`:

```go
// Simple - read directly
seq := binary.LittleEndian.Uint32(pkt.Payload[:4])

// Complex packets - use reader for convenience
reader := NewPacketReader(pkt.Payload)
x, _ := reader.ReadFloat32()
y, _ := reader.ReadFloat32()
z, _ := reader.ReadFloat32()
name, _ := reader.ReadString()
```

Use `PacketReader`/`PacketWriter` for complex packets with many fields. The small struct allocation is negligible compared to packet payload allocations.