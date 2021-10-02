package predict

import (
	"math"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

// moveSnakeBody shifts the snake's body one move forward
// if the snake ate last turn (health restored to full), grow it
// from the tail. Return tail coordinate for projecting next turn
func (p *State) moveSnakeBody(snake *game.Battlesnake, idx int) *game.Coord {
	// If ate, grow tail
	// Else, remove tail guess
	ate := snake.Health == 100
	tail := snake.Body[len(snake.Body)-1]
	if ate {
		snake.Body = append(snake.Body, tail)
	} else {
		p.BodyGuesses[idx].Clear(&tail)
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
func (p *State) eatSnakeFood(snake *game.Battlesnake, idx int, coord, tail *game.Coord, moveProb float64) float64 {
	foodProb := p.FoodGuesses.Prob(coord)
	if foodProb == guess.Impossible {
		return guess.Impossible
	}

	for i, food := range p.Board.Food {
		if food.Hit(coord) {
			p.Board.Food = append(p.Board.Food[:i], p.Board.Food[i+1:]...)
			break
		}
	}

	eatProb := moveProb * foodProb

	// Remaining food probability is current food probability
	// times probability we *didn't* move onto it
	p.FoodGuesses.Mult(coord, (1 - moveProb))

	// Restore health
	healthMissing := 100 - snake.Health
	healthRestored := float64(healthMissing) * eatProb
	snake.Health += int32(math.Round(healthRestored))

	// Grow tail
	p.BodyGuesses[idx].Add(tail, eatProb)

	return eatProb
}
