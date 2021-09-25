package predict

import (
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/guess"
)

type State struct {
	Board game.Board
	You   game.Battlesnake

	HeadGuesses SnakeVision
	BodyGuesses SnakeVision
	FoodGuesses guess.CoordSet
}

func Initialize(gs *game.GameState) *State {
	ps := State{
		Board: gs.Board.Clone(),
		You:   gs.You.Clone(),
	}

	// Initialize body guesses
	ps.BodyGuesses = make(SnakeVision, len(gs.Board.Snakes))

	for i, snake := range gs.Board.Snakes {
		if snake.ID == gs.You.ID {
			continue
		}

		for _, body := range gs.Board.Snakes[i].Body {
			ps.BodyGuesses[i].Set(&body, guess.Certain)
		}
	}

	// Initialize head guesses
	ps.HeadGuesses = make(SnakeVision, len(gs.Board.Snakes))

	for i, snake := range gs.Board.Snakes {
		if snake.ID == gs.You.ID {
			continue
		}

		ps.HeadGuesses[i].Set(&gs.Board.Snakes[i].Head, guess.Certain)
	}

	// Initialize food
	for _, food := range gs.Board.Food {
		ps.FoodGuesses.Set(&food, guess.Certain)
	}

	return &ps
}

func (p *State) Clone() *State {
	clone := State{
		Board:       p.Board.Clone(),
		You:         p.You.Clone(),
		HeadGuesses: p.HeadGuesses.Clone(),
		BodyGuesses: p.BodyGuesses.Clone(),
		FoodGuesses: p.FoodGuesses.Clone(),
	}

	return &clone
}

const (
	ateFoodCutoff = 0.1
)

func (p *State) Move(dir game.Direction) {
	// Eat any food
	var ate bool
	step := p.You.Head.Step(dir)

	if p.FoodGuesses.Prob(&step) > ateFoodCutoff {
		ate = true
		p.FoodGuesses.Clear(&step)

		for i, food := range p.Board.Food {
			if food.Hit(&step) {
				p.Board.Food = append(p.Board.Food[:i], p.Board.Food[i+1:]...)
				break
			}
		}
	}

	// Move body
	p.moveSnakeBody(&p.You, ate)

	// Step the head in direction, and copy to body
	p.You.Head = p.You.Head.Step(dir)
	p.You.Body[0] = p.You.Head
}

func (p *State) MoveEnemies() {
	for i, snake := range p.Board.Snakes {
		if snake.ID == p.You.ID {
			continue
		}

		p.moveEnemy(i)
	}
}

func (p *State) moveEnemy(idx int) {
	enemy := p.Board.Snakes[idx]
	if len(enemy.Body) == 0 {
		return
	}

	// Move enemies probabilistically based on last certain segment
	// If no certain segments remain, do not continue validation
	// Filter out certain death options

	for _, headGuess := range p.HeadGuesses[idx] {
		opts := Options(&headGuess.Coord)
		var legalMoves int

		for i, opt := range opts {
			if SnakeWillDie(p, &opt.Coord, &enemy) {
				opts[i] = nil
				continue
			}

			if p.BodyGuesses[idx].Prob(&opt.Coord) > 0 {
				opts[i] = nil
				continue
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

			p.HeadGuesses[idx].Add(&opt.Coord, headProb)
			p.FoodGuesses.Mult(&opt.Coord, headProb)
		}

		p.HeadGuesses[idx].Clear(&headGuess.Coord)
		p.BodyGuesses[idx].Set(&headGuess.Coord, headGuess.Probability)
	}

	// Clear guess for tail if snake didn't eat
	ate := enemy.Health == 100
	if !ate {
		tail := enemy.Body[len(enemy.Body)-1]
		p.BodyGuesses[idx].Clear(&tail)
	}

	// Move enemy snake probabilistically
	p.moveSnakeBody(&p.Board.Snakes[idx], ate)

	// We don't know where the head went
	// Remove from deterministic structure
	p.Board.Snakes[idx].Body = enemy.Body[1:]
}
