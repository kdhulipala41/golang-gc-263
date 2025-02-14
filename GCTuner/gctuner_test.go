package main

import (
	"cs263/GCTuner/astparse"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"testing"
)

func BenchmarkASTParseDefault(b *testing.B) {
	profileBenchmark(b, "BenchmarkASTParseDefault", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCOff(b *testing.B) {
	b.Setenv("GOGC", "OFF")
	profileBenchmark(b, "BenchmarkASTParseGCOff", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCLow(b *testing.B) {
	b.Setenv("GOGC", "10")
	profileBenchmark(b, "BenchmarkASTParseGCLow", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseGCMid(b *testing.B) {
	b.Setenv("GOGC", "50")
	profileBenchmark(b, "BenchmarkASTParseGCMid", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseLowMemLimit(b *testing.B) {
	b.Setenv("GOMEMLIMIT", "10000MB")
	profileBenchmark(b, "BenchmarkASTParseLowMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseMidMemLimit(b *testing.B) {
	b.Setenv("GOMEMLIMIT", "20000MB")
	profileBenchmark(b, "BenchmarkASTParseMidMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseHighMemLimit(b *testing.B) {
	b.Setenv("GOMEMLIMIT", "40000MB")
	profileBenchmark(b, "BenchmarkASTParseHighMemLimit", func() {
		astparse.BenchmarkN(10000)
	})
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

	// Run the benchmark
	b.ResetTimer()
	benchmarkFunc()
	b.StopTimer()

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
