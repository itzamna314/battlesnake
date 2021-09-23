package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzamna314/battlesnake/model"
)

// Index is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func Index(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func info() model.BattlesnakeInfoResponse {
	log.Println("INFO")
	return model.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Johnny Lawrence",
		Color:      "#000000",
		Head:       "evil",
		Tail:       "small-rattle",
	}
}
