package tree

import (
	"context"
	"log"
	"math/rand"
	"reflect"

	"github.com/itzamna314/battlesnake/game"
	"github.com/pbnjay/memory"
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

	// weight holds the current set of nodes to weight
	weight chan *Node
	// expand holds the current set of nodes to expand
	expand chan *Node
}

type SearchMetadata struct {
	Depth   int
	Memsize int
}

func Search(ctx context.Context,
	state *game.GameState,
	snake *game.Battlesnake,
	brain SnakeBrain,
	cfg ...ConfigFn,
) game.Direction {
	var t Tree

	for _, c := range cfg {
		c(&t)
	}

	brain.Init(state)

	root := Node{
		Coord: &snake.Head,
		Snake: snake,
		Brain: brain,
	}

	freeMem := memory.FreeMemory()
	nodeVal := reflect.TypeOf(root)
	nodeSize := nodeVal.Size()
	maxNodes := freeMem / uint64(nodeSize)

	sizeToUse := maxNodes / 100
	chBuffer := sizeToUse / 2

	t.weight = make(chan *Node, chBuffer)
	t.expand = make(chan *Node, chBuffer)

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
