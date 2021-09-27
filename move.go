package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
	"github.com/itzamna314/battlesnake/tree"
)

type BattlesnakeMoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

func Move(w http.ResponseWriter, r *http.Request) {
	state := game.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	// Read timeout from flag
	to := time.Duration(*timeout)
	ctx, _ := context.WithTimeout(r.Context(), to*time.Millisecond)

	direction := tree.Search(ctx, &state, &state.You, &predict.State{})

	log.Printf("MOVE: %s\n", direction)

	response := BattlesnakeMoveResponse{
		Move: direction.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}
