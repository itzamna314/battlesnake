package predict

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightHazard(coord *game.Coord, me *game.Battlesnake) float64 {
	for _, hazard := range s.Board.Hazards {
		if hazard.Hit(coord) {
			// hazard on top of food won't hurt
			if s.FoodGuesses.Prob(&hazard) > 0.95 {
				continue
			}

			if me.Health <= 20 {
				return Hazard * 2
			}

			return Hazard
		}
	}

	return Nothing
}
