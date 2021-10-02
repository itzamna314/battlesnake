package tree

import (
	"context"
	"log"
	"math/rand"

	"github.com/itzamna314/battlesnake/game"
)

type Tree struct {
	MaxDepth int

	// best holds the tree's current best node(s)
	best []*Node

	// cur fields tracks the level we're currently expanding
	curBest       []*Node
	curBestWeight float64
	curDepth      int
	curLeft       int

	// next fields track the next level we would expand
	nextWidth int
	nextLevel [][]*Node

	// weight holds the current set of nodes to weight
	weight chan *Node
	// expand holds the current set of nodes to expand
	expand chan *Node
}

type SearchMetadata struct {
	Depth  int
	Weight float64
}

func Search(ctx context.Context,
	state *game.GameState,
	snake *game.Battlesnake,
	brain SnakeBrain,
	cfg ...ConfigFn,
) (game.Direction, *SearchMetadata) {
	t := Tree{
		weight:    make(chan *Node),
		expand:    make(chan *Node),
		nextLevel: make([][]*Node, 1),
	}

	for _, c := range cfg {
		c(&t)
	}

	brain.Init(state)

	snakeClone := snake.Clone()
	root := Node{
		Snake: &snakeClone,
		Coord: &snakeClone.Head,
		Brain: brain,
	}

	// Send the root node for expansion
	go func() {
		select {
		case t.expand <- &root:
		case <-ctx.Done():
		}
	}()

	// start a weight worker
	// We could start more workers here safely
	go weightWorker(ctx, t.weight, t.expand)

	// Block until cancelled
	t.expandWorker(ctx)

	var meta SearchMetadata

	// Process results
	if t.best == nil {
		log.Printf("No safe moves detected (best nil)! Moving down\n")
		return game.Down, &meta
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

		// We want to find the move who's parent is root
		for cur := best; cur.Parent != nil; cur = cur.Parent {
			bestMove = cur.Direction
		}

		meta.Depth = best.Depth
		meta.Weight = best.Weight
		bestMoves = append(bestMoves, bestMove)
	}

	if len(bestMoves) == 0 {
		log.Printf("No safe moves detected (no child moves)! Moving down\n")
		return game.Down, &meta
	}

	if len(bestMoves) == 1 {
		return bestMoves[0], &meta
	}

	// Pick randomly if multiple best moves
	return bestMoves[rand.Intn(len(bestMoves))], &meta
}
