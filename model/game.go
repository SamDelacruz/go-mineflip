package model

// Game stores the current state of a game
// All properties can be derived from this struct
// i.e score, win/lose, visible tiles, hints
type Game struct {
	ID    string
	Board [25]byte
	Moves []int
}

// Hint stores the hint metadata for a row/column
// Mines is the number of mines/zeroes in a row/column
// Points is the sum of a row/column
type Hint struct {
	Mines  int
	Points int
}

// Hints is the collection type for grouping hints
// by rows and columns
type Hints struct {
	Rows [5]Hint
	Cols [5]Hint
}

// GetVisible gets the string representation of a gameboard
// Unrevealed tiles represented as '?'
// Revealed tiles represented by their points value
func (g *Game) GetVisible() string {
	v := []byte{
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
		'?', '?', '?', '?', '?',
	}

	for _, m := range g.Moves {
		v[m] = g.Board[m] + 48 // ASCII numeric offset - fine for 0-9
	}
	return string(v)
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
