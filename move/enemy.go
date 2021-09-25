package move

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
)

// Calculate weight for moving You to coord in state
func weightEnemies(state *predict.State, coord *game.Coord) float64 {
	// No enemies
	if len(state.Board.Snakes) <= 1 {
		return Nothing
	}

	var weight float64
	for i, enemy := range state.Board.Snakes {
		if enemy.ID == state.You.ID {
			continue
		}

		prob := state.BodyGuesses[i].Prob(coord)
		if prob == Certain {
			return Death
		}
		weight -= prob

		// If we are shorter, avoid with weight of probability
		// Otherwise, attack with reduced weight of collision probability
		// STRIKE FIRST STRIKE HARD NO MERCY
		// But also don't chase a short snake into a long snake
		prob = state.HeadGuesses[i].Prob(coord)
		if enemy.Length >= state.You.Length {
			weight -= prob
		} else {
			weight += (prob * 0.3)
		}
	}

	return weight
}
