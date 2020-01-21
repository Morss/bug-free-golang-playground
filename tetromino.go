package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Tet struct {
	size int
	mat  [][]bool
}

var (
	source = rand.NewSource(time.Now().UnixNano())
	rng    = rand.New(source)
)

func GetRandom() Tet {
	switch rng.Intn(7) {
	case 0:
		return getBlock()
	case 1:
		return getLine()
	case 2:
		return getLcis()
	case 3:
		return getLtrans()
	case 4:
		return getZcis()
	case 5:
		return getZtrans()
	}
	return getT()
}

func (t *Tet) Rot() {
	res := getSquaredBoolSlice(t.size)

	for m := 0; m < t.size; m++ {
		for n := 0; n < t.size; n++ {
			res[t.size-1-n][m] = t.mat[m][n]
		}
	}
	t.mat = res
}

func (t *Tet) Draw() {
	for m := t.size - 1; m > -1; m-- {
		for n := 0; n < t.size; n++ {
			if t.mat[m][n] {
				fmt.Print("x ")
			} else {
				fmt.Print("o ")
			}
		}
		fmt.Println()
	}
}

func getBlock() Tet {
	tMatSize := 2
	grid := getSquaredBoolSlice(tMatSize)

	grid[0][0] = true
	grid[0][1] = true
	grid[1][0] = true
	grid[1][1] = true

	return Tet{tMatSize, grid}
}

func getLine() Tet {
	tMatSize := 4
	grid := getSquaredBoolSlice(tMatSize)

	grid[1][0] = true
	grid[1][1] = true
	grid[1][2] = true
	grid[1][3] = true

	return Tet{tMatSize, grid}
}

func getLcis() Tet {
	tMatSize := 3
	grid := getSquaredBoolSlice(tMatSize)

	grid[0][1] = true
	grid[1][1] = true
	grid[2][1] = true
	grid[0][2] = true

	return Tet{tMatSize, grid}
}

func getLtrans() Tet {
	tMatSize := 3
	grid := getSquaredBoolSlice(tMatSize)

	grid[0][1] = true
	grid[1][1] = true
	grid[2][1] = true
	grid[0][0] = true

	return Tet{tMatSize, grid}
}

func getZcis() Tet {
	tMatSize := 3
	grid := getSquaredBoolSlice(tMatSize)

	grid[0][0] = true
	grid[0][1] = true
	grid[1][1] = true
	grid[1][2] = true

	return Tet{tMatSize, grid}
}

func getZtrans() Tet {
	tMatSize := 3
	grid := getSquaredBoolSlice(tMatSize)

	grid[1][0] = true
	grid[1][1] = true
	grid[0][1] = true
	grid[0][2] = true

	return Tet{tMatSize, grid}
}

func getT() Tet {
	tMatSize := 3
	grid := getSquaredBoolSlice(tMatSize)

	grid[1][0] = true
	grid[1][1] = true
	grid[1][2] = true
	grid[2][1] = true

	return Tet{tMatSize, grid}
}

func getSquaredBoolSlice(size int) [][]bool {
	res := make([][]bool, size)
	for i := 0; i < size; i++ {
		res[i] = make([]bool, size)
	}
	return res
}
