package game

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

func (b *Board) Clone() Board {
	clone := Board{
		Height: b.Height,
		Width:  b.Width,
		Food:   make([]Coord, len(b.Food)),
		Snakes: make([]Battlesnake, len(b.Snakes)),
	}

	for i, f := range b.Food {
		clone.Food[i] = f
	}

	for i, s := range b.Snakes {
		clone.Snakes[i] = s.Clone()
	}

	return clone
}
