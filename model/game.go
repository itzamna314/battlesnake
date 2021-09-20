package model

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
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

func (gs *GameState) MoveSnake(snake Battlesnake, dir Direction) {
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
