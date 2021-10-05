package multiverse

import "github.com/itzamna314/battlesnake/guess"

// Calculate weight for moving You to coord in state
func (s *State) weightEnemies(snake *Snake) FloatWeight {
	// No enemies
	if len(s.Board.Snakes) <= 1 {
		return Nothing
	}

	var weight FloatWeight
	for i, enemy := range s.Board.Snakes {
		if enemy.ID == snake.ID {
			continue
		}

		prob := s.BodyGuesses[i].Prob(&snake.Head)
		if prob == guess.Certain {
			return Death
		}
		weight += FloatWeight(prob) * EnemyAvoid

		// If we are shorter, avoid with weight of probability
		// Otherwise, attack with reduced weight of collision probability
		// STRIKE FIRST STRIKE HARD NO MERCY
		// But also don't chase a short snake into a long snake
		prob = s.HeadGuesses[i].Prob(&snake.Head)
		if enemy.Length == snake.Length {
			weight += FloatWeight(prob) * EnemyTie
		} else if enemy.Length > snake.Length {
			weight += FloatWeight(prob) * EnemyAvoid
		} else {
			weight += FloatWeight(prob) * EnemyKill
		}
	}

	return weight
}
