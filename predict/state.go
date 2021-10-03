package predict

import (
	"sort"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
	"github.com/itzamna314/battlesnake/tree"
)

type State struct {
	Board Board
	You   Snake

	HeadGuesses SnakeVision
	BodyGuesses SnakeVision
	FoodGuesses guess.CoordSet
}

// Init sets up our predict State for tree traversal
// This assumes we are predicting on behalf of You
func (s *State) Init(gs *game.GameState) {
	s.Board.Init(&gs.Board)
	s.You.Init(&gs.You)

	// Sort snakes by length
	// This allows us to project short snakes avoiding long ones
	sort.Slice(s.Board.Snakes, func(i, j int) bool {
		return s.Board.Snakes[i].Length > s.Board.Snakes[j].Length
	})

	// Initialize body guesses
	s.BodyGuesses = make(SnakeVision, len(s.Board.Snakes))

	for i, snake := range s.Board.Snakes {
		for _, body := range snake.Body {
			s.BodyGuesses[i].Set(&body, guess.Certain)
		}
	}

	// Initialize head guesses
	s.HeadGuesses = make(SnakeVision, len(s.Board.Snakes))

	for i := range s.Board.Snakes {
		s.HeadGuesses[i].Set(&s.Board.Snakes[i].Head, guess.Certain)
	}

	// Initialize food
	for i := 0; i < len(s.Board.Food); i++ {
		food := s.Board.Food[i]
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

func (s *State) Snake(snakeID string) *Snake {
	for i, snake := range s.Board.Snakes {
		if snake.ID == snakeID {
			return &s.Board.Snakes[i]
		}
	}
	return nil
}

func (s *State) Enemies(snakeID string) []*Snake {
	enemies := make([]*Snake, 0, 3)

	for i, snake := range s.Board.Snakes {
		if snake.ID == snakeID {
			continue
		}

		enemies = append(enemies, &s.Board.Snakes[i])
	}

	return enemies
}
