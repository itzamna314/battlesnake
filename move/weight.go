package move

import (
	"fmt"

	"github.com/itzamna314/battlesnake/model"
)

const (
	// Certain death is -1
	Death = -1.0
	// We start evaluating moves at 0.5
	Base = 0.5
	// Lowest-possible score that doesn't lead to death
	Avoid = 0.0

	// From modifier funcs, this means no modifier
	Nothing = 0.0
	// Value of food
	Food = 0.15
)

func moveWeight(state *model.GameState, coord *model.Coord) (w float64) {
	defer func() {
		fmt.Printf("Calculated weight for coord %s: %v\n", coord, w)
	}()

	if isCertainDeath(state, coord) {
		fmt.Printf("Certain death at %s\n", coord)
		return Death
	}

	// Compute food weight
	weight := Base

	weight += weightEnemies(state, coord)

	fmt.Printf("Calculated enemy weight %s: %v\n", coord, weight)

	// Nothing more certain than Death
	if weight <= Death {
		return Death
	}

	weight += weightFood(state, coord)

	fmt.Printf("Calculated food weight %s: %v\n", coord, weight)

	// Don't die over food
	if weight < Avoid {
		return Avoid
	}

	return weight
}
