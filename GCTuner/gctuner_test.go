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

func BenchmarkASTParseWithAIMDTuner(b *testing.B) {
	gctuner.InitGCTuner(0)
	profileBenchmark(b, "BenchmarkWithAIMDTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithRollingAvgTuner(b *testing.B) {
	gctuner.InitGCTuner(1)
	profileBenchmark(b, "BenchmarkASTParseWithRollingAvgTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithLinearTuner(b *testing.B) {
	gctuner.InitGCTuner(2)
	profileBenchmark(b, "BenchmarkASTParseWithLinearTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithFlipFlopTuner(b *testing.B) {
	gctuner.InitGCTuner(3)
	profileBenchmark(b, "BenchmarkASTParseWithFlipFlopTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithThresholdTuner(b *testing.B) {
	gctuner.InitGCTuner(4)
	profileBenchmark(b, "BenchmarkASTParseWithThresholdTuner", func() {
		astparse.BenchmarkN(10000)
	})
}
