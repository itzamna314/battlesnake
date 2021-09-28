package tree

import (
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
	Weight(*game.Coord, *game.Battlesnake) float64
	// Abort takes a weight, and returns true if we can abort the search
	// based on that weight
	Abort(float64) bool
	// Move moves the indicated snake. The brain should track that move
	// in its internal data structure
	Move(*game.Battlesnake, game.Direction)
	// MoveEnemies moves all of the indicated snake's enemies. The brain
	// should predict those moves in its internal data structure
	// Since we don't know for sure where the enemies are going,
	// no direction is provided
	MoveEnemies(*game.Battlesnake)
}
