package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

func getTLM(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mnemonic string
	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 1, &mnemonic); err != nil {
		return nil, err
	}

	// sim som telem data -- but we'd get it from somewhere real in a real system
	fmt.Printf("Getting TLM for %s\n", mnemonic)

	// float value
	return starlark.Float(28.5), nil
}

func sendCMD(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mnemonic string
	var params *starlark.Dict

	// unpack 1 string (cmdName), 1 dict (params)
	if err := starlark.UnpackPositionalArgs(builtin.Name(), args, kwargs, 2, &mnemonic, &params); err != nil {
		return nil, err
	}

	// iterate over params dict
	fmt.Printf("Sending CMD %s with params:\n", mnemonic)
	for _, key := range params.Keys() {
		value, _, _ := params.Get(key)
		fmt.Printf("  %s: %s\n", key.String(), value.String())
	}

	// sim sending command -- but we'd do something real in a real system
	return starlark.None, nil
}

// lock for printing
var printLock sync.Mutex

func safePrint(threadName string, msg string) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Printf("[%s] %s\n", threadName, msg)
}

func sysWait(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var seconds int
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &seconds); err != nil {
		return nil, err
	}

	fmt.Printf("[Go-Engine] Sleeping for %d seconds...\n", seconds)
	// in the real world, might use time.Sleep, or handle a select{} on a context
	time.Sleep(time.Duration(seconds) * time.Second)

	return starlark.None, nil
}

func runScript(threadName string, sourceCode []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	// Globals (modules)
	// Note: We create these inside the function or reuse them if they are stateless.
	// starlark.NewBuiltin is immutable/thread-safe.
	systemMembers := starlark.StringDict{
		"wait": starlark.NewBuiltin("wait", sysWait),
	}
	sysModule := starlarkstruct.FromStringDict(starlark.String("system"), systemMembers)

	globals := starlark.StringDict{
		"system": sysModule,
	}

	// Create the Thread
	thread := &starlark.Thread{Name: threadName}

	// Hook up the print function to our safe logger
	thread.Print = func(_ *starlark.Thread, msg string) {
		safePrint(threadName, fmt.Sprintf("%s", msg))
	}

	// Set the safety limit
	thread.SetMaxExecutionSteps(1000)

	safePrint(threadName, fmt.Sprintf("[%s] Starting..."))

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

	safePrint(threadName, fmt.Sprintf("Finished Successfully."))
}

func runLocalFile(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 1. Setup Environment (Modules)
	sysModule := starlarkstruct.FromStringDict(starlark.String("system"), starlark.StringDict{
		"wait": starlark.NewBuiltin("wait", sysWait),
	})

	globals := starlark.StringDict{
		"system": sysModule,
	}

	// 2. Create Thread
	thread := &starlark.Thread{Name: filePath}
	thread.Print = func(_ *starlark.Thread, msg string) {
		safePrint(filePath, msg)
	}

	// 3. Set Safety Limit (e.g. 2000 steps)
	thread.SetMaxExecutionSteps(2000)

	safePrint(filePath, "STARTING")

	// 4. EXECUTE LOCAL FILE
	// Notice: we pass 'nil' as the 4th argument.
	// This tells Starlark: "Read the source code from the 'filePath' on disk."
	_, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, filePath, nil, globals)

	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			if evalErr.Msg == "too many steps" {
				safePrint(filePath, "KILLED: Infinite Loop detected")
				return
			}
		}
		safePrint(filePath, fmt.Sprintf("ERROR: %v", err))
		return
	}

	safePrint(filePath, "COMPLETE")
}

// main engine
func main() {
	var wg sync.WaitGroup

	// --- CASE 1: PRODUCTION (From Database) ---
	// Imagine this comes from your ClickHouse/Postgres query
	dbScripts := map[string]string{
		"Sat-Pass-101": `print("Running DB Script 101"); system.wait(5)`,
		"Sat-Pass-102": `print("Running DB Script 102"); system.wait(3)`,
	}

	fmt.Println("--- 🚀 Launching Database Scripts ---")
	for id, code := range dbScripts {
		wg.Add(1)
		// We convert the string to []byte
		go runScript("DB:"+id, []byte(code), &wg)
	}

	// --- CASE 2: DEVELOPMENT (From Local Disk) ---
	// Imagine these are passed via CLI args: ./engine --test local_test.star
	localFiles := []string{"mission_star_pass.star"}

	fmt.Println("--- 🛠️ Launching Local Test Scripts ---")
	for _, filename := range localFiles {
		// 1. Read the file into memory FIRST
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Could not read file %s: %v", filename, err)
			continue
		}

		wg.Add(1)
		// 2. Pass the content to the same runner
		go runScript(filename, content, &wg)
	}

	wg.Wait()
	fmt.Println("--- All Operations Complete ---")
}
