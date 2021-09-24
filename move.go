package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/move"
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

	direction := move.Next(state)
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
