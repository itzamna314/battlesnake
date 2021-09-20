package move

import "github.com/itzamna314/battlesnake/model"

func food(state *model.GameState, possible model.PossibleMoves) {
	// Stop eating if we the longest snake
	eat := false
	for _, enemy := range state.Board.Snakes {
		if enemy.ID == state.You.ID {
			continue
		}

		if enemy.Length >= state.You.Length {
			eat = true
		}
	}

	if !eat {
		return
	}

	// Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	var (
		myHead      = state.You.Body[0]
		closestFood model.Coord
		minDist     int
	)
	for _, food := range state.Board.Food {
		dist := myHead.Dist(&food)

		if minDist == 0 || dist < minDist {
			minDist = dist
			closestFood = food
		}
	}

	// Prefer to move toward the nearest food
	if minDist > 0 {
		step := myHead.StepToward(&closestFood)
		possible[step].Weight = 0.75

		if minDist == 1 {
			possible[step].Shout = "OMNOMNOM"
		} else {
			possible[step].Shout = "hungry"
		}
	}
}
