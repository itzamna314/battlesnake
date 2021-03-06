package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzamna314/battlesnake/game"
)

func End(w http.ResponseWriter, r *http.Request) {
	state := game.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}

	end(state)

	// Nothing to respond with here
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state game.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}
