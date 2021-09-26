package move

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
)

func WeightHazard(state *predict.State, coord *game.Coord) float64 {
	for _, hazard := range state.Board.Hazards {
		if hazard.Hit(coord) {
			// hazard on top of food won't hurt
			if state.FoodGuesses.Prob(&hazard) > 0.75 {
				continue
			}

			if state.You.Health <= 0 {
				return Death
			}

			return Nothing
		}
	}

	return Nothing
}
