package predict

import (
	"fmt"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

// Calculate weight for moving You to coord in state
func (s *State) weightEnemies(coord *game.Coord, snake *game.Battlesnake) int32 {
	// No enemies
	if len(s.Board.Snakes) <= 1 {
		return Neutral
	}

	var weight int32
	for i, enemy := range s.Board.Snakes {
		if enemy.ID == snake.ID {
			continue
		}

		if coord.Hit(&game.Coord{2, 10}) {
			fmt.Printf("Enemy %s probability body %v head %v\n", s.Board.Snakes[i].Name, s.BodyGuesses[i].Prob(coord), s.HeadGuesses[i].Prob(coord))
		}

		prob := s.BodyGuesses[i].Prob(coord)
		if prob == guess.Certain {
			return CertainDeath
		}

		// An enemy body is certain death
		// If we're not sure, multiply probability by death
		weight += int32(prob * float64(CertainDeath))

		// If we are shorter, avoid with weight of probability
		// Otherwise, attack with reduced weight of collision probability
		// STRIKE FIRST STRIKE HARD NO MERCY
		// But also don't chase a short snake into a long snake
		prob = s.HeadGuesses[i].Prob(coord)
		if enemy.Length >= snake.Length {
			weight += int32(prob * float64(CertainDeath))
		} else if len(s.Board.Snakes) == 2 {
			weight += int32(prob * float64(CertainWin))
		} else {
			weight += int32(prob * float64(EnemyKill))
		}

		if weight <= CertainDeath {
			return CertainDeath
		}
	}

	return weight
}
