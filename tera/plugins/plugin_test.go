package plugins

import (
	"context"
	"tera/internal/pipeline"
	"tera/internal/protocols"
	"testing"
)

func TestWASM(t *testing.T) {
	steps := []protocols.PipelineStep{}

	steps = append(steps, &protocols.StripProcotol{})

	// wasm plugin
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wasmPlugin, err := protocols.NewWASMStep(testCtx, "plugin.wasm")
	if err != nil {
		t.Fatalf("error loading wasm plugin: %s", err)
	}

	steps = append(steps, wasmPlugin)

	data := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	pipeline.RunPipeline(steps, data)

}
