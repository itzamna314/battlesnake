package predict

import (
	"sort"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
	"github.com/itzamna314/battlesnake/tree"
)

type State struct {
	Board game.Board
	You   game.Battlesnake

	HeadGuesses SnakeVision
	BodyGuesses SnakeVision
	FoodGuesses guess.CoordSet
}

// Init sets up our predict State for tree traversal
// This assumes we are predicting on behalf of You
func (s *State) Init(gs *game.GameState) {
	s.Board = gs.Board.Clone()
	s.You = gs.You.Clone()

	// Sort snakes by length
	// This allows us to project short snakes avoiding long ones
	sort.Slice(s.Board.Snakes, func(i, j int) bool {
		return s.Board.Snakes[i].Length > s.Board.Snakes[j].Length
	})

	// Initialize body guesses
	s.BodyGuesses = make(SnakeVision, len(s.Board.Snakes))

	for i := range s.Board.Snakes {
		for _, body := range s.Board.Snakes[i].Body {
			s.BodyGuesses[i].Set(&body, guess.Certain)
		}
	}

	// Initialize head guesses
	s.HeadGuesses = make(SnakeVision, len(s.Board.Snakes))

	for i := range s.Board.Snakes {
		s.HeadGuesses[i].Set(&s.Board.Snakes[i].Head, guess.Certain)
	}

	// Initialize food
NextFood:
	for i := 0; i < len(s.Board.Food); i++ {
		food := s.Board.Food[i]

		for _, snake := range s.Board.Snakes {
			if snake.ID == s.You.ID {
				continue
			}

			eDist := snake.Head.Dist(&food)
			youDist := s.You.Head.Dist(&food)

			if eDist < youDist {
				continue NextFood
			}

			if eDist == youDist && snake.Length >= s.You.Length {
				continue NextFood
			}
		}

		s.FoodGuesses.Set(&food, guess.Certain)
	}
}

func (s *State) Clone() tree.SnakeBrain {
	clone := State{
		Board:       s.Board.Clone(),
		You:         s.You.Clone(),
		HeadGuesses: s.HeadGuesses.Clone(),
		BodyGuesses: s.BodyGuesses.Clone(),
		FoodGuesses: s.FoodGuesses.Clone(),
	}

	return &clone
}

const (
	ateFoodCutoff = 0.1
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

	// Eat any food
	var ate bool
	step := snake.Head.Step(dir)

	if s.FoodGuesses.Prob(&step) > ateFoodCutoff {
		ate = true
		s.FoodGuesses.Clear(&step)

		for i, food := range s.Board.Food {
			if food.Hit(&step) {
				s.Board.Food = append(s.Board.Food[:i], s.Board.Food[i+1:]...)
				break
			}
		}
	}

	if ate {
		snake.Health = 100
	} else {
		snake.Health -= 1

		for _, hazard := range s.Board.Hazards {
			if hazard.Hit(&step) {
				snake.Health -= 15
				break
			}
		}
	}

	// Move body
	s.moveSnakeBody(snake, ate)

	if !ate {
		s.BodyGuesses[myIdx].Clear(&snake.Body[len(snake.Body)-1])
	}

	// Step the head in direction, and copy to body
	snake.Head = step
	snake.Body[0] = snake.Head

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

	// Clear guess for tail if snake didn't eat
	ate := enemy.Health == 100
	if !ate {
		tail := enemy.Body[len(enemy.Body)-1]
		s.BodyGuesses[idx].Clear(&tail)
	}

	// Move enemy snake probabilistically
	s.moveSnakeBody(&s.Board.Snakes[idx], ate)

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

			eatProb := s.FoodGuesses.Mult(opt, headProb)
			if eatProb > 0 {
				wouldRestore := float64(100 - enemy.Health)
				s.Board.Snakes[idx].Health += int32(wouldRestore * eatProb)
			}

			for _, hazard := range s.Board.Hazards {
				if hazard.Hit(opt) {
					pNotEat := 1 - eatProb
					s.Board.Snakes[idx].Health -= int32(15 * headProb * pNotEat)
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
