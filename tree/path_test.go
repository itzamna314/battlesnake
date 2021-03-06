package tree_test

import (
	"context"
	"testing"
	"time"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
)

func TestRightSearch(t *testing.T) {
	// Arrange
	me := game.Battlesnake{
		// Length 3, facing right
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 10,
			Width:  10,
			Snakes: []game.Battlesnake{me},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	// This brain always pulls to the right
	rightBrain := &testBrain{
		weightFunc: func(coord *game.Coord, snakeID string) float64 {
			return float64(coord.X)
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	mv, meta := tree.Search(ctx, &state, &me, rightBrain)
	t.Logf("Search meta %+v\n", meta)
	if mv != game.Right {
		t.Errorf("quick search with right-only brain did not go right. Went %s", mv)
	}
}

func TestSeekSearch(t *testing.T) {
	// Arrange
	me := game.Battlesnake{
		// Length 3, facing right
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 10,
			Width:  10,
			Snakes: []game.Battlesnake{me},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	// This brain is aiming for 0, 10
	target := &game.Coord{2, 10}
	// abort if we're getting further away
	origDist := me.Head.Dist(target)

	upBrain := &testBrain{
		weightFunc: func(coord *game.Coord, snakeID string) float64 {
			newDist := coord.Dist(target)

			return float64(origDist - newDist)
		},
		abortFunc: func(weight float64) bool {
			return weight < 0
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	mv, meta := tree.Search(ctx, &state, &me, upBrain)
	t.Logf("Search meta %+v\n", meta)
	if mv != game.Up {
		t.Errorf("quick search with right-only brain did not go right. Went %s", mv)
	}
}

func TestDeterioratingPath(t *testing.T) {
	// Arrange
	me := game.Battlesnake{
		// Length 3, facing right
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 10,
	}
	state := game.GameState{}

	world := map[game.Coord]float64{
		game.Coord{0, 4}: 2, game.Coord{1, 4}: 0, game.Coord{2, 4}: 0, game.Coord{3, 4}: 0, game.Coord{4, 4}: 0,
		game.Coord{0, 3}: 2, game.Coord{1, 3}: 0, game.Coord{2, 3}: 0, game.Coord{3, 3}: 0, game.Coord{4, 3}: 0,
		game.Coord{0, 2}: 2, game.Coord{1, 2}: 0, game.Coord{2, 2}: 0, game.Coord{3, 2}: 0, game.Coord{4, 2}: 0,
		game.Coord{0, 1}: 2, game.Coord{1, 1}: 0, game.Coord{2, 1}: 1, game.Coord{3, 1}: -11, game.Coord{4, 1}: -11,
		game.Coord{0, 0}: 2, game.Coord{1, 0}: 2, game.Coord{2, 0}: -11, game.Coord{3, 0}: 10, game.Coord{4, 0}: -11,
	}

	// This brain follows the weights we describe in hardcoded world coordinates
	worldBrain := &testBrain{
		weightFunc: func(coord *game.Coord, snakeID string) float64 {
			w, ok := world[*coord]
			if !ok {
				return -100
			}
			return w
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	mv, meta := tree.Search(ctx, &state, &me, worldBrain)
	t.Logf("Search meta %+v\n", meta)
	if mv != game.Left {
		t.Errorf("did not take best path (left). Went %s", mv)
	}
}
