## Running Benchmarks

### Basic Benchmark Commands

```bash
# Navigate to plugins directory
cd tera/plugins

# Run all benchmarks with memory stats (2 second duration)
go test -bench=. -benchmem -benchtime=2s

# Run all benchmarks with longer duration for more stable results
go test -bench=. -benchmem -benchtime=5s

# Run specific benchmark
go test -bench=BenchmarkPipeline_SmallPacket -benchmem

# Run only WASM-related benchmarks
go test -bench=BenchmarkWASM -benchmem

# Compare performance across packet sizes
go test -bench=BenchmarkPipeline -benchmem
```

### Advanced Profiling

```bash
# CPU profiling
go test -bench=BenchmarkPipeline_SmallPacket -cpuprofile=cpu.out
go tool pprof cpu.out

# Memory profiling
go test -bench=BenchmarkPipeline_SmallPacket -memprofile=mem.out
go tool pprof mem.out

# Allocation profiling
go test -bench=BenchmarkPipeline_SmallPacket -benchmem -memprofile=mem.out
go tool pprof -alloc_space mem.out

# Generate CPU profile graph (requires graphviz)
go test -bench=BenchmarkPipeline_SmallPacket -cpuprofile=cpu.out
go tool pprof -http=:8080 cpu.out
```

### Benchmark Comparison

```bash
# Save baseline results
go test -bench=. -benchmem > baseline.txt

# After making changes, compare
go test -bench=. -benchmem > optimized.txt

# Use benchstat for statistical comparison (install: go install golang.org/x/perf/cmd/benchstat@latest)
benchstat baseline.txt optimized.txt
```