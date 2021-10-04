package tree_test

import (
	"fmt"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
)

type weight int32

func (w weight) Compare(other tree.Weight) int8 {
	o := other.(weight)
	switch {
	case w == o:
		return 0
	case w < o:
		return -1
	default:
		return 1
	}
}

func (w weight) String() string {
	return fmt.Sprintf("%d", w)
}

type testBrain struct {
	weightFunc func(*tree.Node) tree.Weight
	abortFunc  func(*tree.Node) bool
}

func (b *testBrain) Init(*game.GameState) {}
func (b *testBrain) Clone() tree.SnakeBrain {
	return b
}

// testBrain moves to the right
func (b *testBrain) Weight(nd *tree.Node) tree.Weight {
	if b.weightFunc == nil {
		return weight(0)
	}

	if nd.Parent.Weight == nil {
		return b.weightFunc(nd)
	}

	return b.weightFunc(nd).(weight) + nd.Parent.Weight.(weight)
}

// testBrain never aborts
func (b *testBrain) Abort(nd *tree.Node) bool {
	if b.abortFunc == nil {
		return false
	}

	return b.abortFunc(nd)
}

func (b *testBrain) Move(string, game.Direction) {}
func (b *testBrain) MoveEnemies(string)          {}
