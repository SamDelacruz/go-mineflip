package model

import (
	"errors"
	"math/rand"
	"time"

	"math"

	"github.com/jmcvetta/randutil"
)

func GenBoardForLevel(level int) [25]byte {
	vals := genVals(level)
	idxs := nRandomIndices(25)
	board, _ := GenBoardFromVals(vals, idxs)
	return *board
}

// GenBoardFromVals given a set of tile values, and respective indices, return a board.
func GenBoardFromVals(vals []byte, idxs []int) (*[25]byte, error) {
	if len(vals) != len(idxs) {
		return nil, errors.New("Number of vals != number of idxs")
	}
	var board [25]byte
	boardLen := len(board)

	for i, val := range vals {
		boardIdx := idxs[i]
		if boardIdx+1 > boardLen || boardIdx < 0 {
			return nil, errors.New("Board index out out bounds -> " + string(boardIdx))
		}
		board[boardIdx] = val
	}
	return &board, nil
}

func genVals(level int) (vals []byte) {
	vals = make([]byte, 0)
	for i, n := range genValCounts(level) {
		counts := make([]byte, n)
		for j := range counts {
			counts[j] = byte(i)
		}
		vals = append(vals, counts...)
	}
	return vals
}

func genValCounts(level int) []int {
	numZeros, _ := randutil.IntRange(minMaxZero(level))
	numTwos, _ := randutil.IntRange(minMaxTwo(level))
	numThrees, _ := randutil.IntRange(minMaxThree(level))
	numOnes := 25 - numZeros - numTwos - numThrees
	return []int{numZeros, numOnes, numTwos, numThrees}
}

func minMaxZero(level int) (min int, max int) {
	return minMax(level, 0.5, 5, 0.5, 6)
}

func minMaxTwo(level int) (min int, max int) {
	return minMax(level, 0.5, 2.5, 0.6, 3.5)
}

func minMaxThree(level int) (min int, max int) {
	return minMax(level, 0.6, 0.4, 0.6, 2)
}

func minMax(level int, m, c, n, d float64) (min int, max int) {
	min = int(math.Ceil(float64(level)*m + c))
	max = int(math.Floor(float64(level)*n+d)) + 1
	return min, max
}

func nRandomIndices(n int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idxs := r.Perm(n)
	return idxs
}
