package main

import (
	"cs263/GCTuner/astparse"
	"cs263/GCTuner/gctuner"
	"testing"
)

func BenchmarkASTParseNoTuner(b *testing.B) {
	profileBenchmark(b, "BenchmarkNoTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithTuner(b *testing.B) {
	gctuner.InitGCTuner()
	profileBenchmark(b, "BenchmarkWithDefaultTuner", func() {
		astparse.BenchmarkN(10000)
	})
}
