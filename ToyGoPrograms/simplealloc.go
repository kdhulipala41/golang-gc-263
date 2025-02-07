package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// Allocates 1024 bytes every 0.5 seconds.
func simpleAllocLoop() {
	var allocations [][]byte
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			allocations = nil
		default:
			x := make([]byte, 1024)
			allocations = append(allocations, x)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	// Set GOMEMLIMIT to 100MB
	err := os.Setenv("GOMEMLIMIT", "750000B")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("runs/run-" + time.Now().Format("15-04-05"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("time,heap_inuse,gc_cpu_fraction\n")
	go simpleAllocLoop()

	var memStats runtime.MemStats
	start := time.Now()
	for time.Since(start) < 1*time.Minute {
		runtime.ReadMemStats(&memStats)

		// Log heap usage (bytes)
		secondsSinceStart := time.Since(start).Seconds()
		line := fmt.Sprintf("%.0f,%d,%.6f\n", secondsSinceStart, memStats.HeapInuse, memStats.GCCPUFraction)
		file.WriteString(line)

		time.Sleep(1 * time.Second)
	}
}
