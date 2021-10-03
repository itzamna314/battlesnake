package predict

import (
	"github.com/itzamna314/battlesnake/game"
)

type Board struct {
	Height int          `json:"height"`
	Width  int          `json:"width"`
	Food   []game.Coord `json:"food"`
	Snakes []Snake      `json:"snakes"`

	// Used in non-standard game modes
	Hazards []game.Coord `json:"hazards"`
}

func (b *Board) Clone() Board {
	clone := Board{
		Height:  b.Height,
		Width:   b.Width,
		Food:    make([]game.Coord, len(b.Food)),
		Hazards: make([]game.Coord, len(b.Hazards)),
		Snakes:  make([]Snake, len(b.Snakes)),
	}

	for i, f := range b.Food {
		clone.Food[i] = f
	}

	for i, h := range b.Hazards {
		clone.Hazards[i] = h
	}

	for i, s := range b.Snakes {
		clone.Snakes[i] = s.Clone()
	}

	return clone
}

func (b *Board) Init(gb *game.Board) {
	b.Height = gb.Height
	b.Width = gb.Width

	b.Food = make([]game.Coord, len(gb.Food))
	for i := range gb.Food {
		b.Food[i] = gb.Food[i]
	}

	b.Snakes = make([]Snake, len(gb.Snakes))
	for i := range gb.Snakes {
		b.Snakes[i].Init(&gb.Snakes[i])
	}

	b.Hazards = make([]game.Coord, len(gb.Hazards))
	for i := range gb.Hazards {
		b.Hazards[i] = gb.Hazards[i]
	}
}
