package tree

import (
	"github.com/itzamna314/battlesnake/game"
)

type Node struct {
	Snake  *game.Battlesnake
	Parent *Node
	Coord  *game.Coord

	Direction game.Direction
	Brain     SnakeBrain

	Depth  int
	Weight int32
}
