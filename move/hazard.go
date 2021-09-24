package move

import (
	"github.com/itzamna314/battlesnake/model"
)

func WeightHazard(state *model.GameState, coord *model.Coord) float64 {
	for _, hazard := range state.Board.Hazards {
		if hazard.Hit(coord) {
			return Hazard
		}
	}

	return Nothing
}
