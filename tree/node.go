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
	Weight float64
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}

	var s string
	for cur := n; cur != nil; cur = cur.Parent {
		s += fmt.Sprintf("%s[%.2f] ", cur.Coord, cur.Weight)
	}

	return s
}
