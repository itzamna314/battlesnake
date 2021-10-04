package multiverse

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

func (s *State) Move(snakeID string, dir game.Direction) {
	// Get my index
	var (
		myIdx int
		snake *Snake
	)

	for i, snk := range s.Board.Snakes {
		if snk.ID == snakeID {
			myIdx = i
			snake = &s.Board.Snakes[i]
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
			snake.Health -= hazardDmg
			break
		}
	}

	// Step the head in direction, and copy to body
	snake.Head = step
	snake.Body[0] = snake.Head

	s.BodyGuesses[myIdx].Set(&snake.Body[1], guess.Certain)
	s.HeadGuesses[myIdx].Clear(&snake.Body[1])
	s.HeadGuesses[myIdx].Set(&snake.Head, guess.Certain)

	// Destroy any enemy head guesses that we would have eaten
	for i, enemy := range s.Board.Snakes {
		if enemy.ID == snake.ID {
			continue
		}

		if snake.Length <= enemy.Length {
			continue
		}

		s.HeadGuesses[i].Clear(&snake.Head)
	}

	// If we may have eaten, re-append our tail
	if ateProb > 0 {
		snake.Body = append(snake.Body, *tail)
	}

	if snake.ID == s.You.ID {
		s.You = *snake
	}
}

func (s *State) MoveEnemies(snakeID string) {
	for i, snake := range s.Board.Snakes {
		if snake.ID == snakeID {
			continue
		}

		s.moveEnemy(i)
	}
}

func (s *State) moveEnemy(idx int) {
	enemy := &s.Board.Snakes[idx]
	if len(s.Board.Snakes[idx].Body) == 0 {
		return
	}

	// Move enemy snake
	tail := s.moveSnakeBody(enemy, idx)

	// Move enemies probabilistically based on last certain segment
	// If no certain segments remain, do not continue projection
	// Filter out certain death options
	// For each guess, write the next level into nextGuesses, and replace
	var (
		nextGuesses guess.CoordSet
		maybeAte    bool
	)

	for _, headGuess := range s.HeadGuesses[idx] {
		opts := game.Options(&headGuess.Coord)
		var legalMoves int

	NextOpt:
		for i, opt := range opts {
			if SnakeWillDie(s, opt, enemy) {
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

				if len(eEnemy.Body) <= 1 {
					continue
				}

				eeBodyProb := s.BodyGuesses[eeIdx].Prob(opt)
				if eeBodyProb > 0.333 {
					tail := eEnemy.Body[len(eEnemy.Body)-1]

					// If this is the tail skip it...
					if tail.Hit(opt) {
						tNeck := eEnemy.Body[len(eEnemy.Body)-2]

						// ... unless they just ate
						if !tail.Hit(&tNeck) {
							continue
						}
					}

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

			eatProb := s.eatSnakeFood(enemy, idx, opt, tail, headProb)
			if eatProb > guess.Impossible {
				maybeAte = true
			}

			for _, hazard := range s.Board.Hazards {
				if hazard.Hit(opt) {
					pNotEat := 1 - eatProb
					enemy.Health -= 15 * headProb * pNotEat
					break
				}
			}
		}

		s.BodyGuesses[idx].Set(&headGuess.Coord, headGuess.Probability)
	}

	s.HeadGuesses[idx] = nextGuesses
	enemy.Health -= 1

	// We don't know where the head went
	enemy.Body = enemy.Body[1:]

	// If we maybe ate, re-attach our tail so we can clean up guesses
	if maybeAte {
		enemy.Body = append(enemy.Body, *tail)
	}
}
