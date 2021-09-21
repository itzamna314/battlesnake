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
		isEnemy := snake.ID != state.You.ID

		for _, sBody := range snake.Body {
			if coord.Hit(&sBody) {
				return false
			}
		}

		// If this is an enemy that can eat us, avoid its next moves
		if isEnemy && snake.Length > state.You.Length {
			opts := model.Options(&snake.Head)
			for _, opt := range opts {
				if coord.Hit(&opt.Coord) {
					return false
				}
			}
		}
	}

	return true
}
