package protocols

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WASMStep struct {
	Module      api.Module
	ProcessFunc api.Function

	SharedPtr  uint64
	BufferSize uint64
}

func NewWASMStep(ctx context.Context, wasmFile string) (*WASMStep, error) {
	// Initialize Runtime
	r := wazero.NewRuntime(ctx)

	// Enable WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Read source
	wasmBytes, err := os.ReadFile(wasmFile)
	if err != nil {
		return nil, err
	}

	// Compile the module
	compiledMod, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	// Instantiate without calling _start
	modConfig := wazero.NewModuleConfig().WithName("").WithStartFunctions()
	mod, err := r.InstantiateModule(ctx, compiledMod, modConfig)
	if err != nil {
		return nil, err
	}

	// Call _initialize if it exists (for WASI reactor modules)
	if initFn := mod.ExportedFunction("_initialize"); initFn != nil {
		if _, err := initFn.Call(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize WASM module: %v", err)
		}
	}

	malloc := mod.ExportedFunction("allocate_buffer")

	// Reserve 1MB (1024 * 1024)
	bufferSize := uint64(1 * 1024 * 1024)

	results, err := malloc.Call(ctx, bufferSize)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate shared buffer: %v", err)
	}

	sharedPtr := results[0]

	return &WASMStep{
		Module:      mod,
		ProcessFunc: mod.ExportedFunction("process_packet"),
		SharedPtr:   sharedPtr,
		BufferSize:  bufferSize,
	}, nil
}

func (w *WASMStep) Name() string { return "User-WASM" }

func (w *WASMStep) Process(ctx context.Context, data []byte) ([]byte, error) {

	if uint64(len(data)) > w.BufferSize {
		return nil, fmt.Errorf("packet too large")
	}

	if !w.Module.Memory().Write(uint32(w.SharedPtr), data) {
		return nil, fmt.Errorf("failed to write to shared buffer")

	}

	ret, err := w.ProcessFunc.Call(ctx, w.SharedPtr, uint64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to call process_packet: %v", err)
	}

	// unpacket
	packed := ret[0]
	resPtr := uint32(packed >> 32)
	resLen := uint32(packed)

	// read back
	out, _ := w.Module.Memory().Read(resPtr, resLen)

	// find a better way than just do a make here -- not a fan
	finalData := make([]byte, len(out))
	copy(finalData, out)

	return finalData, nil

}
