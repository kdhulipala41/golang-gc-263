package main

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

const (
	minGCPercent = 60
	maxGCPercent = 500
)

var memStats runtime.MemStats
var prevHeapAlloc uint64
var prevDiffAlloc int64
var currGCPercent int

// Idea for finalizers taken from: https://www.uber.com/blog/how-we-saved-70k-cores-across-30-mission-critical-services/
// The intuition is that finalizers act "sort of" like a destructor, but are run whenever Go decides to "clean" this object
// during GC. This can be somewhat unpredicatable, but exactly what we want to determine when to run the GC Tuner and change
// the value.
func readHeapAllocStats() uint64 {
	runtime.ReadMemStats(&memStats)
	return memStats.HeapAlloc
}

func setGCValue(percent int) {
	currGCPercent = percent
	debug.SetGCPercent(percent)
}

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
func setDynamicGCValue(f *finalizerRef) {
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		// fmt.Printf("Current: %v, Past: %v, Diff: %v\n", currHeapAlloc, prevHeapAlloc, currDiffAlloc)
		var newGCPercent int
		if currDiffAlloc-prevDiffAlloc > 0 {
			newGCPercent = max(minGCPercent, int(0.7*float32(currGCPercent)))
		} else {
			newGCPercent = min(maxGCPercent, currGCPercent+40)
		}
		if newGCPercent != currGCPercent {
			// fmt.Printf("CurrGCPercent: %v, NewGCPercent%v\n", currGCPercent, newGCPercent)
			setGCValue(newGCPercent)
		}
		prevDiffAlloc = currDiffAlloc
		prevHeapAlloc = currHeapAlloc
	}
}

// Add options, and finish above function to read memory limit and set it based on the option.
func InitGCTuner() *finalizer {
	// Set the GOMEMLIMIT to 90% of the cgroup's memory limit or the system's memory limit.
	memlimit.SetGoMemLimitWithOpts(
		memlimit.WithRatio(0.9),
		memlimit.WithProvider(
			memlimit.ApplyFallback(
				memlimit.FromCgroup,
				memlimit.FromSystem,
			),
		),
		memlimit.WithRefreshInterval(1*time.Minute),
	)
	currGCPercent = 100
	f := &finalizer{
		ch: make(chan time.Time, 1),
	}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	go setDynamicGCValue(f.ref)
	f.ref = nil

	return f
}
