package main

// Encapsulates row and column positions
type Pos struct {
	Row int
	Col int
}

// Memo map for valid moves
var moves map[Pos][]Pos = make(map[Pos][]Pos)

// All valid knight moves, relative to a position
var knightMoves []Pos = []Pos{Pos{-2, -1}, Pos{-1, -2}, Pos{1, -2}, Pos{2, -1}, Pos{2, 1}, Pos{1, 2}, Pos{-1, 2}, Pos{-2, 1}}

// Returns a list of valid moves for a given position, pruning ones not on the board
// moves are relative [r, c] coordinates
func (p Pos) GetMoves() (results []Pos) {
	var ok bool
	results, ok = moves[p]
	if ok {
		return
	}
	for _, pos := range knightMoves {
		if p.IsMoveInGrid(pos) {
			results = append(results, p.Add(pos))
		}
	}
	moves[p] = results
	return
}

// Checks whether or not a given move will be inside the grid
func (p Pos) IsMoveInGrid(move Pos) bool {
	new := p.Add(move)
	return new.Row > 0 && new.Col > 0 && new.Row < len(Grid) && new.Col < len(Grid[0])
}

// Adds two positions together, just like adding vectors
func (p Pos) Add(other Pos) Pos {
	return Pos{p.Row + other.Row, p.Col + other.Col}
}
