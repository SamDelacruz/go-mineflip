package model

import (
	"github.com/jmcvetta/randutil"
)

func GenBoard() [25]byte {
	var b [25]byte

	choices := []randutil.Choice{
		{Item: 0, Weight: 6},
		{Item: 1, Weight: 13},
		{Item: 2, Weight: 4},
		{Item: 3, Weight: 2},
	}

	for i, _ := range b {
		v, _ := randutil.WeightedChoice(choices)
		b[i] = byte(v.Item.(int))
	}

	return b
}
