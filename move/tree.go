package move

import "github.com/itzamna314/battlesnake/game"

func BuildTree(start *game.GameState, depth int) *TreeNode {
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

	// Guess where enemies will go before assessing options
	node.State.MoveEnemies()

	opts := game.Options(&node.State.You.Head)
	for dir, opt := range opts {

		// Calculate weight based on enemy predictions
		weight := Weight(node.State, &opt.Coord)

		// Advance game state in direction
		next := node.State.Clone()
		next.Move(game.Direction(dir))

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

	if bestChild <= -1.0 {
		node.Weight = -1.0
	} else {
		node.Weight = node.Weight + (bestChild * 0.5)
	}
}
