package main

import (
 "fmt"
 "os"
 "runtime/debug"
 "runtime/trace"
 "sync"
)

const (
 NumWorkers    = 4
 NumTasks      = 500
 MemoryIntense = 10000 // Size of memory-intensive task (number of elements).
)


func main() {
 // Write to the trace file.
 f, _ := os.Create("trace.out")
 trace.Start(f)
 defer trace.Stop()

 // Set the target percentage for the garbage collector. Default is 100%.
 debug.SetGCPercent(100)

 // Task queue and result queue.
 taskQueue := make(chan int, NumTasks)
 resultQueue := make(chan int, NumTasks)

 // Start workers.
 var wg sync.WaitGroup
 wg.Add(NumWorkers)
 for i := 0; i < NumWorkers; i++ {
  go worker(taskQueue, resultQueue, &wg)
 }

 // Send tasks to the queue.
 for i := 0; i < NumTasks; i++ {
  taskQueue <- i


 }
 close(taskQueue)

 // Retrieve results from the queue.
 go func() {
  wg.Wait()
  close(resultQueue)
 }()

 // Process the results.
 for result := range resultQueue {
  fmt.Println("Result:", result)
 }

 fmt.Println("Done!")
}

// Worker function.
func worker(tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
 defer wg.Done()

 for task := range tasks {
  result := performMemoryIntensiveTask(task)
  results <- result
 }
}


func performMemoryIntensiveTask(task int) int {
 // Create a large-sized slice.
 data := make([]int, MemoryIntense)
 for i := 0; i < MemoryIntense; i++ {
  data[i] = i + task
 }

 // Imitation of latency.
 time.Sleep(10 * time.Millisecond)

 // Calculate the result.
 result := 0
 for eachValue := range data {
  result += eachValue
 }
 return result
}
