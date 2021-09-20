package move

import "github.com/itzamna314/battlesnake/model"

func food(state *model.GameState, possible model.PossibleMoves) {
	// Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	var (
		myHead       = state.You.Body[0]
		closestFoods []model.Coord
		minDist      int
	)
	for _, food := range state.Board.Food {
		dist := myHead.Dist(&food)

		if minDist == 0 || dist < minDist {
			minDist = dist
			closestFoods = []model.Coord{food}
		} else if dist == minDist {
			closestFoods = append(closestFoods, food)
		}
	}

	if minDist == 0 {
		return
	}

	// Prefer to move toward the nearest food
	for _, food := range closestFoods {
		step := myHead.StepToward(&food)
		possible[step].Weight += 0.15

		if minDist == 1 {
			possible[step].Shout = "OMNOMNOM"
		} else {
			possible[step].Shout = "hungry"
		}
	}
}
