package move

import "github.com/itzamna314/battlesnake/model"

// Calculate weight for moving You to coord in state
func isSafe(state *model.GameState, coord *model.Coord) bool {
	// Don't hit walls.
	if coord.X < 0 {
		return false
	} else if coord.X >= state.Board.Width {
		return false
	}

	if coord.Y < 0 {
		return false
	} else if coord.Y >= state.Board.Height {
		return false
	}

	// Don't hit snakes
	for _, snake := range state.Board.Snakes {
		for _, sBody := range snake.Body {
			if coord.Hit(&sBody) {
				return false
			}
		}
	}

	return true
}
