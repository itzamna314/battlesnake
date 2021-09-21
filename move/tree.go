package move

import "github.com/itzamna314/battlesnake/model"

func BuildTree(start *model.GameState, depth int) *model.TreeNode {
	root := model.TreeNode{
		State: start,
	}
	expand(&root, depth)

	return &root
}

// expand expands children of node via DFS up to depth
func expand(node *model.TreeNode, depth int) {
	// If we're at the bottom, stop
	if depth <= 0 {
		return
	}

	// If this move is death, no need to calculate children
	if node.Weight < 0 {
		return
	}

	moves := model.Options(&node.State.You.Head)
	for dir, move := range moves {
		weight := weightSafe(node.State, &move.Coord)

		if weight >= 0 {
			weight += weightFood(node.State, &move.Coord)
		}

		// Advance game state in direction
		next := node.State.Clone()
		next.MoveSnake(next.You, model.Direction(dir))

		// Build child and recurse
		child := model.TreeNode{
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
