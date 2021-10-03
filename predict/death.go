package predict

import "github.com/itzamna314/battlesnake/game"

func YouWillDie(state *State, coord *game.Coord) bool {
	return SnakeWillDie(state, coord, &state.You)
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
	for i, body := range snake.Body {
		// Tail - it will move as long as we didn't eat last turn
		if i == len(snake.Body)-1 && snake.Health < 100 {
			break
		}

		if coord.Hit(&body) {
			return true
		}
	}

	// Don't starve
	if snake.Health <= 0 {
		return true
	}

	return false
}
