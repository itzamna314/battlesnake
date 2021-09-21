package move

import (
	"github.com/itzamna314/battlesnake/model"
)

const (
	Death = -1.0
	Base  = 0.5
	Avoid = 0.0
	Food  = 0.15
)

func moveWeight(state *model.GameState, coord *model.Coord) float64 {
	if !isSafe(state, coord) {
		return Death
	}

	// Compute food weight
	weight := Base
	weight += weightFood(state, coord)

	// Don't die over food
	if weight < Avoid {
		return Avoid
	}

	return weight
}
