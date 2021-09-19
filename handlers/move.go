package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzamna314/battlesnake/handlers/move"
	"github.com/itzamna314/battlesnake/model"
)

func Move(w http.ResponseWriter, r *http.Request) {
	state := model.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	response := move.Next(state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}
