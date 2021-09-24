package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzamna314/battlesnake/game"
)

func Start(w http.ResponseWriter, r *http.Request) {
	state := game.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(state)

	// Nothing to respond with here
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state game.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}
