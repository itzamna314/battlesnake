package move

import "github.com/itzamna314/battlesnake/model"

func BuildTree(root *model.GameState, depth int) *model.GameTree {
	tree := &model.GameTree{
		Root: treeNode(root, depth),
	}

	return tree
}

// treeNode returns a node of the tree with all of its children,
// expanding children via DFS up to depth
func treeNode(state *model.GameState, depth int) *model.TreeNode {
	if depth < 0 {
		return nil
	}

	nd := &model.TreeNode{
		State: state,
		Moves: model.Options(&state.You.Head),
	}

	safe(nd.State, nd.Moves)
	food(nd.State, nd.Moves)

	for dir, move := range nd.Moves {
		if !move.Safe {
			continue
		}

		// Build next game state for move
		gs := *state

		// Copy each body segment to next
		// Head will remain copied into neck
		for i := len(gs.You.Body) - 1; i > 0; i-- {
			prev := i - 1
			gs.You.Body[i] = gs.You.Body[prev]
		}

		gs.You.Head = move.Coord
		gs.You.Body[0] = gs.You.Head

		nd.Children[dir] = treeNode(&gs, depth-1)
	}

	return nd
}
