package main

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var connectedCount int64

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	numClients := 32

	// Stagger connections to avoid overwhelming the server
	for i := 0; i < numClients; i++ {
		go fakeClient(mainCtx, i)
		time.Sleep(time.Millisecond) // 1ms between each connection
	}

	log.Printf("Launched %d clients, connected: %d", numClients, atomic.LoadInt64(&connectedCount))

	// Periodically log connection count
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-mainCtx.Done():
				return
			case <-ticker.C:
				log.Printf("Connected clients: %d", atomic.LoadInt64(&connectedCount))
			}
		}
	}()

	<-mainCtx.Done()
	log.Println("Received termination signal, shutting down...")
}

func fakeClient(ctx context.Context, id int) {
	conn, err := net.Dial("tcp", "localhost:8085")
	if err != nil {
		log.Printf("[Client %d] Failed to connect: %v", id, err)
		return
	}
	atomic.AddInt64(&connectedCount, 1)
	defer func() {
		atomic.AddInt64(&connectedCount, -1)
		conn.Close()
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var seqCount uint32
	var totalLatency time.Duration

	for {
		select {
		case <-ctx.Done():
			if seqCount > 0 {
				avgLatency := totalLatency / time.Duration(seqCount)
				log.Printf("[Client %d] Avg latency: %v over %d pings", id, avgLatency, seqCount)
			}
			return
		case <-ticker.C:
			start := time.Now()
			sendPing(conn, seqCount)
			readPong(conn)
			latency := time.Since(start)
			totalLatency += latency
			seqCount++
		}
	}
}

func sendPing(conn net.Conn, seq uint32) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], 6)         // size = 2 + 4
	binary.LittleEndian.PutUint16(buf[2:4], 0x0002) // CMSG_PING
	binary.LittleEndian.PutUint32(buf[4:8], seq)
	conn.Write(buf)
}

func readPong(conn net.Conn) {
	header := make([]byte, 4)
	conn.Read(header)
	size := binary.BigEndian.Uint16(header[0:2])

	payload := make([]byte, size-2)
	conn.Read(payload)
}
