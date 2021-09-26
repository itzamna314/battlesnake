package move

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
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
	// Base value of food
	Food = 0.15
	// Value of hazard
	Hazard = -0.2
	// Required to live
	Mandatory = 1.0

	// Probabilities
	Certain = 1.0
)

func Weight(state *predict.State, coord *game.Coord) float64 {
	if predict.YouWillDie(state, coord) {
		return Death
	}

	// Compute food weight
	weight := Base

	enemy := weightEnemies(state, coord)
	if enemy <= Death {
		return Death
	}
	weight += enemy

	hazard := WeightHazard(state, coord)
	weight += hazard
	if weight <= Death {
		return Death
	}

	food := WeightFood(state, coord)
	weight += food

	// Don't die over food
	if weight < Avoid {
		return Avoid
	}

	return weight
}
