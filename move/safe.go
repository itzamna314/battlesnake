package move

import "github.com/itzamna314/battlesnake/model"

// Calculate weight for moving You to coord in state
func weightSafe(state *model.GameState, coord *model.Coord) float64 {
	// Don't hit walls.
	if coord.X < 0 {
		return WeightDeath
	} else if coord.X >= state.Board.Width {
		return WeightDeath
	}

	if coord.Y < 0 {
		return WeightDeath
	} else if coord.Y >= state.Board.Height {
		return WeightDeath
	}

	// Don't hit snakes
	for _, snake := range state.Board.Snakes {
		for _, sBody := range snake.Body {
			if coord.Hit(&sBody) {
				return WeightDeath
			}
		}
	}

	return WeightNothing
}
