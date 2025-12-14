package protocols

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
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

	// Read source
	wasmBytes, err := os.ReadFile(wasmFile)
	if err != nil {
		return nil, err
	}

	// Instantiate
	mod, err := r.Instantiate(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	malloc := mod.ExportedFunction("allocate_buffer")

	// Reserve 1MB (1024 * 1024)
	bufferSize := uint64(1 * 1024 * 1024)

	results, err := malloc.Call(ctx, bufferSize)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate shared buffer: %v", err)
	}

	sharedPtr := results[0] // This is our permanent "PO Box"

	return &WASMStep{
		Module:      mod,
		ProcessFunc: mod.ExportedFunction("process_packet"),
		SharedPtr:   sharedPtr,
		BufferSize:  bufferSize,
	}, nil
}

func (w *WASMStep) Name() string { return "User-WASM-C++" }

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
	copy(finalData, data)

	return finalData, nil

}
