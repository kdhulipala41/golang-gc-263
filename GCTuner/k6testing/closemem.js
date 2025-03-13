// Similar to the loadtest.js, but has another scenario which runs concurrently that adds 700MB of bloat in live-heap to
// test application load when memory stays close to the system limit.

import http from "k6/http";
import { check, sleep } from "k6";
import { Gauge, Trend } from "k6/metrics";

// Custom metrics to track heap allocation and GC stats
const heapAlloc = new Trend("heap_alloc_bytes");
const heapObjects = new Trend("heap_objects");
const gcRuns = new Gauge("gc_runs");
const lastGcPause = new Trend("last_gc_pause_ns");

// Define the options for the load test
export const options = {
  scenarios: {
    // Scenario 1: Main load test on /astparse
    load_test: {
      executor: "ramping-vus",
      stages: [
        // Initial ramp-up to warm up the system
        { duration: "30s", target: 50 }, // Ramp up to 50 users over 30 seconds
        { duration: "30s", target: 50 }, // Stay at 50 users

        // Burst 1: Sudden spike to 200 users
        { duration: "10s", target: 120 }, // Spike to 200 users for 10 seconds
        { duration: "30s", target: 50 }, // Drop back to 50 users for 20 seconds
        { duration: "40s", target: 0 }, // Stay at 50 users

        // Burst 2: Another spike to 300 users
        { duration: "10s", target: 135 }, // Spike to 300 users for 10 seconds
        { duration: "30s", target: 50 }, // Drop back to 50 users for 20 seconds
        { duration: "40s", target: 0 }, // Stay at 50 users

        // Burst 3: Final spike to 400 users
        { duration: "10s", target: 150 }, // Spike to 400 users for 10 seconds
        { duration: "20s", target: 50 }, // Drop back to 50 users for 20 seconds

        // Ramp down to 0 users
        { duration: "30s", target: 0 }, // Gradually reduce load to 0 users
      ],
      exec: "mainLoadTest", // Function to execute for this scenario
    },
    // Scenario 2: Single VU to query /stats
    stats_monitor: {
      executor: "per-vu-iterations",
      vus: 1, // Only 1 VU
      iterations: 1, // Large number of iterations (will stop when load_test ends)
      maxDuration: "230s",
      exec: "statsMonitor", // Function to execute for this scenario
      startTime: "0s", // Start immediately
      gracefulStop: "5s", // Allow 5 seconds to finish after load_test ends
    },
    // Memory Bloat
    memory_bloat: {
      executor: "per-vu-iterations",
      vus: 1, // Only 1 VU
      iterations: 1, // Large number of iterations (will stop when load_test ends)
      maxDuration: "230s",
      exec: "bloatMemory", // Function to execute for this scenario
      startTime: "0s", // Start immediately
      gracefulStop: "5s", // Allow 5 seconds to finish after load_test ends
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
    http_req_failed: ["rate<0.01"], // Less than 1% of requests should fail
  },
  cloud: {
    // Test runs with the same name groups test runs together
    name: "CS263AstParse",
  },
};

// Function to query the /stats endpoint and collect metrics
function getStats() {
  const statsUrl = `http://${__ENV.MY_HOSTNAME}/stats`;
  const res = http.get(statsUrl);

  if (res.status === 200) {
    const stats = JSON.parse(res.body);
    heapAlloc.add(stats.heap_alloc); // Track heap allocation
    heapObjects.add(stats.heap_objects); // Track heap objects
    gcRuns.add(stats.gc_runs); // Track total GC runs
    lastGcPause.add(stats.last_gc_pause); // Track last GC pause time
  }
}

// Scenario 1: Main load test on /astparse
export function mainLoadTest() {
  const url = `http://${__ENV.MY_HOSTNAME}/astparse`;
  const res = http.get(url);

  // Check if the response status is 200
  check(res, {
    "is status 200": (r) => r.status === 200,
  });

  // Add a sleep to simulate think time between requests
  sleep(1);
}

// Scenario 2: Single VU to query /stats
export function statsMonitor() {
  while (true) {
    getStats();
    sleep(5); // Query every 5 seconds
  }
}

// Scenario 3: Bloat Memory
export function bloatMemory() {
  const statsUrl = `http://${__ENV.MY_HOSTNAME}/bigallocate`;
  const res = http.get(statsUrl);

  // Check if the response status is 200
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
}
