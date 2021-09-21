package model

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

func (g *GameState) Clone() GameState {
	return GameState{
		Game:  g.Game,
		Turn:  g.Turn,
		Board: g.Board.Clone(),
		You:   g.You.Clone(),
	}
}

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

func (b *Board) Clone() Board {
	clone := Board{
		Height: b.Height,
		Width:  b.Width,
		Food:   make([]Coord, len(b.Food)),
		Snakes: make([]Battlesnake, len(b.Snakes)),
	}

	for i, f := range b.Food {
		clone.Food[i] = f
	}

	for i, s := range b.Snakes {
		clone.Snakes[i] = s.Clone()
	}

	return clone
}

type Battlesnake struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Health  int32   `json:"health"`
	Body    []Coord `json:"body"`
	Head    Coord   `json:"head"`
	Length  int32   `json:"length"`
	Latency string  `json:"latency"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

func (b *Battlesnake) Clone() Battlesnake {
	clone := Battlesnake{
		ID:      b.ID,
		Name:    b.Name,
		Health:  b.Health,
		Head:    b.Head,
		Length:  b.Length,
		Latency: b.Latency,
		Body:    make([]Coord, len(b.Body)),
	}

	for i, bd := range b.Body {
		clone.Body[i] = bd
	}

	return clone
}

func (gs *GameState) MoveSnake(snake Battlesnake, dir Direction) {
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

	for i, s := range gs.Board.Snakes {
		if s.ID == snake.ID {
			gs.Board.Snakes[i] = snake
			break
		}
	}

	if snake.ID == gs.You.ID {
		gs.You = snake
	}

	// Eat any food
	for i, food := range gs.Board.Food {
		if food.Hit(&snake.Head) {
			gs.Board.Food = append(gs.Board.Food[:i], gs.Board.Food[i+1:]...)
			break
		}
	}
}
