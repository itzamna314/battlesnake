package multiverse

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightFood(snake *Snake) FloatWeight {
	baseWeight := FoodAvoid

	if s.wantFood(snake) {
		baseWeight = Food
	}

	// Prefer to move toward or away from foods
	// Weight foods more strongly by the likelihood that they will remain
	// Divide by number of foods where this move changes the distance
	var finalWeight FloatWeight

	// If we just ate, apply the full weight
	if snake.Health == 100 {
		finalWeight += baseWeight
	}

	// Add a small weight adjustment based on the distance squared to food
	// The closer we are to food, the higher the adjustment
	// This must be insignificant compared to the value of just eating
	// Otherwise the snake will circle around food without eating it
	for _, food := range s.FoodGuesses {
		var headDist = snake.Head.Dist(&food.Coord)
		headDist *= headDist

		contestFactor := s.foodContestFactor(&snake.Head, snake)
		if contestFactor == 0 {
			continue
		}

		distWeight := 0.0001 * (100 - FloatWeight(headDist))

		finalWeight += (baseWeight * distWeight * FloatWeight(food.Probability) * contestFactor)
	}

	return finalWeight
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

func (s *State) foodContestFactor(food *game.Coord, me *Snake) FloatWeight {
	contestFactor := FloatWeight(1.0)
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
		contestFactor += -0.25 * FloatWeight(distDiff)
		if contestFactor <= 0 {
			return 0
		}
	}

	return contestFactor
}
