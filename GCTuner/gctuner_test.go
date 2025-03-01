package main

import (
	"cs263/GCTuner/astparse"
	"testing"
)

func BenchmarkNoTuner(b *testing.B) {
	profileBenchmark(b, "BenchmarkNoTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkWithDefaultTuner(b *testing.B) {
	InitGCTuner()
	profileBenchmark(b, "BenchmarkWithDefaultTuner", func() {
		astparse.BenchmarkN(10000)
	})
}
