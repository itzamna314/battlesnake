package tree

import (
	"context"
	"fmt"

	"github.com/itzamna314/battlesnake/game"
)

// expandWorker listens to the tree's expand channel
// It checks each next node to see if its our new best, and
// adds all of its children to the weight channel
func (t *Tree) expandWorker(ctx context.Context) {
	defer func() {
		close(t.weight)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case exp, ok := <-t.expand:
			if !ok {
				fmt.Println("expand channel closed. Exit")
				return
			}

			if t.curBest == nil || exp.Weight > t.curBestWeight {
				t.curBest = []*Node{exp}
				t.curBestWeight = exp.Weight
			} else if exp.Weight == t.curBestWeight {
				t.curBest = append(t.curBest, exp)
			}

			// Expand this node if:
			// * We want to explore the next level
			// * The brain thinks its worth exploring this node further
			if t.MaxDepth == 0 || t.curDepth < t.MaxDepth {
				if !exp.Brain.Abort(exp.Weight) {
					t.expandNode(ctx, exp)
				}
			}

			// We're processing a finished node at the current level
			t.curLeft--

			// Finishing a level. Store the results of the last depth.
			// Re-set the current best
			if t.curLeft <= 0 {
				// Store the results of this level
				t.completeLevel()

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
		snakeClone := node.Snake.Clone()
		child := Node{
			Direction: game.Direction(dir),
			Snake:     &snakeClone,
			Coord:     opt,
			Brain:     node.Brain.Clone(),
			Parent:    node,
			Depth:     node.Depth + 1,
		}

		select {
		case <-ctx.Done():
			return
		case t.weight <- &child:
			t.nextWidth++
		}
	}
}

// completeLevel marks the current level as complete
// Called from the expand worker when we finish processing a level
// of the tree.
// Return true to continue traversal, false to exit
func (t *Tree) completeLevel() {
	t.best = t.curBest
	t.curBest = nil

	t.curLeft = t.nextWidth
	t.nextWidth = 0

	t.curDepth++
}
