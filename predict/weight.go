package predict

import (
	"math"

	"github.com/itzamna314/battlesnake/game"
)

const (
	// We start evaluating moves at 0.5
	Base = 0.5
	// Lowest-possible score we like
	Avoid = 0.0

	// From modifier funcs, this means no modifier
	Nothing = 0.0

	// Base value of food
	Food = 0.15
	// Food value when we want to avoid eating
	FoodAvoid = -0.1

	// Value of enemies that can kill us
	EnemyAvoid = -4.0
	// Value of enemies we tie (both die) with
	EnemyTie = -1.0
	// Value of enemies that we can kill
	EnemyKill = 0.3

	// Value of hazard
	Hazard = -0.2
	// Required to live
	Mandatory = 1.0
)

var (
	// Certain death is -infinity
	Death = math.Inf(-1)
)

func (s *State) Weight(coord *game.Coord, snakeID string) float64 {
	snake := s.Snake(snakeID)

	if SnakeWillDie(s, coord, snake) {
		return Death
	}

	var weight float64
	// Start with our remaining health
	health := s.weightHealth(snake)
	weight += health

	// Weight for enemy encounters
	enemy := s.weightEnemies(coord, snake)
	if enemy <= Death {
		return Death
	}
	weight += enemy

	// Weight food
	food := s.weightFood(coord, snake)
	weight += food

	// Clip to a maximum to prevent runaway scores due to bugs
	if weight > Mandatory {
		return Mandatory
	}

	return weight
}

func (s *State) Abort(weight float64) bool {
	return weight <= Death
}
