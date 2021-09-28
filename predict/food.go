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

	// Prefer to move toward foods
	for _, food := range s.FoodGuesses {
		if !coord.Hit(&food.Coord) {
			continue
		}

		return baseWeight
	}

	return Neutral
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
