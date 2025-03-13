package gctuner

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

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Setup a func with an inf. for-select loop on f.parent.ch, which will trigger
// grabbing memory metrics, calculating new value of GOGC/GOMEMLIMIT and setting it.
func setGCValueAIMD(f *finalizerRef) {
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		// fmt.Printf("Current: %v, Past: %v, Diff: %v\n", currHeapAlloc, prevHeapAlloc, currDiffAlloc)
		var newGCPercent int
		if currDiffAlloc-prevDiffAlloc > 0 {
			// GC more aggressive
			newGCPercent = max(minGCPercent, int(0.7*float32(currGCPercent)))
		} else {
			// GC more relaxed
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

// Dynamically modifies GOGC/GOMEMLIMIT to either 60 or 500
// based on positive/negative difference in allocation speed
func setGCValueFlipFlop(f *finalizerRef) {
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		var newGCPercent int
		if currDiffAlloc-prevDiffAlloc > 0 {
			// GC more aggressive
			newGCPercent = minGCPercent
		} else {
			// GC more relaxed
			newGCPercent = maxGCPercent
		}
		if newGCPercent != currGCPercent {
			setGCValue(newGCPercent)
		}
		prevDiffAlloc = currDiffAlloc
		prevHeapAlloc = currHeapAlloc
	}
}

// Dynamically modifies GOGC/GOMEMLIMIT to either 60 or 500
// based on allocation speed being above or below 1GB threshold
func setGCValueThreshold(f *finalizerRef) {
	const threshold = 1024 * 1024 * 1024
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		var newGCPercent int
		if currDiffAlloc > threshold {
			// GC more aggressive
			newGCPercent = minGCPercent
		} else {
			// GC more relaxed
			newGCPercent = maxGCPercent
		}
		if newGCPercent != currGCPercent {
			setGCValue(newGCPercent)
		}
		prevHeapAlloc = currHeapAlloc
	}
}

// Dynamically scales GOGC/GOMEMLIMIT linearly based on allocation rate
// (normalized into the range of [minGCPercent, maxGCPercent])
func setGCValueLinear(f *finalizerRef) {
	// alloc rates are min of 1MB/sec and max of 50MB/sec, these are somewhat arbitrarily chosen but can possibly be extended
	// to be read off a profile?
	const maxAllocRate = 1024 * 1024 * 50
	const minAllocRate = 1024 * 1024 * 1
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		normalizedGC := maxGCPercent - int((float64(currDiffAlloc-minAllocRate)/float64(maxAllocRate-minAllocRate))*(float64(maxGCPercent-minGCPercent)))
		normalizedGC = clamp(normalizedGC, minGCPercent, maxGCPercent)
		if normalizedGC != currGCPercent {
			setGCValue(normalizedGC)
		}
		prevHeapAlloc = currHeapAlloc
	}
}

// Dynamically scales GOGC/GOMEMLIMIT from a rolling average
// from a history of the 5 most recent allocation speeds
func setGCValueRollingAvg(f *finalizerRef) {
	const windowSize = 5
	var allocHistory [windowSize]int64
	var index int
	for range f.parent.ch {
		currHeapAlloc := readHeapAllocStats()
		currDiffAlloc := int64(currHeapAlloc) - int64(prevHeapAlloc)
		allocHistory[index%windowSize] = currDiffAlloc
		index++
		var total int64
		for _, v := range allocHistory {
			total += v
		}
		avgAllocRate := total / windowSize
		var newGCPercent int
		if avgAllocRate > prevDiffAlloc {
			// GC more aggressive
			newGCPercent = minGCPercent
		} else {
			// GC more relaxed
			newGCPercent = maxGCPercent
		}
		if newGCPercent != currGCPercent {
			setGCValue(newGCPercent)
		}
		prevHeapAlloc = currHeapAlloc
		prevDiffAlloc = avgAllocRate
	}
}

// Add options, and finish above function to read memory limit and set it based on the option.
func InitGCTuner(tunerType int) *finalizer {
	// Set the GOMEMLIMIT to 90% of the cgroup's memory limit or the system's memory limit.
	memlimit.SetGoMemLimitWithOpts(
		memlimit.WithRatio(0.7),
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

	switch tunerType {
	case 0:
		go setGCValueAIMD(f.ref)
	case 1:
		go setGCValueRollingAvg(f.ref)
	case 2:
		go setGCValueLinear(f.ref)
	case 3:
		go setGCValueFlipFlop(f.ref)
	case 4:
		go setGCValueThreshold(f.ref)
	}
	f.ref = nil

	return f
}
