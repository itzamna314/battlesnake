package move

import (
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/game"
)

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func Next(state game.GameState) game.Direction {
	moveTree := BuildTree(&state, 6)

	log.Println("Choosing from possible moves:")
	var (
		bestMoves []game.Direction
		maxWeight float64
	)
	for dir, coord := range moveTree.Children {
		if coord == nil {
			continue
		}

		log.Printf("%s: Weight: %f\n",
			game.Direction(dir),
			coord.Weight,
		)

		if coord.Weight > maxWeight {
			maxWeight = coord.Weight
			bestMoves = []game.Direction{game.Direction(dir)}
		} else if coord.Weight == maxWeight {
			bestMoves = append(bestMoves, game.Direction(dir))
		}
	}

	// If we have no best moves, admit defeat
	if len(bestMoves) == 0 {
		log.Printf("%s MOVE %d: No safe moves detected! Moving down\n", state.Game.ID, state.Turn)
		return game.Down
	}

	var nextMove game.Direction

	// If we have a single best move, do it
	if len(bestMoves) == 1 {
		nextMove = bestMoves[0]
	} else {
		nextMove = bestMoves[rand.Intn(len(bestMoves))]
	}

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	return nextMove
}
