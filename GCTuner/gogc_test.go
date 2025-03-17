package main

// This test suite contains our initial findings on benchmarking the results of statically toggling the GOGC and GOMEMLIMIT,
// which can show slight improvements on the AST Parsing Benchmark. In addition, it includes a vital function profileBenchmark
// which we use across all our benchmarks.

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"testing"
	"time"

	"github.com/kdhulipala41/golang-gc-263/GCTuner/astparse"
)

func BenchmarkASTParseDefault(b *testing.B) {
	profileBenchmark(b, "BenchmarkASTParseDefault", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCOff(b *testing.B) {
	debug.SetGCPercent(-1) // Not sure why, but this works over setting GOGC to a certain value.
	profileBenchmark(b, "BenchmarkASTParseGCOff", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCLow(b *testing.B) {
	debug.SetGCPercent(10)
	profileBenchmark(b, "BenchmarkASTParseGCLow", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCMid(b *testing.B) {
	debug.SetGCPercent(50)
	profileBenchmark(b, "BenchmarkASTParseGCMid", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCLarge(b *testing.B) {
	debug.SetGCPercent(200)
	profileBenchmark(b, "BenchmarkASTParseGCLarge", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseLowMemLimit(b *testing.B) {
	debug.SetGCPercent(-1)
	debug.SetMemoryLimit(10000 >> 20) // 10000MB
	profileBenchmark(b, "BenchmarkASTParseLowMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseMidMemLimit(b *testing.B) {
	debug.SetGCPercent(-1)
	debug.SetMemoryLimit(20000 >> 20) // 20000MB
	profileBenchmark(b, "BenchmarkASTParseMidMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseHighMemLimit(b *testing.B) {
	debug.SetGCPercent(-1)
	debug.SetMemoryLimit(40000 >> 20) // 40000MB
	profileBenchmark(b, "BenchmarkASTParseHighMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
}

func printGCCPUFrac() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("GC CPU Fraction: %.6f\n", memStats.GCCPUFraction)
}

func profileBenchmark(b *testing.B, name string, benchmarkFunc func()) {
	// Create directory for the benchmark
	dir := filepath.Join(".", "profiles", name)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		b.Fatalf("failed to create directory: %v", err)
	}

	// Start CPU profiling
	cpuProfile, err := os.Create(filepath.Join(dir, "cpu.prof"))
	if err != nil {
		b.Fatalf("could not create CPU profile: %v", err)
	}
	defer cpuProfile.Close()
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		b.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()
	defer printGCCPUFrac()

	// Channel to stop the goroutine
	stopChan := make(chan struct{})
	// Slice to store heap allocations
	var heapAllocs []uint64

	// Goroutine to query heap_alloc
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				heapAllocs = append(heapAllocs, memStats.HeapAlloc)
			case <-stopChan:
				return
			}
		}
	}()

	// Run the benchmark
	b.ResetTimer()
	benchmarkFunc()
	b.StopTimer()

	// Stop the goroutine
	close(stopChan)

	// Calculate average and p99
	var total uint64
	for _, alloc := range heapAllocs {
		total += alloc
	}
	avgHeapAlloc := total / uint64(len(heapAllocs))

	// Sort heapAllocs to find p99
	sort.Slice(heapAllocs, func(i, j int) bool { return heapAllocs[i] < heapAllocs[j] })
	p99HeapAlloc := heapAllocs[len(heapAllocs)*99/100]

	fmt.Printf("Average HeapAlloc: %d bytes\n", avgHeapAlloc)
	fmt.Printf("P99 HeapAlloc: %d bytes\n", p99HeapAlloc)

	// Capture memory profile
	memProfile, err := os.Create(filepath.Join(dir, "mem.prof"))
	if err != nil {
		b.Fatalf("could not create memory profile: %v", err)
	}
	defer memProfile.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		b.Fatalf("could not write memory profile: %v", err)
	}
}
