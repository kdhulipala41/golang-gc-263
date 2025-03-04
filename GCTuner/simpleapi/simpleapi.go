package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

var (
	hugeMap = make(map[int][]byte)
	mu      sync.Mutex
)

// Handler to insert random data into the map
func insertHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Insert 1000 random entries
	for i := 0; i < 1000; i++ {
		key := rand.Intn(1_000_000_000)
		hugeMap[key] = make([]byte, 1024) // 1KB per entry
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inserted 1000 entries\n"))
}

// Handler to delete random entries from the map
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Delete 1000 random keys
	for i := 0; i < 1000; i++ {
		key := rand.Intn(1_000_000_000)
		delete(hugeMap, key)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted 1000 entries\n"))
}

// Handler to return current map size and memory stats
func statsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	size := len(hugeMap)
	mu.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := map[string]interface{}{
		"map_size":      size,
		"heap_alloc":    memStats.HeapAlloc,
		"heap_objects":  memStats.HeapObjects,
		"gc_runs":       memStats.NumGC,
		"last_gc_pause": memStats.PauseNs[(memStats.NumGC+255)%256], // Most recent GC pause time
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/stats", statsHandler)

	http.ListenAndServe(":8080", nil)
}
