package move

import "github.com/itzamna314/battlesnake/model"

func weightFood(state *model.GameState, coord *model.Coord) float64 {
	baseWeight := WeightFood

	if !wantFood(state) {
		baseWeight *= -1
	}

	// Find closest food to our head
	var (
		closestFoods []model.Coord
		minDist      int
	)
	for _, food := range state.Board.Food {
		dist := state.You.Head.Dist(&food)

		if minDist == 0 || dist < minDist {
			minDist = dist
			closestFoods = []model.Coord{food}
		} else if dist == minDist {
			closestFoods = append(closestFoods, food)
		}
	}

	// No food
	if minDist == 0 {
		return WeightNothing
	}

	// Prefer to move toward or away from the closest foods
	// Divide by number of closest foods
	// If there are many tied at the same distance from head,
	// lower priority and only count the ones in the right direction
	var finalWeight, numWeights float64

	for _, food := range closestFoods {
		myDist := coord.Dist(&food)
		if myDist < minDist {
			amtCloser := float64(minDist-myDist) / float64(minDist)
			finalWeight += baseWeight * amtCloser
			numWeights++
		}
	}

	if numWeights == 0 {
		return WeightNothing
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
