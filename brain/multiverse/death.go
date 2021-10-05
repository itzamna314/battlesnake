package multiverse

import "github.com/itzamna314/battlesnake/game"

func YouWillDie(state *State) bool {
	return SnakeWillDie(state, &state.You.Head, &state.You)
}

func SnakeWillDie(state *State, coord *game.Coord, snake *Snake) bool {
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
	for i := 1; i < len(snake.Body); i++ {
		if coord.Hit(&snake.Body[i]) {
			return true
		}
	}

	// Don't starve
	if snake.Health <= 0 {
		return true
	}

	return false
}
