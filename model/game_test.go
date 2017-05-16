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

func TestGame_GetVisibleMines(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{0, 24}}

	want := "*???????????????????????2"
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

func TestGame_ToJSON(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{10, 24}}

	want := `{"id":"1","board":[["?","?","?","?","?"],["?","?","?","?","?"],["1","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","2"]],"score":2,"hints":{"rows":[{"mines":1,"points":6},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":6}],"cols":[{"mines":1,"points":4},{"mines":0,"points":7},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":6}]},"game_won":false,"game_lost":false}`
	j, _ := g.ToJSON()
	got := string(j)

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_ToJSONLose(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 2,
	}, Moves: []int{0, 24}}

	want := `{"id":"1","board":[["*","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","2"]],"score":0,"hints":{"rows":[{"mines":1,"points":6},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":6}],"cols":[{"mines":1,"points":4},{"mines":0,"points":7},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":6}]},"game_won":false,"game_lost":true}`
	j, _ := g.ToJSON()
	got := string(j)

	if want != got {
		t.Error("expected", want, "got", got)
	}
}

func TestGame_ToJSONWin(t *testing.T) {
	g := model.Game{ID: "1", Board: [25]byte{
		0, 3, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
	}, Moves: []int{1, 24}}

	want := `{"id":"1","board":[["?","3","?","?","?"],["?","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","?"],["?","?","?","?","1"]],"score":3,"hints":{"rows":[{"mines":1,"points":6},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":5}],"cols":[{"mines":1,"points":4},{"mines":0,"points":7},{"mines":0,"points":5},{"mines":0,"points":5},{"mines":0,"points":5}]},"game_won":true,"game_lost":false}`
	j, _ := g.ToJSON()
	got := string(j)

	if want != got {
		t.Error("expected", want, "got", got)
	}
}
