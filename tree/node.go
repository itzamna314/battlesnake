package tree

import (
	"bytes"
	"fmt"

	"github.com/itzamna314/battlesnake/model"
)

type TreeNode struct {
	State    *model.GameState
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
