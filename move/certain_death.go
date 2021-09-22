package move

import "github.com/itzamna314/battlesnake/model"

func isCertainDeath(state *model.GameState, coord *model.Coord) bool {
	if coord.X < 0 {
		return true
	} else if coord.X >= state.Board.Width {
		return true
	}

	if coord.Y < 0 {
		return true
	} else if coord.Y >= state.Board.Height {
		return true
	}

	// Don't hit self
	for _, body := range state.You.Body {
		if coord.Hit(&body) {
			return true
		}
	}

	return false
}
