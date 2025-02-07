package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var sink []byte // Global var, doesn't allow compiler to optimize it out.

// Allocates 20480 bytes every 0.1 seconds.
func simpleAllocLoop() {
	_ = make([]byte, 500000)
	for {
		sink = make([]byte, 20480)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Set GOMEMLIMIT to 100MB
	err := os.Setenv("GOMEMLIMIT", "750000B")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("runs", os.ModePerm)
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
