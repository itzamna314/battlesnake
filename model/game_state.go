package model

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`

	EnemyGuesses SnakeVision
	FoodGuesses  GuessCoordSet
}

func (g *GameState) initGuesses() {
	// Initialize enemy guesses if necessary
	if len(g.EnemyGuesses) == 0 {
		g.EnemyGuesses = make(SnakeVision, len(g.Board.Snakes))
	}

	// Initialize food if necessary
	if len(g.FoodGuesses) != len(g.Board.Food) {
		for _, food := range g.Board.Food {
			g.FoodGuesses.Set(&food, Certain)
		}
	}
}

func (g *GameState) Clone() GameState {
	clone := GameState{
		Game:         g.Game,
		Turn:         g.Turn,
		Board:        g.Board.Clone(),
		You:          g.You.Clone(),
		EnemyGuesses: g.EnemyGuesses.Clone(),
		FoodGuesses:  g.FoodGuesses.Clone(),
	}

	clone.initGuesses()
	return clone
}

const (
	ateFoodCutoff = 0.1
)

func (g *GameState) Move(dir Direction) {
	// Eat any food
	var ate bool
	step := g.You.Head.Step(dir)

	if g.FoodGuesses.Prob(&step) > ateFoodCutoff {
		ate = true
		g.FoodGuesses.Clear(&step)

		for i, food := range g.Board.Food {
			if food.Hit(&step) {
				g.Board.Food = append(g.Board.Food[:i], g.Board.Food[i+1:]...)
				break
			}
		}
	}

	// Move self (deterministic)
	g.You.MoveDet(dir, ate)

	for i, snake := range g.Board.Snakes {
		if snake.ID == g.You.ID {
			continue
		}

		g.MoveEnemy(i)
	}
}

func (g *GameState) MoveEnemy(idx int) {
	enemy := g.Board.Snakes[idx]
	if len(enemy.Body) == 0 {
		return
	}

	// Move enemies probabilistically based on last certain segment
	// If no certain segments remain, do not continue validation
	// Filter out certain death options
	opts := Options(&enemy.Body[0])
	var legalMoves int

	for i, opt := range opts {
		if SnakeWillDie(g, &opt.Coord, &enemy) {
			opts[i] = nil
			continue
		}

		legalMoves++
	}

	// Distribute move evenly among non-death options
	// Check for eating at each one
	var ate bool
	headProb := 1.0 / float64(legalMoves)
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		g.EnemyGuesses[idx].Add(&opt.Coord, headProb)

		// Arbitrary cutoff to assume they ate
		if g.FoodGuesses.Prob(&opt.Coord) > 0.1 {
			g.FoodGuesses.Mult(&opt.Coord, headProb)
			ate = true
		}
	}

	// Move enemy snake probabilistically
	g.Board.Snakes[idx].MoveProb(ate)

	// Record certainties
	for _, body := range g.Board.Snakes[idx].Body {
		g.EnemyGuesses[idx].Set(&body, Certain)
	}
}
