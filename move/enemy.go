package move

import "github.com/itzamna314/battlesnake/model"

// Calculate weight for moving You to coord in state
func weightEnemies(state *model.GameState, coord *model.Coord) float64 {
	// No enemies
	if len(state.Board.Snakes) <= 1 {
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

		// Avoid with weight of negative probability
		weight -= prob
	}

	return weight / float64(len(state.Board.Snakes)-1)
}
