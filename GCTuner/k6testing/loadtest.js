import http from "k6/http";
import { check, sleep } from "k6";

// Define the options for the load test
export const options = {
  stages: [
    // Initial ramp-up to warm up the system
    { duration: "30s", target: 50 }, // Ramp up to 50 users over 30 seconds
    { duration: "30s", target: 50 }, // Stay at 50 users

    // Burst 1: Sudden spike to 200 users
    { duration: "10s", target: 100 }, // Spike to 200 users for 10 seconds
    { duration: "20s", target: 50 }, // Drop back to 50 users for 20 seconds
    { duration: "30s", target: 50 }, // Stay at 50 users

    // Burst 2: Another spike to 300 users
    { duration: "10s", target: 120 }, // Spike to 300 users for 10 seconds
    { duration: "20s", target: 50 }, // Drop back to 50 users for 20 seconds
    { duration: "30s", target: 50 }, // Stay at 50 users

    // Burst 3: Final spike to 400 users
    { duration: "10s", target: 150 }, // Spike to 400 users for 10 seconds
    { duration: "20s", target: 50 }, // Drop back to 50 users for 20 seconds

    // Ramp down to 0 users
    { duration: "30s", target: 0 }, // Gradually reduce load to 0 users
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
    http_req_failed: ["rate<0.01"], // Less than 1% of requests should fail
  },
};

// The main function that k6 will execute
export default function () {
  // Define the endpoint you want to test
  const url = `http://${__ENV.MY_HOSTNAME}/astparse`;

  // Send a GET request to the /mergesort endpoint
  const res = http.get(url);

  // Check if the response status is 200
  check(res, {
    "is status 200": (r) => r.status === 200,
  });

  // Add a sleep to simulate think time between requests
  sleep(1);
}
