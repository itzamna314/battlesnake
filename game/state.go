package game

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`

	HeadGuesses SnakeVision
	BodyGuesses SnakeVision
	FoodGuesses GuessCoordSet
}

func (g *GameState) initGuesses() {
	// Initialize enemy guesses if necessary
	if len(g.BodyGuesses) == 0 {
		g.BodyGuesses = make(SnakeVision, len(g.Board.Snakes))

		for i, snake := range g.Board.Snakes {
			if snake.ID == g.You.ID {
				continue
			}

			for _, body := range g.Board.Snakes[i].Body {
				g.BodyGuesses[i].Set(&body, Certain)
			}
		}
	}

	if len(g.HeadGuesses) == 0 {
		g.HeadGuesses = make(SnakeVision, len(g.Board.Snakes))

		for i, snake := range g.Board.Snakes {
			if snake.ID == g.You.ID {
				continue
			}

			g.HeadGuesses[i].Set(&g.Board.Snakes[i].Head, Certain)
		}
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
		Game:        g.Game,
		Turn:        g.Turn,
		Board:       g.Board.Clone(),
		You:         g.You.Clone(),
		HeadGuesses: g.HeadGuesses.Clone(),
		BodyGuesses: g.BodyGuesses.Clone(),
		FoodGuesses: g.FoodGuesses.Clone(),
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
}

func (g *GameState) MoveEnemies() {
	for i, snake := range g.Board.Snakes {
		if snake.ID == g.You.ID {
			continue
		}

		g.moveEnemy(i)
	}
}

func (g *GameState) moveEnemy(idx int) {
	enemy := g.Board.Snakes[idx]
	if len(enemy.Body) == 0 {
		return
	}

	// Move enemies probabilistically based on last certain segment
	// If no certain segments remain, do not continue validation
	// Filter out certain death options

	for _, headGuess := range g.HeadGuesses[idx] {
		opts := Options(&headGuess.Coord)
		var legalMoves int

		for i, opt := range opts {
			if SnakeWillDie(g, &opt.Coord, &enemy) {
				opts[i] = nil
				continue
			}

			if g.BodyGuesses[idx].Prob(&opt.Coord) > 0 {
				opts[i] = nil
				continue
			}

			legalMoves++
		}

		// Distribute move evenly among non-death options
		// Check for eating at each one
		headProb := 1.0 / float64(legalMoves)
		for _, opt := range opts {
			if opt == nil {
				continue
			}

			g.HeadGuesses[idx].Add(&opt.Coord, headProb)
			g.FoodGuesses.Mult(&opt.Coord, headProb)
		}

		g.HeadGuesses[idx].Clear(&headGuess.Coord)
		g.BodyGuesses[idx].Set(&headGuess.Coord, headGuess.Probability)
	}

	// Clear guess for tail if snake didn't eat
	ate := enemy.Health == 100
	if !ate {
		tail := enemy.Body[len(enemy.Body)-1]
		g.BodyGuesses[idx].Clear(&tail)
	}

	// Move enemy snake probabilistically
	g.Board.Snakes[idx].MoveProb(ate)
}
