// Copyright 2014 The Go Authors. All rights reserved.
// This code has been adapted from benchmark material on the GO GC from the open-source go code.
// The nethttp.go just holds an old snapshot of net/http as a large string.

// This package contains code to parse a string and return a ptr to an AST File object, iterated N times.
package astparse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"sync"
	"sync/atomic"
)

type ParsedPackage *ast.File

var (
	parsed []ParsedPackage
)

func BenchmarkN(N uint64) {
	P := runtime.GOMAXPROCS(0)
	parsed = make([]ParsedPackage, 15)
	// Create G goroutines, but only 2*P of them parse at the same time.
	G := 1024
	gate := make(chan bool, 2*P)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(G)
	remain := int64(N)
	pos := 0
	for g := 0; g < G; g++ {
		go func() {
			defer wg.Done()
			for atomic.AddInt64(&remain, -1) >= 0 {
				gate <- true
				p := parsePackage()
				mu.Lock()
				// Overwrite only half of the array,
				// the other part represents "old" generation.
				parsed[pos%(len(parsed)/2)] = p
				pos++
				mu.Unlock()
				<-gate
			}
		}()
	}
	wg.Wait()
}

// parsePackage parses and returns net/http package.
func parsePackage() ParsedPackage {
	pkgs, err := parser.ParseFile(token.NewFileSet(), "net/http", src, parser.ParseComments)
	if err != nil {
		println("parse", err.Error())
		panic("fail")
	}
	return pkgs
}
