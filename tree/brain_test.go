package tree_test

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
)

type testBrain struct {
	weightFunc func(*game.Coord, string) float64
	abortFunc  func(float64) bool
}

func (b *testBrain) Init(*game.GameState) {}
func (b *testBrain) Clone() tree.SnakeBrain {
	return b
}

// testBrain moves to the right
func (b *testBrain) Weight(coord *game.Coord, snakeID string) float64 {
	if b.weightFunc == nil {
		return 0
	}

	return b.weightFunc(coord, snakeID)
}

// testBrain never aborts
func (b *testBrain) Abort(weight float64) bool {
	if b.abortFunc == nil {
		return false
	}

	return b.abortFunc(weight)
}

func (b *testBrain) Move(string, game.Direction) {}
func (b *testBrain) MoveEnemies(string)          {}
