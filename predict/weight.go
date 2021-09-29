package predict

import "github.com/itzamna314/battlesnake/game"

const (
	// Certain death is -1
	Death = -1.0
	// We start evaluating moves at 0.5
	Base = 0.5
	// Lowest-possible score we like
	Avoid = 0.0

	// From modifier funcs, this means no modifier
	Nothing = 0.0

	// Base value of food
	Food = 0.15
	// Food value when we want to avoid eating
	FoodAvoid = -0.0375

	// Value of hazard
	Hazard = -0.3
	// Required to live
	Mandatory = 1.0
)

func (s *State) Weight(coord *game.Coord, snake *game.Battlesnake) float64 {
	if SnakeWillDie(s, coord, snake) {
		return Death
	}

	// Compute food weight
	weight := Base

	enemy := s.weightEnemies(coord, snake)
	if enemy <= Death {
		return Death
	}
	weight += enemy

	hazard := s.weightHazard(coord, snake)
	weight += hazard
	if weight <= Death {
		return Death
	}

	food := s.weightFood(coord, snake)
	weight += food

	// Don't die over food
	if weight < Avoid {
		return Avoid
	} else if weight > Mandatory {
		return Mandatory
	}

	return weight
}

func (s *State) Abort(weight float64) bool {
	return weight <= Death
}
