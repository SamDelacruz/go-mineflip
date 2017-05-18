package model

import (
	"encoding/json"
	"errors"
)

// Game stores the current state of a game
// All properties can be derived from this struct
// i.e score, win/lose, visible tiles, hints
type Game struct {
	ID    string
	Board [25]byte
	Moves []int
}

// GameRepr is the public derived state of a game
type GameRepr struct {
	ID    string       `json:"id"`
	Board [5][5]string `json:"board"`
	Score int          `json:"score"`
	Hints Hints        `json:"hints"`
	Won   bool         `json:"game_won"`
	Lost  bool         `json:"game_lost"`
}

// ToRepr returns the public derived state of a game
func (g *Game) ToRepr() GameRepr {
	var repr GameRepr
	repr.ID = g.ID
	repr.Score = g.GetScore()
	repr.Hints = g.GetHints()
	repr.Board = unsquashBoard(g.GetVisible())
	repr.Won = g.Won()
	repr.Lost = g.Lost()
	return repr
}

// Hint stores the hint metadata for a row/column
// Mines is the number of mines/zeroes in a row/column
// Points is the sum of a row/column
type Hint struct {
	Mines  int `json:"mines"`
	Points int `json:"points"`
}

// Hints is the collection type for grouping hints
// by rows and columns
type Hints struct {
	Rows [5]Hint `json:"rows"`
	Cols [5]Hint `json:"cols"`
}

// ToJSON returns the json representation of a Game
func (g *Game) ToJSON() ([]byte, error) {
	return json.Marshal(g.ToRepr())
}

func unsquashBoard(b string) [5][5]string {
	var s [5][5]string
	bb := []byte(b)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			s[j][i] = string(bb[j*5+i])
		}
	}
	return s
}

// GetVisible gets the string representation of a gameboard
// Unrevealed tiles represented as '?'
// Revealed tiles represented by their points value
// Mines (0) represented by "*"
func (g *Game) GetVisible() string {
	v := []byte{
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
	}

	if g.Lost() {
		for i, t := range g.Board {
			v[i] = byteToTile(t)
		}
	} else {
		for _, m := range g.Moves {
			v[m] = byteToTile(g.Board[m])
		}
	}
	return string(v)
}

func byteToTile(b byte) byte {
	if b == 0 {
		return '*'
	} // Mine is a *
	return b + 48 // ASCII numeric offset - fine for 0-9
}

// GetScore returns the current score for a game
// Score is multiplicative - the product of all revealed tiles
// Edge case: when no tiles have been revealed, score is zero.
func (g *Game) GetScore() int {
	if len(g.Moves) == 0 {
		return 0
	}
	s := 1
	for _, m := range g.Moves {
		s *= int(g.Board[m])
	}
	return s
}

// GetMaxScore returns the maximum number of points available
// for the given game board
func (g *Game) GetMaxScore() int {
	t := 0
	for _, v := range g.Board {
		if v > 0 {
			if t == 0 {
				t++
			}
			t *= int(v)
		}
	}
	return t
}

// GetHints returns the set of hints for a game
func (g *Game) GetHints() Hints {
	var h Hints
	for i := 0; i < 5; i++ {
		h.Cols[i] = g.getColHint(i)
		h.Rows[i] = g.getRowHint(i)
	}
	return h
}

func (g *Game) getColHint(col int) Hint {
	m, p := 0, 0
	for i := col; i < len(g.Board); i += 5 {
		v := g.Board[i]
		if v == 0 {
			m++
		}
		p += int(v)
	}
	return Hint{m, p}
}

func (g *Game) getRowHint(row int) Hint {
	m, p := 0, 0
	for i := 5 * row; i < 5*row+5; i++ {
		v := g.Board[i]
		if v == 0 {
			m++
		}
		p += int(v)
	}
	return Hint{m, p}
}

func (g *Game) Lost() bool {
	return len(g.Moves) > 0 && g.GetScore() == 0
}

func (g *Game) Won() bool {
	return g.GetScore() == g.GetMaxScore()
}

func (g *Game) AddMove(i int) (byte, error) {
	if g.moveAvail(i) {
		g.Moves = append(g.Moves, i)
		return g.Board[i], nil
	}
	return 0, errors.New("INVALID_MOVE")
}

func (g *Game) moveAvail(i int) bool {
	for _, m := range g.Moves {
		if m == i {
			return false
		}
	}
	return true
}
