package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// Uncomment when implementing:
	// "context"
	// "realm-server/pkg/config"
)

// =============================================================================
// WORLD SERVER ENTRY POINT
// =============================================================================
//
// This is the main game server process.
// Run one instance per shard/zone cluster.
//
// Startup sequence:
// 1. Load configuration
// 2. Connect to MySQL
// 3. Connect to NATS
// 4. Load zone data
// 5. Start game loop
// 6. Accept player connections
//
// Shutdown sequence:
// 1. Stop accepting new connections
// 2. Notify connected players
// 3. Save all player data
// 4. Persist hot state to NATS KV
// 5. Close connections

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting World Server...")

	// TODO: Implement startup:
	//
	// // 1. Load configuration
	// cfg := config.DefaultConfig()
	// // Or: cfg, err := config.Load()
	//
	// // 2. Connect to MySQL
	// db, err := mysql.NewRepository(cfg.Database.DSN())
	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer db.Close()
	//
	// // 3. Load templates into memory
	// log.Println("Loading game data...")
	// itemTemplates, _ := db.LoadItemTemplates(context.Background())
	// npcTemplates, _ := db.LoadNPCTemplates(context.Background())
	// spellTemplates, _ := db.LoadSpellTemplates(context.Background())
	// log.Printf("Loaded %d items, %d NPCs, %d spells",
	//     len(itemTemplates), len(npcTemplates), len(spellTemplates))
	//
	// // 4. Connect to NATS
	// nc, err := nats.Connect(cfg.NATS.URL)
	// if err != nil {
	//     log.Fatalf("Failed to connect to NATS: %v", err)
	// }
	// defer nc.Close()
	// js, _ := jetstream.New(nc)
	//
	// // 5. Initialize cache
	// cache, err := cache.NewCache(js)
	// if err != nil {
	//     log.Fatalf("Failed to initialize cache: %v", err)
	// }
	//
	// // 6. Create and start world server
	// server := world.NewServer(world.Config{
	//     BindAddr: cfg.Server.BindAddr,
	//     // ... other config
	// })
	//
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	//
	// go func() {
	//     if err := server.Run(ctx); err != nil {
	//         log.Printf("Server error: %v", err)
	//     }
	// }()
	//
	// log.Printf("World Server started on %s", cfg.Server.BindAddr)

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("Received signal %v, initiating shutdown...", sig)

	// TODO: Graceful shutdown:
	//
	// // Give players warning
	// server.BroadcastSystemMessage("Server shutting down in 10 seconds...")
	// time.Sleep(10 * time.Second)
	//
	// // Stop accepting new connections
	// cancel()
	//
	// // Wait for server to finish (with timeout)
	// shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer shutdownCancel()
	// server.Shutdown(shutdownCtx)

	log.Println("World Server stopped")
}
