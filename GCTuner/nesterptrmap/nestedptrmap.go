package nestedptrmap

import (
	"math/rand"
	"sync"
	"time"
)

const totalSize = 50 * 1024 * 1024 // 10MB in bytes
const numKeys = 100000
const depth = 30

type DeepPtr struct {
	Next      *DeepPtr
	ByteSlice []byte
}

func NewDeepPtr(depth int) *DeepPtr {
	if depth == 0 {
		return nil
	}
	return &DeepPtr{Next: NewDeepPtr(depth - 1), ByteSlice: make([]byte, totalSize/(numKeys*depth))}
}

func InitAndMutateNestedPtrMap() {
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localMap := make(map[int]*DeepPtr, numKeys)

			for j := 0; j < numKeys; j++ {
				localMap[j] = NewDeepPtr(depth)
			}

			for j := 0; j < 50; j++ {
				DeleteN(localMap, 1000, rand.Intn(numKeys))
				InsertN(localMap, numKeys, rand.Intn(numKeys))
				time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

func DeleteN(localMap map[int]*DeepPtr, n int, start int) {
	i := start
	for n > 0 {
		_, ok := localMap[i]
		if ok {
			delete(localMap, i)
			n -= 1
		}
		i = (i + 1) % numKeys
	}
}

func InsertN(localMap map[int]*DeepPtr, n int, start int) {
	i := start
	for n > 0 {
		localMap[i] = NewDeepPtr(depth)
		n -= 1
		i = (i + 1) % numKeys
	}
}
