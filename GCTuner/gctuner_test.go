package main

import (
	"cs263/GCTuner/astparse"
	"testing"
)

func TestGCTunerBasic(t *testing.T) {
	InitGCTuner()
	astparse.BenchmarkN(100)
}
