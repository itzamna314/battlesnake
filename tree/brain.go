package tree

import (
	"fmt"

	"github.com/itzamna314/battlesnake/game"
)

type SnakeBrain interface {
	// Init must initialize the brain's data structures
	// It will be called once with an initial game state
	Init(*game.GameState)
	// Clone returns a deep copy of this SnakeBrain
	Clone() SnakeBrain
	// Weight calculates the weight for the indicated snake index,
	// based on the brain's internal data
	Weight(*Node) Weight
	// Abort takes a weight, and returns true if we can abort the search
	// based on that weight
	Abort(*Node) bool
	// Move moves the indicated snake. The brain should track that move
	// in its internal data structure
	Move(string, game.Direction)
	// MoveEnemies moves all of the indicated snake's enemies. The brain
	// should predict those moves in its internal data structure
	// Since we don't know for sure where the enemies are going,
	// no direction is provided
	MoveEnemies(string)
}

// Weight represents the value that a brain can assign to a move
type Weight interface {
	// Compare compares this Weight to other.
	// It returns < 0 if we are less, > 0 if greater, and 0 if equal
	Compare(other Weight) int8
	fmt.Stringer
}
