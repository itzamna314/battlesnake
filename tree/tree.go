package tree

import (
	"github.com/itzamna314/battlesnake/model"
	"github.com/itzamna314/battlesnake/move"
)

func Build(start *model.GameState, depth int) *TreeNode {
	initialState := start.Clone()
	root := TreeNode{
		State: &initialState,
	}

	expand(&root, depth)

	return &root
}

// expand expands children of node via DFS up to depth
func expand(node *TreeNode, depth int) {
	// If we're at the bottom, stop
	if depth <= 0 {
		return
	}

	// If this move is death, no need to calculate children
	if node.Weight < 0 {
		return
	}

	opts := model.Options(&node.State.You.Head)
	for dir, opt := range opts {
		// Handle for enemy movement
		node.State.MoveEnemies()

		// Calculate weight based on enemy predictions
		weight := move.Weight(node.State, &opt.Coord)

		// Advance game state in direction
		next := node.State.Clone()
		next.Move(model.Direction(dir))

		// Build child and recurse
		child := TreeNode{
			Parent: node,
			Weight: weight,
			State:  &next,
		}

		expand(&child, depth-1)
		node.Children[dir] = &child
	}

	// Add best child to our weight
	bestChild := -1.0
	for _, child := range node.Children {
		if child == nil {
			continue
		}

		if child.Weight > bestChild {
			bestChild = child.Weight
		}
	}

	node.Weight = node.Weight + (bestChild * 0.5)
	return
}
