package main

import (
	"cs263/GCTuner/mergesort"
	"encoding/json"
	"net/http"
	"runtime"
)

// Handler to insert random data into the map
func mergesortHandler(w http.ResponseWriter, r *http.Request) {
	mergesort.AllocateNAndSort(10000)

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

func main() {
	http.HandleFunc("/mergesort", mergesortHandler)
	http.HandleFunc("/stats", statsHandler)

	http.ListenAndServe(":8080", nil)
}
