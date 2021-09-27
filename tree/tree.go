package tree

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/game"
)

type Tree struct {
	// best holds the tree's current best node(s)
	best []*Node

	// curDepth tracks how far we traversed
	curBest       []*Node
	curBestWeight int32
	curDepth      int

	// weight holds the current set of nodes to weight
	weight chan *Node
	// expand holds the current set of nodes to expand
	expand chan *Node
}

func Search(ctx context.Context, state *game.GameState, snake *game.Battlesnake, brain SnakeBrain) game.Direction {
	brain.Init(state)

	root := Node{
		Coord: &snake.Head,
		Snake: snake,
		Brain: brain,
	}

	t := Tree{
		// 100k buffer will hold 3^10 expansions, even without pruning
		weight: make(chan *Node, 100000),
		expand: make(chan *Node, 100000),
	}

	// start a weight worker
	// We could start more workers here safely
	go weightWorker(ctx, t.weight, t.expand)

	// Send the root node for expansion
	go func() {
		t.expand <- &root
	}()

	// Block until cancelled
	t.expandWorker(ctx)

	// Process results
	if t.best == nil {
		log.Printf("No safe moves detected (best nil)! Moving down\n")
		return game.Down
	}

	var (
		bestMoves []game.Direction
		bestMove  game.Direction
	)

	// Traverse the ancestry of best to find its starting move
	for _, best := range t.best {
		// best is our root node. We didn't explore anything
		if best.Parent == nil {
			continue
		}

		fmt.Printf("Best move %s[%v] (%v)\n", best.Coord, best.Depth, best.Weight)

		// We want to find the move who's parent is root
		for cur := best; cur.Parent != nil; cur = cur.Parent {
			bestMove = cur.Direction
		}
		bestMoves = append(bestMoves, bestMove)
	}

	if len(bestMoves) == 0 {
		log.Printf("No safe moves detected (no child moves)! Moving down\n")
		return game.Down
	}

	if len(bestMoves) == 1 {
		return bestMoves[0]
	}

	// Pick randomly if multiple best moves
	return bestMoves[rand.Intn(len(bestMoves))]
}
