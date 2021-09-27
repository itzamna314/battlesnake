package predict

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightFood(coord *game.Coord, me *game.Battlesnake) int32 {
	baseWeight := FoodNotHungry

	if s.wantFood(me) {
		baseWeight = FoodHungry
	}

	if s.needFood(me) {
		baseWeight = FoodStarving
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight int32

NextFood:
	for _, food := range s.FoodGuesses {
		var (
			coordDist = coord.Dist(&food.Coord)
		)

		// Check to see if this food is contested
		// If any snake is closer than us, ignore this food
		// If a longer or equal snake is the same distance, ignore
		for _, snake := range s.Board.Snakes {
			if snake.ID == me.ID {
				continue
			}

			eDist := snake.Head.Dist(&food.Coord)
			if eDist < coordDist {
				continue NextFood
			}

			if eDist == coordDist && snake.Length >= me.Length {
				continue NextFood
			}
		}

		finalWeight += int32(float64(baseWeight) * food.Probability)
	}

	return finalWeight
}

func (s *State) wantFood(me *game.Battlesnake) bool {
	if me.Health < 50 {
		return true
	}

	acceptableLength := me.Length - 2

	for _, snake := range s.Board.Snakes {
		if snake.ID == me.ID {
			continue
		}

		if snake.Length >= acceptableLength {
			return true
		}
	}

	return false
}

func (s *State) needFood(me *game.Battlesnake) bool {
	return me.Health < 20
}
