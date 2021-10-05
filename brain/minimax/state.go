package minimax

import (
	"github.com/itzamna314/battlesnake/game"
)

type State struct {
	MaxScore *int32
	MinScore *int32

	Height int
	Width  int

	MaxSnake *game.Battlesnake
	MinSnake *game.Battlesnake

	Hazards []game.Coord
	Food    []game.Coord
}