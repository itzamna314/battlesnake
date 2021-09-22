package move

import "github.com/itzamna314/battlesnake/model"

// Calculate weight for moving You to coord in state
func weightEnemies(state *model.GameState, coord *model.Coord) float64 {
	vis := state.Future[coord.X][coord.Y]

	// No enemies
	if len(state.Board.Snakes) <= 1 {
		return Nothing
	}

	var weight float64
	for i, enemy := range state.Board.Snakes {
		if enemy.ID == state.You.ID {
			continue
		}

		// Ignore low-probability enemies
		if vis.Enemies[i] < 0.75 {
			continue
		}

		// Avoid with weight of negative probability
		if enemy.Length >= state.You.Length {
			weight += -1 * vis.Enemies[i]
		}
	}

	return weight / float64(len(state.Board.Snakes)-1)
}
