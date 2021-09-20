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
	tree := BuildTree(&state, 1)

	log.Println("Choosing from possible moves:")
	for dir, move := range tree.Root.Moves {
		if move == nil {
			continue
		}

		fmt.Printf("%s: Safe? %t, Weight: %f, Shout: %s\n",
			model.Direction(dir),
			move.Safe,
			move.Weight,
			move.Shout,
		)
	}

	var (
		nextMove  string
		nextShout string
		safeMoves []string
		maxWeight float64
	)
	for dir, coord := range tree.Root.Moves {
		if !coord.Safe {
			continue
		}

		if coord.Weight > maxWeight {
			maxWeight = coord.Weight
			safeMoves = []string{model.Direction(dir).String()}
			nextShout = coord.Shout
		} else if coord.Weight == maxWeight {
			safeMoves = append(safeMoves, model.Direction(dir).String())
			nextShout = ""
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		nextShout = "bye"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}

	return model.BattlesnakeMoveResponse{
		Move:  nextMove,
		Shout: nextShout,
	}
}
