package predict

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

// Calculate weight for moving You to coord in state
func (s *State) weightEnemies(coord *game.Coord, snake *Snake) float64 {
	// No enemies
	if len(s.Board.Snakes) <= 1 {
		return Nothing
	}

	var weight float64
	for i, enemy := range s.Board.Snakes {
		if enemy.ID == snake.ID {
			continue
		}

		prob := s.BodyGuesses[i].Prob(coord)
		if prob == guess.Certain {
			return Death
		}
		weight += (prob * EnemyAvoid)

		// If we are shorter, avoid with weight of probability
		// Otherwise, attack with reduced weight of collision probability
		// STRIKE FIRST STRIKE HARD NO MERCY
		// But also don't chase a short snake into a long snake
		prob = s.HeadGuesses[i].Prob(coord)
		if enemy.Length == snake.Length {
			weight += (prob * EnemyTie)
		} else if enemy.Length > snake.Length {
			weight += (prob * EnemyAvoid)
		} else {
			weight += (prob * EnemyKill)
		}
	}

	return weight
}
