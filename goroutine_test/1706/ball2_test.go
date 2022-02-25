package main

import (
	"runtime"
	"sync"
	"testing"
)

func Benchmark_Sync(b *testing.B) {
	runtime.GOMAXPROCS(1)
	N := 20000
	grid := make([][]int, N)
	for i := 0; i < N; i++ {
		grid[i] = make([]int, N)
		tmp := 1
		if i%2 == 0 {
			tmp = -1
		}
		for j := 0; j < N; j++ {
			grid[i][j] = tmp
		}
	}
	b.StartTimer()
	findBallSync(grid)
	// fmt.Println(res)
}
func findBallSync(grid [][]int) []int {
	var wg sync.WaitGroup
	lr := len(grid)
	lc := len(grid[0])
	res := make([]int, lc)
	for index, _ := range grid[0] {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			i := 0
			j := k
			for i < lr {
				dir := grid[i][j]
				j += dir
				if j < 0 || j >= lc || grid[i][j] != dir {
					j = -1
					break
				}
				i++
			}
			res[k] = j
		}(index)
	}
	wg.Wait()
	return res
}
