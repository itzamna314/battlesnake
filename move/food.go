package move

import "github.com/itzamna314/battlesnake/model"

func weightFood(state *model.GameState, coord *model.Coord) float64 {
	baseWeight := Food

	if !wantFood(state) {
		baseWeight *= -1
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight, numWeights float64

	for _, food := range state.FoodGuesses {
		var (
			headDist = state.You.Head.Dist(&food.Coord)
			myDist   = coord.Dist(&food.Coord)
		)

		// We didn't get closer. Ignore
		if myDist == headDist {
			continue
		}

		distDiffPct := float64(headDist-myDist) / float64(headDist)

		finalWeight += baseWeight * distDiffPct * food.Probability
		numWeights++
	}

	if numWeights == 0 {
		return 0
	}

	return finalWeight / numWeights
}

func wantFood(state *model.GameState) bool {
	if state.You.Health < 20 {
		return true
	}

	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		if snake.Length >= state.You.Length {
			return true
		}
	}

	return false
}
