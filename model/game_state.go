package model

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`

	Future BoardVision
}

func (g *GameState) Clone() GameState {
	var future BoardVision
	if g.Future != nil {
		future = g.Future.Clone()
	} else {
		future = g.initFutureToPresent()
	}

	return GameState{
		Game:   g.Game,
		Turn:   g.Turn,
		Board:  g.Board.Clone(),
		You:    g.You.Clone(),
		Future: future,
	}
}

func (g *GameState) MoveSnake(snake Battlesnake, dir Direction) {
	if !snake.Head.Hit(&snake.Body[0]) {
		panic("illegal snake")
	}

	// Copy each body segment to next
	// Head will remain copied into neck
	for i := len(snake.Body) - 1; i > 0; i-- {
		prev := i - 1
		snake.Body[i] = snake.Body[prev]
	}

	snake.Head = snake.Head.Step(dir)
	snake.Body[0] = snake.Head

	for i, s := range g.Board.Snakes {
		if s.ID == snake.ID {
			g.Board.Snakes[i] = snake
			break
		}
	}

	if snake.ID == g.You.ID {
		g.You = snake
	}

	// Eat any food
	for i, food := range g.Board.Food {
		if food.Hit(&snake.Head) {
			g.Board.Food = append(g.Board.Food[:i], g.Board.Food[i+1:]...)
			break
		}
	}
}

func (g *GameState) initFutureToPresent() BoardVision {
	future := make(BoardVision, g.Board.Width)

	for x := 0; x < g.Board.Width; x++ {
		future[x] = make([]VisionCell, g.Board.Height)
	}

	for i, snake := range g.Board.Snakes {
		if snake.ID == g.You.ID {
			continue
		}

		for _, b := range snake.Body {
			future[b.X][b.Y].Enemies[i] = Certain
		}
	}

	for _, f := range g.Board.Food {
		future[f.X][f.Y].Food = Certain
	}

	return future
}
