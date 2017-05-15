package model_test

import (
	"github.com/samdelacruz/go-mineflip/model"
	"testing"
)

func TestGame_GetScoreSimple(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		1, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{0, 1, 2, 24}}

	want := 6
	got := g.GetScore()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetScoreZero(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{0, 1, 2, 24}}

	want := 0
	got := g.GetScore()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetVisibleNewGame(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{}}

	want := "?????????????????????????"
	got := g.GetVisible()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetVisibleRevealed(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{10, 24}}

	want := "??????????1?????????????2"
	got := g.GetVisible()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}
