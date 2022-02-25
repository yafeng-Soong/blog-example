package main

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func Benchmark_muti(b *testing.B) {
	runtime.GOMAXPROCS(1)
	N := 50
	b.StartTimer()
	for i := 0; i < N; i++ {
		wg.Add(1)
		go IoWithGoroutine()
	}
	wg.Wait()
}

func IoWithGoroutine() {
	defer wg.Done()
	// 在sleep前后做少量运算
	for i := 0; i < 10; {
		i++
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; {
		i++
	}
}
