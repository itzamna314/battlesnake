package move

import (
	"fmt"

	"github.com/itzamna314/battlesnake/model"
)

// Calculate weight for moving You to coord in state
func weightEnemies(state *model.GameState, coord *model.Coord) float64 {
	// No enemies
	if len(state.Board.Snakes) <= 1 {
		fmt.Println("no enemies")
		return Nothing
	}

	var weight float64
	for i, enemy := range state.Board.Snakes {
		if enemy.ID == state.You.ID {
			continue
		}

		prob := state.EnemyGuesses[i].Prob(coord)
		if prob == Certain {
			return Death
		}

		// Avoid with weight of twice negative probability
		weight -= prob
	}

	// Double weight to increase enemy avoid priority
	return weight * 2
}
