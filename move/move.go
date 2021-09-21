package move

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/model"
)

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func Next(state model.GameState) model.BattlesnakeMoveResponse {
	tree := BuildTree(&state, 10)

	log.Println("Choosing from possible moves:")
	var (
		nextMove  string
		nextShout string
		bestMoves []string
		maxWeight float64
	)
	for dir, coord := range tree.Children {
		if coord == nil {
			continue
		}

		fmt.Printf("%s: Weight: %f\n",
			model.Direction(dir),
			coord.Weight,
		)

		if coord.Weight > maxWeight {
			maxWeight = coord.Weight
			bestMoves = []string{model.Direction(dir).String()}
		} else if coord.Weight == maxWeight {
			bestMoves = append(bestMoves, model.Direction(dir).String())
			nextShout = ""
		}
	}

	if len(bestMoves) == 0 {
		nextMove = "down"
		nextShout = "bye"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = bestMoves[rand.Intn(len(bestMoves))]
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}

	return model.BattlesnakeMoveResponse{
		Move:  nextMove,
		Shout: nextShout,
	}
}
