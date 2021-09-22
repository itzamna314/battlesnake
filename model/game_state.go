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

func (g *GameState) Move(dir Direction) {
	// Copy each body segment to next
	// Head will remain copied into neck
	for i := len(g.You.Body) - 1; i > 0; i-- {
		prev := i - 1
		g.You.Body[i] = g.You.Body[prev]
	}

	g.You.Head = g.You.Head.Step(dir)
	g.You.Body[0] = g.You.Head

	for i, s := range g.Board.Snakes {
		if s.ID == g.You.ID {
			g.Board.Snakes[i] = g.You
			break
		}
	}

	// Eat any food
	for i, food := range g.Board.Food {
		if food.Hit(&g.You.Head) {
			g.Board.Food = append(g.Board.Food[:i], g.Board.Food[i+1:]...)
			break
		}
	}
}

func (g *GameState) initFutureToPresent() BoardVision {
	future := make(BoardVision, g.Board.Width)

	for x := 0; x < g.Board.Width; x++ {
		future[x] = make([]VisionCell, g.Board.Height)
		for y := 0; y < g.Board.Height; y++ {
			future[x][y].Enemies = make([]float64, len(g.Board.Snakes))
		}
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
