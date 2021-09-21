package move

import "github.com/itzamna314/battlesnake/model"

const FoodWeight = 0.15

func weightFood(state *model.GameState, coord *model.Coord) float64 {
	foodWeight := FoodWeight

	if !wantFood(state) {
		weight *= -1
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
		return 0
	}

	// Prefer to move toward the closest foods
	// Divide by number of closest foods
	// If there are many tied at the same distance from head,
	// lower priority and only count the ones in the right direction
	var weight float64
	for _, food := range closestFoods {
		myDist := coord.Dist(&food)
		if myDist < minDist {
			amtCloser := float64(minDist-myDist) / float64(minDist)
			weight = weight + foodWeight*amtCloser
		}
	}

	return weight / float64(len(closestFoods))
}

func wantFood(state *model.GameState) bool {
	if state.You.Health < 20 {
		return true
	}

	var myLen = state.You.Length

	for _, snake := range state.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		if snake.Length >= state.You.Length {
			return true
		}
	}

	return false
}
