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

	for i, snake := range s.Board.Snakes {
		for _, body := range snake.Body {
			// Only set tail to certain if this snake ate (full health)
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
