package move

import (
	"bytes"
	"fmt"

	"github.com/itzamna314/battlesnake/predict"
)

type TreeNode struct {
	State    *predict.State
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
