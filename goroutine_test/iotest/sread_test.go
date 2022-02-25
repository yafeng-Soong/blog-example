package main

import (
	"testing"
	"time"
)

func Benchmark_single(b *testing.B) {
	N := 50
	b.StartTimer()
	for i := 0; i < N; i++ {
		IO()
	}
}

func IO() {
	// 在sleep前后做少量运算
	for i := 0; i < 10; {
		i++
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; {
		i++
	}
}
