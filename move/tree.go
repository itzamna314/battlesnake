package move

import (
	"fmt"

	"github.com/itzamna314/battlesnake/model"
)

func BuildTree(root *model.GameState, depth int) *model.GameTree {
	tree := &model.GameTree{
		Root: treeNode(root, 0, depth),
	}

	return tree
}

// treeNode returns a node of the tree with all of its children,
// expanding children via DFS up to depth
func treeNode(state *model.GameState, weight float64, depth int) *model.TreeNode {
	nd := &model.TreeNode{
		State:    state,
		Children: [4]*model.TreeNode{},
		Weight:   weight,
	}

	moves := model.Options(&state.You.Head)

	safe(nd.State, moves)
	food(nd.State, moves)

	// If we're at the bottom, return
	if depth <= 0 {
		return nd
	}

	for dir, move := range moves {
		// Prune branches that die
		if move.Weight < 0 {
			continue
		}

		// Build next game state for move
		gs := *state
		gs.MoveSnake(gs.You, model.Direction(dir))

		nd.Children[dir] = treeNode(&gs, move.Weight, depth-1)
	}

	var (
		numChildren  float64
		childWeights float64
	)

	// Copy child node weights up to this node
	for _, child := range nd.Children {
		// Pruned child, weight -1
		if child == nil {
			continue
		}

		numChildren += 1
		childWeights += child.Weight
	}

	nd.Weight = nd.Weight + (childWeights / numChildren)
	return nd
}

func depthPrintf(depth int, str string, args ...interface{}) {
	d := 3 - depth
	for i := 0; i < d; i++ {
		fmt.Printf(".")
	}
	fmt.Printf(str, args...)
}
