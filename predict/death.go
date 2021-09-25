package predict

import (
	"github.com/itzamna314/battlesnake/game"
)

func YouWillDie(state *State, coord *game.Coord) bool {
	return SnakeWillDie(state, coord, &state.You)
}

func SnakeWillDie(state *State, coord *game.Coord, snake *game.Battlesnake) bool {
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
	for _, body := range snake.Body {
		if coord.Hit(&body) {
			return true
		}
	}

	return false
}
