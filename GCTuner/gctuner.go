package main

import (
	"runtime"
	"time"
)

// Idea for finalizers taken from: https://www.uber.com/blog/how-we-saved-70k-cores-across-30-mission-critical-services/
// The intuition is that finalizers act "sort of" like a destructor, but are run whenever Go decides to "clean" this object
// during GC. This can be somewhat unpredicatable, but exactly what we want to determine when to run the GC Tuner and change
// the value.
type finalizer struct {
	ch  chan time.Time
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

func finalizerHandler(f *finalizerRef) {
	select {
	case f.parent.ch <- time.Time{}:
	default:
	}
	runtime.SetFinalizer(f, finalizerHandler)
}

// Setup a func with an inf. for-select loop on f.parent.ch, which will trigger
// grabbing memory metrics, calculating new value of GOGC/GOMEMLIMIT and setting it.

func InitGCTuner() *finalizer {
	// TODO: add actual stats read from mem stats, and do work off of this
	f := &finalizer{
		ch: make(chan time.Time, 1),
	}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil

	return f
}
