package tree

import (
	"context"

	"github.com/itzamna314/battlesnake/game"
)

// expandWorker listens to the tree's expand channel
// It checks each next node to see if its our new best, and
// adds all of its children to the weight channel
func (t *Tree) expandWorker(ctx context.Context) {
Listen:
	for {
		select {
		case <-ctx.Done():
			break Listen
		case exp, ok := <-t.expand:
			if !ok {
				break Listen
			}

			// Starting a new level. Store the results of the last depth.
			// Re-set the current depth
			if exp.Depth > t.curDepth {
				t.curDepth = exp.Depth
				t.best = t.curBest
				t.curBest = nil
			}

			// Ask the brain if we should abort search
			if exp.Brain.Abort(exp.Weight) {
				continue
			}

			if t.curBest == nil || exp.Weight > t.curBestWeight {
				t.curBest = []*Node{exp}
				t.curBestWeight = exp.Weight
			} else if exp.Weight == t.curBestWeight {
				t.curBest = append(t.curBest, exp)
			}

			for dir, opt := range game.Options(exp.Coord) {
				snakeClone := exp.Snake.Clone()
				child := Node{
					Direction: game.Direction(dir),
					Snake:     &snakeClone,
					Coord:     opt,
					Brain:     exp.Brain.Clone(),
					Parent:    exp,
					Depth:     exp.Depth + 1,
				}

				select {
				case <-ctx.Done():
					break Listen
				case t.weight <- &child:
				}
			}
		}
	}
}
