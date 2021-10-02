package predict

import (
	"math"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

func (s *State) Move(snake *game.Battlesnake, dir game.Direction) {
	// Get my index
	var myIdx int
	for i, snk := range s.Board.Snakes {
		if snk.ID == snake.ID {
			myIdx = i
			break
		}
	}

	step := snake.Head.Step(dir)

	// Move body
	tail := s.moveSnakeBody(snake, myIdx)

	// Consume food
	ateProb := s.eatSnakeFood(snake, myIdx, &step, tail, guess.Certain)

	// Take standard damage
	snake.Health -= 1

	// Take hazard damage
	for _, hazard := range s.Board.Hazards {
		if hazard.Hit(&step) {
			hazardDmg := float64(15) * (1 - ateProb)
			snake.Health -= int32(math.Round(hazardDmg))
			break
		}
	}

	// Step the head in direction, and copy to body
	snake.Head = step
	snake.Body[0] = snake.Head

	s.BodyGuesses[myIdx].Set(&snake.Body[1], guess.Certain)
	s.HeadGuesses[myIdx].Clear(&snake.Body[1])
	s.HeadGuesses[myIdx].Set(&snake.Head, guess.Certain)

	if snake.ID == s.You.ID {
		s.You = *snake
	}
}

func (s *State) MoveEnemies(me *game.Battlesnake) {
	for i, snake := range s.Board.Snakes {
		if snake.ID == me.ID {
			continue
		}

		s.moveEnemy(i)
	}
}

func (s *State) moveEnemy(idx int) {
	enemy := s.Board.Snakes[idx]
	if len(enemy.Body) == 0 {
		return
	}

	// Move enemy snake
	tail := s.moveSnakeBody(&s.Board.Snakes[idx], idx)

	// Move enemies probabilistically based on last certain segment
	// If no certain segments remain, do not continue projection
	// Filter out certain death options
	// For each guess, write the next level into nextGuesses, and replace
	var nextGuesses guess.CoordSet

	for _, headGuess := range s.HeadGuesses[idx] {
		opts := game.Options(&headGuess.Coord)
		var legalMoves int

	NextOpt:
		for i, opt := range opts {
			if SnakeWillDie(s, opt, &enemy) {
				opts[i] = nil
				continue
			}

			if s.BodyGuesses[idx].Prob(opt) >= 0.25 {
				opts[i] = nil
				continue
			}

			for eeIdx, eEnemy := range s.Board.Snakes {
				if eEnemy.ID == enemy.ID {
					continue
				}

				eeHeadProb := s.HeadGuesses[eeIdx].Prob(opt)
				if eeHeadProb > 0.333 {
					if eEnemy.Length >= enemy.Length {
						opts[i] = nil
						continue NextOpt
					} else {
						s.HeadGuesses[eeIdx].Clear(opt)
					}
				}

				eeBodyProb := s.BodyGuesses[eeIdx].Prob(opt)
				if eeBodyProb > 0.333 {
					opts[i] = nil
					continue NextOpt
				}
			}

			legalMoves++
		}

		// Distribute move evenly among non-death options
		// Check for eating at each one
		headProb := 1.0 / float64(legalMoves) * headGuess.Probability
		for _, opt := range opts {
			if opt == nil {
				continue
			}

			nextGuesses.Add(opt, headProb)

			eatProb := s.eatSnakeFood(&s.Board.Snakes[idx], idx, opt, tail, headProb)

			for _, hazard := range s.Board.Hazards {
				if hazard.Hit(opt) {
					pNotEat := 1 - eatProb
					s.Board.Snakes[idx].Health -= int32(15 * headProb * pNotEat)
					break
				}
			}
		}

		s.BodyGuesses[idx].Set(&headGuess.Coord, headGuess.Probability)
	}

	s.HeadGuesses[idx] = nextGuesses
	s.Board.Snakes[idx].Health -= 1

	// We don't know where the head went
	// Remove from deterministic structure
	s.Board.Snakes[idx].Body = enemy.Body[1:]
}
