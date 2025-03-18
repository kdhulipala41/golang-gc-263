# A Survey on the Golang Garbage Collector

## Project Vision and Plan

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
