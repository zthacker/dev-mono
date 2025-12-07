package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"

	"example_automation/pkg/backend/mock"
	"example_automation/pkg/modules/data"
	"example_automation/pkg/modules/ground"
	"example_automation/pkg/modules/satellite"
	"example_automation/pkg/modules/system"
	"example_automation/pkg/registry"
)

// lock for printing
var printLock sync.Mutex

func safePrint(threadName string, msg string) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Printf("[%s] %s\n", threadName, msg)
}

func runScript(threadName string, sourceCode []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	// Build globals from registry
	globals := registry.DefaultRegistry.BuildGlobals()

	// Inject introspection functions (help and dir)
	registry.InjectIntrospection(registry.DefaultRegistry, globals)

	// Create the Thread
	thread := &starlark.Thread{Name: threadName}

	// Hook up the print function to our safe logger
	thread.Print = func(_ *starlark.Thread, msg string) {
		safePrint(threadName, fmt.Sprintf("%s", msg))
	}

	// Set the safety limit
	thread.SetMaxExecutionSteps(1000)

	safePrint(threadName, "Starting...")

	// Execute
	_, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, threadName, sourceCode, globals)

	if err != nil {
		// Detect if it was a step limit error
		if evalErr, ok := err.(*starlark.EvalError); ok {
			if evalErr.Msg == "too many steps" {
				safePrint(threadName, "KILLED: Exceeded Step Limit!")
				return
			}
		}
		safePrint(threadName, fmt.Sprintf("ERROR: %v", err))
		return
	}

	safePrint(threadName, "Finished Successfully.")
}

// main engine
func main() {
	fmt.Println("=== Mission Control Starlark Automation System ===\n")

	// Create backend services
	telemetryService := mock.NewMockTelemetryService()
	commandService := mock.NewMockCommandService()
	groundService := mock.NewMockGroundStationService()
	dataService := mock.NewMockDataProcessingService()

	// Create and register modules
	fmt.Println("Registering modules...")
	if err := registry.DefaultRegistry.Register(system.NewSystemModule()); err != nil {
		log.Fatalf("Failed to register system module: %v", err)
	}
	if err := registry.DefaultRegistry.Register(satellite.NewSatelliteModule(telemetryService, commandService)); err != nil {
		log.Fatalf("Failed to register satellite module: %v", err)
	}
	if err := registry.DefaultRegistry.Register(ground.NewGroundModule(groundService)); err != nil {
		log.Fatalf("Failed to register ground module: %v", err)
	}
	if err := registry.DefaultRegistry.Register(data.NewDataModule(dataService)); err != nil {
		log.Fatalf("Failed to register data module: %v", err)
	}

	// Display registered modules
	modules := registry.DefaultRegistry.All()
	fmt.Printf("Registered %d modules:\n", len(modules))
	for _, mod := range modules {
		meta := mod.Metadata()
		fmt.Printf("  - %s v%s (%s)\n", meta.Name, meta.Version, meta.Description)
	}
	fmt.Println()

	var wg sync.WaitGroup

	// This comes from a db query
	dbScripts := map[string]string{
		"Sat-Pass-101": `print("Running DB Script 101"); system.wait(2)`,
		"Sat-Pass-102": `print("Running DB Script 102"); system.wait(1)`,
	}

	fmt.Println("--- Launching Database Scripts ---")
	for id, code := range dbScripts {
		wg.Add(1)
		// We convert the string to []byte
		go runScript("DB:"+id, []byte(code), &wg)
	}

	// These are passed via CLI args: ./engine --test local_test.star
	localFiles := []string{
		"scripts/mission_star_pass.star",
		"scripts/test_satellite.star",
		"scripts/test_introspection.star",
	}

	fmt.Println("--- Launching Local Test Scripts ---")
	for _, filename := range localFiles {
		// Read the file into memory FIRST
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Could not read file %s: %v", filename, err)
			continue
		}

		wg.Add(1)
		// Pass the content to the same runner
		go runScript(filename, content, &wg)
	}

	wg.Wait()
	fmt.Println("\n--- All Operations Complete ---")
}
