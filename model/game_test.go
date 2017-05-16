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

func TestGame_GetHintsBasic(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
	}}

	want := model.Hints{
		Rows: [5]model.Hint{
			{0, 5}, {0, 5}, {0, 5}, {0, 5}, {0, 5},
		},
		Cols: [5]model.Hint{
			{0, 5}, {0, 5}, {0, 5}, {0, 5}, {0, 5},
		},
	}
	got := g.GetHints()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetHintsBombs(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 1, 1, 1, 1,
		1, 0, 1, 1, 1,
		0, 1, 0, 1, 1,
		1, 0, 1, 1, 1,
		1, 1, 0, 1, 1,
	}}

	want := model.Hints{
		Rows: [5]model.Hint{
			{1, 4}, {1, 4}, {2, 3}, {1, 4}, {1, 4},
		},
		Cols: [5]model.Hint{
			{2, 3}, {2, 3}, {2, 3}, {0, 5}, {0, 5},
		},
	}
	got := g.GetHints()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetHintsPoints(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 1, 1, 1, 2,
		1, 0, 1, 3, 1,
		0, 4, 0, 1, 3,
		1, 0, 2, 1, 1,
		2, 1, 0, 1, 1,
	}}

	want := model.Hints{
		Rows: [5]model.Hint{
			{1, 5}, {1, 6}, {2, 8}, {1, 5}, {1, 5},
		},
		Cols: [5]model.Hint{
			{2, 4}, {2, 6}, {2, 4}, {0, 7}, {0, 8},
		},
	}
	got := g.GetHints()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetMaxScoreZero(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
	}}

	want := 0
	got := g.GetMaxScore()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetMaxScoreOne(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 1,
	}}

	want := 1
	got := g.GetMaxScore()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_GetMaxScoreGeneral(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		1, 0, 0, 0, 0,
		0, 2, 0, 4, 0,
		0, 3, 1, 0, 0,
		0, 0, 4, 1, 0,
		0, 0, 0, 0, 1,
	}}

	want := 96
	got := g.GetMaxScore()

	if want != got {
		t.Error("expected", want, "got", got)
	}
}
