package predict

import "github.com/itzamna314/battlesnake/game"

func (s *State) weightHealth(snake *game.Battlesnake) float64 {
	return Base * float64(snake.Health) / 100
}
