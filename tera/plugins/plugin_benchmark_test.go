package plugins

import (
	"context"
	"tera/internal/protocols"
	"testing"
)

// BenchmarkPipeline_SmallPacket benchmarks the full pipeline with a 6-byte packet
func BenchmarkPipeline_SmallPacket(b *testing.B) {
	// Setup pipeline (same as test)
	steps := []protocols.PipelineStep{}
	steps = append(steps, &protocols.StripProcotol{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasmPlugin, err := protocols.NewWASMStep(ctx, "plugin.wasm")
	if err != nil {
		b.Fatalf("error loading wasm plugin: %s", err)
	}
	defer wasmPlugin.Module.Close(ctx)

	steps = append(steps, wasmPlugin)

	data := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Run pipeline (without printing)
		currentData := data
		for _, step := range steps {
			currentData, err = step.Process(ctx, currentData)
			if err != nil {
				b.Fatalf("pipeline failed: %v", err)
			}
		}
	}
}

// BenchmarkPipeline_MediumPacket benchmarks with 1KB packets
func BenchmarkPipeline_MediumPacket(b *testing.B) {
	// Setup pipeline
	steps := []protocols.PipelineStep{}
	steps = append(steps, &protocols.StripProcotol{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasmPlugin, err := protocols.NewWASMStep(ctx, "plugin.wasm")
	if err != nil {
		b.Fatalf("error loading wasm plugin: %s", err)
	}
	defer wasmPlugin.Module.Close(ctx)

	steps = append(steps, wasmPlugin)

	// 1KB packet
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		currentData := data
		for _, step := range steps {
			currentData, err = step.Process(ctx, currentData)
			if err != nil {
				b.Fatalf("pipeline failed: %v", err)
			}
		}
	}
}

// BenchmarkPipeline_LargePacket benchmarks with 64KB packets
func BenchmarkPipeline_LargePacket(b *testing.B) {
	// Setup pipeline
	steps := []protocols.PipelineStep{}
	steps = append(steps, &protocols.StripProcotol{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasmPlugin, err := protocols.NewWASMStep(ctx, "plugin.wasm")
	if err != nil {
		b.Fatalf("error loading wasm plugin: %s", err)
	}
	defer wasmPlugin.Module.Close(ctx)

	steps = append(steps, wasmPlugin)

	// 64KB packet
	data := make([]byte, 64*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		currentData := data
		for _, step := range steps {
			currentData, err = step.Process(ctx, currentData)
			if err != nil {
				b.Fatalf("pipeline failed: %v", err)
			}
		}
	}
}

// BenchmarkWASM_DirectCall isolates WASM call overhead
func BenchmarkWASM_DirectCall(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasmPlugin, err := protocols.NewWASMStep(ctx, "plugin.wasm")
	if err != nil {
		b.Fatalf("error loading wasm plugin: %s", err)
	}
	defer wasmPlugin.Module.Close(ctx)

	data := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := wasmPlugin.Process(ctx, data)
		if err != nil {
			b.Fatalf("wasm call failed: %v", err)
		}
	}
}

// BenchmarkPipeline_Batch simulates realistic FEP workload
func BenchmarkPipeline_Batch(b *testing.B) {
	// Process 1000 packets per iteration
	const batchSize = 1000

	// Setup pipeline
	steps := []protocols.PipelineStep{}
	steps = append(steps, &protocols.StripProcotol{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasmPlugin, err := protocols.NewWASMStep(ctx, "plugin.wasm")
	if err != nil {
		b.Fatalf("error loading wasm plugin: %s", err)
	}
	defer wasmPlugin.Module.Close(ctx)

	steps = append(steps, wasmPlugin)

	data := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < batchSize; j++ {
			currentData := data
			for _, step := range steps {
				currentData, err = step.Process(ctx, currentData)
				if err != nil {
					b.Fatalf("pipeline failed: %v", err)
				}
			}
		}
	}
}

// BenchmarkStripProtocol_Only benchmarks just the Strip protocol for baseline
func BenchmarkStripProtocol_Only(b *testing.B) {
	ctx := context.Background()
	strip := &protocols.StripProcotol{}
	data := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := strip.Process(ctx, data)
		if err != nil {
			b.Fatalf("strip failed: %v", err)
		}
	}
}
