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
				return
			}

			// Starting a new level. Store the results of the last depth.
			// Re-set the current best
			if exp.Depth > t.curDepth {
				fmt.Printf("Finished depth %v, best:\n", t.curDepth)
				for _, b := range t.curBest {
					fmt.Printf("\t%s\n", b)
				}
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

			// This is the final level. Stop expanding, and stream
			// cur into best
			if t.MaxDepth != 0 && exp.Depth == t.MaxDepth {
				t.best = t.curBest

				if exp.Coord.Hit(&game.Coord{2, 10}) {
					fmt.Printf("Draining expand channel [%d]:\n\t%s\n", exp.Depth, exp)
				}
				continue
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
					return
				case t.weight <- &child:
				}
			}
		}
	}
}
