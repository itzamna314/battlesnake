package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/itzamna314/battlesnake/model"
	"github.com/itzamna314/battlesnake/tree"
)

func Move(w http.ResponseWriter, r *http.Request) {
	state := model.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	response := NextMove(state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func NextMove(state model.GameState) model.BattlesnakeMoveResponse {
	moveTree := tree.Build(&state, 1)

	log.Println("Choosing from possible moves:")
	var (
		nextMove  string
		nextShout string
		bestMoves []string
		maxWeight float64
	)
	for dir, coord := range moveTree.Children {
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
