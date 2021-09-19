package move

import (
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/model"
)

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func Next(state model.GameState) model.BattlesnakeMoveResponse {
	myHead := state.You.Body[0] // Coordinates of your head
	possibleMoves := model.Options(&myHead)

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	if myNeck.X < myHead.X {
		possibleMoves[model.Left].Safe = false
	} else if myNeck.X > myHead.X {
		possibleMoves[model.Right].Safe = false
	} else if myNeck.Y < myHead.Y {
		possibleMoves[model.Down].Safe = false
	} else if myNeck.Y > myHead.Y {
		possibleMoves[model.Up].Safe = false
	}

	// TODO: Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	// boardWidth := state.Board.Width
	// boardHeight := state.Board.Height
	if myHead.X-1 < 0 {
		possibleMoves[model.Left].Safe = false
	} else if myHead.X+1 == state.Board.Width {
		possibleMoves[model.Right].Safe = false
	}

	if myHead.Y-1 < 0 {
		possibleMoves[model.Down].Safe = false
	} else if myHead.Y+1 == state.Board.Height {
		possibleMoves[model.Up].Safe = false
	}

	// TODO: Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	//mybody := state.You.Body

	// TODO: Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := []string{}
	for dir, coord := range possibleMoves {
		if coord.Safe {
			safeMoves = append(safeMoves, model.Direction(dir).String())
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return model.BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
