package model

type Game struct {
	ID    string
	Board [25]byte
	Moves []int
}

type Hint struct {
	Mines  int
	Points int
}

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

func (g *Game) GetScore() int {
	s := 1
	for _, m := range g.Moves {
		s *= int(g.Board[m])
	}
	return s
}
