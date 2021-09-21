package model

import (
	"bytes"
	"fmt"
)

type TreeNode struct {
	State    *GameState
	Parent   *TreeNode
	Children [4]*TreeNode
	Weight   float64
}

func (n *TreeNode) Ancestry() string {
	var buf bytes.Buffer

	for cur := n; cur.Parent != nil; cur = cur.Parent {
		buf.WriteString(fmt.Sprintf(" %s ", cur.State.You.Head))
	}

	return buf.String()
}
