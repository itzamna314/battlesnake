package move

import (
	"github.com/itzamna314/battlesnake/model"
)

func safe(state model.GameState, possible model.PossibleMoves) {
	var (
		myBody = state.You.Body
		myHead = myBody[0]
	)

	// Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	if myHead.X-1 < 0 {
		possible[model.Left].Safe = false
	} else if myHead.X+1 == state.Board.Width {
		possible[model.Right].Safe = false
	}

	if myHead.Y-1 < 0 {
		possible[model.Down].Safe = false
	} else if myHead.Y+1 == state.Board.Height {
		possible[model.Up].Safe = false
	}

	// Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	for dir, poss := range possible {
		if !poss.Safe {
			continue
		}

		for _, body := range myBody {
			if poss.Coord.Hit(&body) {
				possible[dir].Safe = false
			}
		}
	}

	// Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.
Enemies:
	for _, enemy := range state.Board.Snakes {
		if enemy.ID == state.You.ID {
			continue
		}
		for eIdx, eBody := range enemy.Body {
			isHead := eIdx == 0
			for dir, poss := range possible {
				if !poss.Safe {
					continue
				}

				if poss.Hit(&eBody) {
					if isHead && enemy.Length < state.You.Length {
						possible[dir].Safe = true
						possible[dir].Weight = 1
						possible[dir].Shout = "KILL"
						break Enemies
					}
					possible[dir].Safe = false
				}
			}
		}
	}
}
