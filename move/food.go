package move

import "github.com/itzamna314/battlesnake/game"

func WeightFood(state *game.GameState, coord *game.Coord) float64 {
	baseWeight := Food

	if !wantFood(state) {
		baseWeight *= -0.25
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight, numWeights float64

NextFood:
	for _, food := range state.FoodGuesses {
		var (
			headDist = state.You.Head.Dist(&food.Coord)
			myDist   = coord.Dist(&food.Coord)
		)

		// We didn't get closer. Ignore
		if myDist == headDist {
			continue
		}

		// Check to see if this food is contested
		// If any snake is closer than us, ignore this food
		// If a longer or equal snake is the same distance, ignore
		for _, snake := range state.Board.Snakes {
			if snake.ID == state.You.ID {
				continue
			}

			eDist := snake.Head.Dist(&food.Coord)
			if eDist < myDist {
				continue NextFood
			}

			if eDist == myDist && snake.Length >= state.You.Length {
				continue NextFood
			}
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

func wantFood(state *game.GameState) bool {
	if state.You.Health < 50 {
		return true
	}

	acceptableLength := state.You.Length - 2

	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		if snake.Length >= acceptableLength {
			return true
		}
	}

	return false
}
