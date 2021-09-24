package move

import "github.com/itzamna314/battlesnake/game"

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
	// Value of hazard
	Hazard = -0.3

	// Probabilities
	Certain = 1.0
)

func Weight(state *game.GameState, coord *game.Coord) float64 {
	if game.YouWillDie(state, coord) {
		return Death
	}

	// Compute food weight
	weight := Base

	enemy := weightEnemies(state, coord)

	// Nothing more certain than Death
	if enemy <= Death {
		return Death
	}

	weight += enemy

	food := WeightFood(state, coord)

	weight += food

	// Don't die over food
	if weight < Avoid {
		return Avoid
	}

	return weight
}
