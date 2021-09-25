package guess

import (
	"fmt"

	"github.com/itzamna314/battlesnake/game"
)

const (
	Certain    = 1.0
	Impossible = 0.0
)

type Coord struct {
	game.Coord
	Probability float64
}

func (g *Coord) String() string {
	return fmt.Sprintf("%s[%4f]", g.Coord, g.Probability)
}
