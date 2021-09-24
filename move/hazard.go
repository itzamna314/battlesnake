package move

import (
	"github.com/itzamna314/battlesnake/game"
)

func WeightHazard(state *game.GameState, coord *game.Coord) float64 {
	for _, hazard := range state.Board.Hazards {
		if hazard.Hit(coord) {
			return Hazard
		}
	}

	return Nothing
}
