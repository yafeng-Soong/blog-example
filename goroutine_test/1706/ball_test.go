package main

import (
	"testing"
)

func Benchmark_Nomal(b *testing.B) {
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
	findBall(grid)
	// fmt.Println(res)
}

func findBall(grid [][]int) []int {
	lr := len(grid)
	lc := len(grid[0])
	res := make([]int, lc)
	for k, _ := range grid[0] {
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
	}
	return res
}
