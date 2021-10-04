package tree

import (
	"fmt"

	"github.com/itzamna314/battlesnake/game"
)

type Node struct {
	SnakeID string
	Parent  *Node
	Coord   *game.Coord

	Direction game.Direction
	Brain     SnakeBrain

	Depth  int
	Weight Weight
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}

	var s string
	for cur := n; cur != nil; cur = cur.Parent {
		s += fmt.Sprintf("%s[%s] ", cur.Coord, cur.Weight)
	}

	return s
}
