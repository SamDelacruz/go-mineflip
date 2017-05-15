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
