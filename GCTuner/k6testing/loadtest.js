import http from "k6/http";
import { check, sleep } from "k6";

// Define the options for the load test
export const options = {
  stages: [
    { duration: "30s", target: 20 }, // Ramp up to 20 within first 30s
    { duration: "10s", target: 20 }, // Stay 20VUs
    { duration: "30s", target: 0 }, // Drop out in last 30s
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
    http_req_failed: ["rate<0.01"], // Less than 1% of requests should fail
  },
};

// The main function that k6 will execute
export default function () {
  // Define the endpoint you want to test
  const url = `http://${__ENV.MY_HOSTNAME}/mergesort`;

  // Send a GET request to the /mergesort endpoint
  const res = http.get(url);

  // Check if the response status is 200
  check(res, {
    "is status 200": (r) => r.status === 200,
  });

  // Add a sleep to simulate think time between requests
  sleep(1);
}
