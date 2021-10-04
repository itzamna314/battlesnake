package multiverse

import (
	"math"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

type Snake struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Health float64      `json:"health"`
	Body   []game.Coord `json:"body"`
	Head   game.Coord   `json:"head"`
	Length float64      `json:"length"`
}

func (s *Snake) Clone() Snake {
	clone := Snake{
		ID:     s.ID,
		Name:   s.Name,
		Health: s.Health,
		Head:   s.Head,
		Length: s.Length,
		Body:   make([]game.Coord, len(s.Body)),
	}

	for i, bd := range s.Body {
		clone.Body[i] = bd
	}

	return clone
}

func (s *Snake) Init(gs *game.Battlesnake) {
	s.ID = gs.ID
	s.Name = gs.Name
	s.Health = float64(gs.Health)
	s.Head = gs.Head
	s.Length = float64(gs.Length)

	s.Body = make([]game.Coord, len(gs.Body))
	for i := range gs.Body {
		s.Body[i] = gs.Body[i]
	}
}

// moveSnakeBody shifts the snake's body one move forward
// if the snake ate last turn (health restored to full), grow it
// from the tail. Return tail coordinate for projecting next turn
func (p *State) moveSnakeBody(snake *Snake, idx int) *game.Coord {
	// If ate, grow tail
	// Else, remove tail guess
	ate := snake.Health == 100
	tail := snake.Body[len(snake.Body)-1]
	if ate {
		snake.Body = append(snake.Body, tail)
	} else {
		tailProb := p.BodyGuesses[idx].Clear(&tail)
		if len(snake.Body) > 1 {
			newTail := snake.Body[len(snake.Body)-2]
			p.BodyGuesses[idx].Set(&newTail, tailProb)
		}
	}

	// Copy each body segment to next
	// Head will remain copied into neck
	for i := len(snake.Body) - 1; i > 0; i-- {
		next := i - 1
		snake.Body[i] = snake.Body[next]
	}

	return &tail
}

// eatSnakeFood simulates the snake eating food at coord, with probability moveProb
// returns the probability that the snake ate
func (p *State) eatSnakeFood(snake *Snake, idx int, coord, tail *game.Coord, moveProb float64) float64 {
	foodProb := p.FoodGuesses.Prob(coord)
	if foodProb == guess.Impossible {
		return guess.Impossible
	}

	eatProb := moveProb * foodProb

	for i, food := range p.Board.Food {
		if food.Hit(coord) {
			p.Board.Food = append(p.Board.Food[:i], p.Board.Food[i+1:]...)
			break
		}
	}

	// Remaining food probability is current food probability
	// times probability we *didn't* move onto it
	p.FoodGuesses.Mult(coord, (1 - moveProb))

	// Restore health
	healthMissing := 100 - snake.Health
	healthRestored := healthMissing * eatProb
	snake.Health += math.Round(healthRestored)

	// Grow tail
	p.BodyGuesses[idx].Add(tail, eatProb)
	snake.Length += eatProb

	return eatProb
}
