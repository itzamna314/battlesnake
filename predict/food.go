package predict

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightFood(coord *game.Coord, me *game.Battlesnake) float64 {
	baseWeight := FoodAvoid

	if s.wantFood(me) {
		baseWeight = Food
	}

	if s.needFood(me) {
		baseWeight = Mandatory
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight, numWeights float64

NextFood:
	for _, food := range s.FoodGuesses {
		var (
			headDist = me.Head.Dist(&food.Coord)
			myDist   = coord.Dist(&food.Coord)
		)

		// We didn't get closer. Ignore
		if myDist == headDist {
			continue
		}

		// Check to see if this food is contested
		// If any snake is closer than us, ignore this food
		// If a longer or equal snake is the same distance, ignore
		for _, snake := range s.Board.Snakes {
			if snake.ID == me.ID {
				continue
			}

			eDist := snake.Head.Dist(&food.Coord)
			if eDist < myDist {
				continue NextFood
			}

			if eDist == myDist && snake.Length >= s.You.Length {
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
