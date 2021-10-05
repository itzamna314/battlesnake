package multiverse

import (
	"fmt"
	"math"

	"github.com/itzamna314/battlesnake/tree"
)

const (
	// We start evaluating moves at 0.5
	Base FloatWeight = 0.5
	// Lowest-possible score we like
	Avoid FloatWeight = 0.0

	// From modifier funcs, this means no modifier
	Nothing FloatWeight = 0.0

	// Base value of food
	Food FloatWeight = 0.15
	// Food value when we want to avoid eating
	FoodAvoid FloatWeight = -0.1

	// Value of enemies that can kill us
	EnemyAvoid FloatWeight = -4.0
	// Value of enemies we tie (both die) with
	EnemyTie FloatWeight = -1.0
	// Value of enemies that we can kill
	EnemyKill FloatWeight = 0.3
)

var (
	// Certain death is -infinity
	Death = FloatWeight(math.Inf(-1))
)

type FloatWeight float64

func (w FloatWeight) Compare(other tree.Weight) int8 {
	o := other.(FloatWeight)
	switch {
	case w == o:
		return 0
	case w < o:
		return -1
	default:
		return 1
	}
}

func (w FloatWeight) String() string {
	return fmt.Sprintf("%.2f", w)
}

func (s *State) Weight(nd *tree.Node) tree.Weight {
	snake := s.Snake(nd.SnakeID)

	if SnakeWillDie(s, &snake.Head, snake) {
		return Death
	}

	var weight FloatWeight
	if nd.Parent != nil && nd.Parent.Weight != nil {
		weight = nd.Parent.Weight.(FloatWeight)
	}

	// Start with our remaining health
	health := s.weightHealth(snake)
	weight += health

	// Weight for enemy encounters
	enemy := s.weightEnemies(snake)
	if enemy <= Death {
		return Death
	}
	weight += enemy

	// Weight food
	food := s.weightFood(snake)
	weight += food

	return weight
}

func (s *State) Abort(nd *tree.Node) bool {
	return nd.Weight.Compare(Death) <= 0
}
