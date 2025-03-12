package main

import (
	"cs263/GCTuner/astparse"
	"cs263/GCTuner/gctuner"
	nestedptrmap "cs263/GCTuner/nesterptrmap"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"runtime"
)

func nestedptrmapHandler(w http.ResponseWriter, r *http.Request) {
	nestedptrmap.InitAndMutateNestedPtrMap()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success\n"))
}

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

func allocateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	liveMemory = make([]byte, 700*1024*1024) // Allocate 700MB
	for i := range liveMemory {
		liveMemory[i] = 1 // Touch the memory to ensure it's allocated
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Allocated 700MB of memory\n"))
}

func main() {
	var tunerType int
	flag.IntVar(&tunerType, "tunerType", -1, "Type of GC tuner to use: 0 - AIMD, 1 - Rolling Avg, 2 - Linear, 3- Flip Flop, 4 - GC Value Threshold")
	flag.Parse()

	if tunerType != -1 {
		fmt.Print("turning on gctuner")
		gctuner.InitGCTuner(tunerType)
	}

	http.HandleFunc("/nestedptrmap", nestedptrmapHandler)
	http.HandleFunc("/astparse", astparseHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/bigallocate", allocateMemoryHandler)

	http.ListenAndServe(":8080", nil)
}
