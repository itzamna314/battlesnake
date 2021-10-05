package tree

import (
	"context"

	"github.com/itzamna314/battlesnake/game"
)

// expandWorker listens to the tree's expand channel
// It checks each next node to see if its our new best, and
// adds all of its children to the weight channel
func (t *Tree) expandWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case exp, ok := <-t.expand:
			if !ok {
				return
			}

			if t.curBest == nil {
				t.curBest = []*Node{exp}
				t.curBestWeight = exp.Weight
			} else if diff := exp.Weight.Compare(t.curBestWeight); diff > 0 {
				t.curBest = []*Node{exp}
				t.curBestWeight = exp.Weight
			} else if diff == 0 {
				t.curBest = append(t.curBest, exp)
			}

			// Expand this node if:
			// * We want to explore the next level
			// * The brain thinks its worth exploring this node further
			if t.MaxDepth == 0 || t.curDepth < t.MaxDepth {
				// TODO: Rename Abort => ExpandOrPrune
				// Return nodes for next level
				// Minimax - flip active snake
				// Multiverse - Keep same active snake
				if exp.Weight == nil || !exp.Brain.Abort(exp) {
					t.expandNode(ctx, exp)
				}
			}

			// We're processing a finished node at the current level
			t.curLeft--

			// Finishing a level. Store the results of the last depth.
			// Re-set the current best
			if t.curLeft <= 0 {
				// Store the results of this level, and queue up the next one
				// When all of the next level has been queued, we can begin
				// expanding it out
				t.completeLevel(ctx)

				// Exit if this is the final level
				if t.MaxDepth != 0 && t.curDepth > t.MaxDepth {
					return
				}
			}
		}
	}
}

func (t *Tree) expandNode(ctx context.Context, node *Node) {
	for dir, opt := range game.Options(node.Coord) {
		child := Node{
			Direction: game.Direction(dir),
			SnakeID:   node.SnakeID,
			Coord:     opt,
			Brain:     node.Brain.Clone(),
			Parent:    node,
			Depth:     node.Depth + 1,
		}

		t.nextLevel[node.Depth] = append(t.nextLevel[node.Depth], &child)
		t.nextWidth++
	}
}

// completeLevel marks the current level as complete
// Called from the expand worker when we finish processing a level
// of the tree.
// Return true to continue traversal, false to exit
func (t *Tree) completeLevel(ctx context.Context) {
	t.best = t.curBest
	t.curBest = nil

	t.curLeft = t.nextWidth
	t.nextWidth = 0

	t.nextLevel = append(t.nextLevel, []*Node{})

	go func(depth int) {
		for _, n := range t.nextLevel[depth] {
			select {
			case t.weight <- n:
			case <-ctx.Done():
				return
			}
		}
		t.nextLevel[depth] = nil
	}(t.curDepth)

	t.curDepth++
}
