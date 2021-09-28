package tree_test

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
)

type testBrain struct {
	weightFunc func(*game.Coord, *game.Battlesnake) float64
	abortFunc  func(float64) bool
}

func (b *testBrain) Init(*game.GameState) {}
func (b *testBrain) Clone() tree.SnakeBrain {
	return b
}

// testBrain moves to the right
func (b *testBrain) Weight(coord *game.Coord, snake *game.Battlesnake) float64 {
	if b.weightFunc == nil {
		panic("called testBrain.Weight with nil weightFunc")
	}

	return b.weightFunc(coord, snake)
}

// testBrain never aborts
func (b *testBrain) Abort(weight float64) bool {
	if b.abortFunc == nil {
		return false
	}

	return b.abortFunc(weight)
}

func (b *testBrain) Move(*game.Battlesnake, game.Direction) {}
func (b *testBrain) MoveEnemies(*game.Battlesnake)          {}
