package main

import (
	"cs263/GCTuner/astparse"
	"cs263/GCTuner/gctuner"
	nestedptrmap "cs263/GCTuner/nesterptrmap"
	"testing"
)

func BenchmarkASTParseNoTuner(b *testing.B) {
	profileBenchmark(b, "BenchmarkNoTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithAIMDTuner(b *testing.B) {
	gctuner.InitGCTuner(0, 0.8)
	profileBenchmark(b, "BenchmarkWithAIMDTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithRollingAvgTuner(b *testing.B) {
	gctuner.InitGCTuner(1, 0.8)
	profileBenchmark(b, "BenchmarkASTParseWithRollingAvgTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithLinearTuner(b *testing.B) {
	gctuner.InitGCTuner(2, 0.8)
	profileBenchmark(b, "BenchmarkASTParseWithLinearTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithFlipFlopTuner(b *testing.B) {
	gctuner.InitGCTuner(3, 0.8)
	profileBenchmark(b, "BenchmarkASTParseWithFlipFlopTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkASTParseWithThresholdTuner(b *testing.B) {
	gctuner.InitGCTuner(4, 0.8)
	profileBenchmark(b, "BenchmarkASTParseWithThresholdTuner", func() {
		astparse.BenchmarkN(10000)
	})
}

func BenchmarkNestedPtrMapNoTuner(b *testing.B) {
	profileBenchmark(b, "BenchmarkNestedPtrMapNoTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}

func BenchmarkNestedPtrMapWithAIMDTuner(b *testing.B) {
	gctuner.InitGCTuner(0, 0.8)
	profileBenchmark(b, "BenchmarkNestedPtrMapWithAIMDTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}

func BenchmarkNestedPtrMapWithRollingAvgTuner(b *testing.B) {
	gctuner.InitGCTuner(1, 0.8)
	profileBenchmark(b, "BenchmarkNestedPtrMapWithRollingAvgTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}

func BenchmarkNestedPtrMapWithLinearTuner(b *testing.B) {
	gctuner.InitGCTuner(2, 0.8)
	profileBenchmark(b, "BenchmarkNestedPtrMapWithLinearTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}

func BenchmarkNestedPtrMapWithFlipFlopTuner(b *testing.B) {
	gctuner.InitGCTuner(3, 0.8)
	profileBenchmark(b, "BenchmarkNestedPtrMapNoTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}

func BenchmarkNestedPtrMapWithThresholdTuner(b *testing.B) {
	gctuner.InitGCTuner(4, 0.8)
	profileBenchmark(b, "BenchmarkNestedPtrMapWithThresholdTuner", func() {
		nestedptrmap.InitAndMutateNestedPtrMap()
	})
}
