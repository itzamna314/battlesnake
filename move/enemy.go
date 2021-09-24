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

		// If we are shorter, avoid with weight of probability
		// Otherwise, attack with weight of collision probability
		// STRIKE FIRST STRIKE HARD NO MERCY
		if enemy.Length >= state.You.Length {
			weight -= prob
		} else {
			weight += prob
		}
	}

	return weight
}
