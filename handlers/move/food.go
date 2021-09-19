package move

import (
	"fmt"

	"github.com/itzamna314/battlesnake/model"
)

func food(state model.GameState, possible model.PossibleMoves) {
	// Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	var (
		myHead      = state.You.Body[0]
		closestFood model.Coord
		minDist     int
	)
	fmt.Printf("Closest food to my head %s\n", myHead)
	for _, food := range state.Board.Food {
		fmt.Printf("Considering food %s\n", food)
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
