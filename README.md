# A Survey on the Golang Garbage Collector

### Vision

The Go garbage collector implements a concurrent, tri-color, mark and sweep algorithm that is triggered periodically or based on memory pressure (heap target has been exceeded). Although the Go garbage collector provides a knob GOGC to adjust the size of the total heap relative to the size of live objects, applications that see frequent and/or high variations in load can still struggle with high CPU usage from GC cycles and possibly even OOMs due to the GC being too aggressive or too passive. We would like to investigate the intricacies of Go’s garbage collector, hoping to improve on its runtime pacer by leveraging the GOGC knob dynamically along with other GC knobs such as GOMEMLIMIT.

### Installation

If the directory is not already in a module, create one in the root directory:

```shell
go mod init example.com/test
```

### Usage

Import the following line and call the constructor for the gctuner:

```go
package main

import "github.com/kdhulipala41/golang-gc-263/GCTuner/gctuner"

func main() {
  // Giving it values of 3 = Flip Flop Tuner, and 0.8 or 80% of container/system limit
	gctuner.InitGCTuner(3, 0.8)
}
```

Then proceed to run:

```shell
go mod tidy
```

This will install the dependencies needed. Now you should be ready to Go!
