package move

import "github.com/itzamna314/battlesnake/model"

// Calculate weight for moving You to coord in state
func weightSafety(state *model.GameState, coord *model.Coord) float64 {
	// Don't hit walls.
	if coord.X < 0 {
		return Death
	} else if coord.X >= state.Board.Width {
		return Death
	}

	if coord.Y < 0 {
		return Death
	} else if coord.Y >= state.Board.Height {
		return Death
	}

	// Don't hit snakes
	for _, snake := range state.Board.Snakes {
		isEnemy := snake.ID != state.You.ID

		for _, sBody := range snake.Body {
			if coord.Hit(&sBody) {
				return Death
			}
		}

		// If this is an enemy that can eat us, avoid its next moves
		if isEnemy && snake.Length >= state.You.Length {
			opts := model.Options(&snake.Head)
			for _, opt := range opts {
				if coord.Hit(&opt.Coord) {
					// assume enemy snake has 3 valid moves
					// this could be improved
					return Death * 0.33
				}
			}
		}
	}

	return Nothing
}
