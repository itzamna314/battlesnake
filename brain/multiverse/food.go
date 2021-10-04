package multiverse

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightFood(coord *game.Coord, me *Snake) float64 {
	baseWeight := FoodAvoid

	if s.wantFood(me) {
		baseWeight = Food
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight, numWeights float64

	for _, food := range s.FoodGuesses {
		var (
			headDist = me.Head.Dist(&food.Coord)
			myDist   = coord.Dist(&food.Coord)
		)

		// We didn't get closer or farther. Ignore
		if myDist == headDist {
			continue
		}

		// If we don't want food, only avoid eating
		// Don't penalize getting closer to food.
		// Not eating isn't very important
		if baseWeight < 0 && myDist > 0 {
			continue
		}

		contestFactor := s.foodContestFactor(&food.Coord, me)
		if contestFactor == 0 {
			continue
		}

		distDiffPct := float64(headDist-myDist) / float64(headDist)

		finalWeight += (baseWeight * distDiffPct * food.Probability * contestFactor)
		numWeights++
	}

	if numWeights == 0 {
		return 0
	}

	return finalWeight / numWeights
}

func (s *State) wantFood(me *Snake) bool {
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

func (s *State) foodContestFactor(food *game.Coord, me *Snake) float64 {
	contestFactor := 1.0
	for _, enemy := range s.Board.Snakes {
		if enemy.ID == me.ID {
			continue
		}

		eDist := enemy.Head.Dist(food)
		myDist := me.Head.Dist(food)

		// Treat enemy as 1 closer if they can bully us off
		if enemy.Length > me.Length {
			eDist--
		}

		distDiff := myDist - eDist

		// If we are closer, no contest factor
		if distDiff < 0 {
			continue
		}

		// Subtract 0.25 for each dist diff
		// If factor reaches 0, return 0
		contestFactor += -0.25 * float64(distDiff)
		if contestFactor <= 0 {
			return 0
		}
	}

	return contestFactor
}
