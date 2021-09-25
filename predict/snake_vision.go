package predict

import (
	"github.com/itzamna314/battlesnake/guess"
)

type SnakeVision []guess.CoordSet

func (v SnakeVision) Clone() SnakeVision {
	clone := make(SnakeVision, len(v))
	for x := 0; x < len(v); x++ {
		clone[x] = v[x].Clone()
	}

	return clone
}
