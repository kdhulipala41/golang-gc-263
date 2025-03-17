package main

// This package contains our API hosted on the EC2 instance. It contains some simple endpoints that will run minimized
// versions of the short-lived tests, and also includes a stats endpoint which queries the runtime memstats to obtain values
// like size of the live heap or # of gc runs.

import (
	"encoding/json"
	"flag"
	"net/http"
	"runtime"

	"github.com/kdhulipala41/golang-gc-263/GCTuner/astparse"
	"github.com/kdhulipala41/golang-gc-263/GCTuner/gctuner"
	nestedptrmap "github.com/kdhulipala41/golang-gc-263/GCTuner/nesterptrmap"
)

// Calls nestedptrmap.
func nestedptrmapHandler(w http.ResponseWriter, r *http.Request) {
	nestedptrmap.InitAndMutateNestedPtrMap()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success\n"))
}

// Calls one iteration of AST parse.
func astparseHandler(w http.ResponseWriter, r *http.Request) {
	astparse.ParsePackage()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success\n"))
}

// Handler to return current map size and memory stats
func statsHandler(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := map[string]interface{}{
		"heap_alloc":    memStats.HeapAlloc,
		"heap_objects":  memStats.HeapObjects,
		"gc_runs":       memStats.NumGC,
		"last_gc_pause": memStats.PauseNs[(memStats.NumGC+255)%256], // Most recent GC pause time
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

var liveMemory []byte

// Bloats memory by allocating a large 600MB global var that will stay during the lifetime of the program.
func allocateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	liveMemory = make([]byte, 600*1024*1024) // Allocate 600MB
	for i := range liveMemory {
		liveMemory[i] = 1 // Touch the memory to ensure it's allocated
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Allocated 600MB of memory\n"))
}

func main() {
	var tunerType int
	flag.IntVar(&tunerType, "tunerType", -1, "Type of GC tuner to use: 0 - AIMD, 1 - Rolling Avg, 2 - Linear, 3- Flip Flop, 4 - GC Value Threshold")

	var memLimitFrac float64
	flag.Float64Var(&memLimitFrac, "memLimitFrac", 0.8, "The fraction of the container or system memory limit that should be set as GOMEMLIMIT. Defaults to 0.8")
	flag.Parse()

	if tunerType != -1 {
		gctuner.InitGCTuner(tunerType, memLimitFrac)
	}

	http.HandleFunc("/nestedptrmap", nestedptrmapHandler)
	http.HandleFunc("/astparse", astparseHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/bigallocate", allocateMemoryHandler)

	http.ListenAndServe(":8080", nil)
}
