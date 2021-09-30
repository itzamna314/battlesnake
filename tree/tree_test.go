package tree_test

import (
	"context"
	"testing"
	"time"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
)

func TestExitContext(t *testing.T) {
	me := game.Battlesnake{
		// Length 3, facing right
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 11,
			Width:  11,
			Snakes: []game.Battlesnake{me},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	var brain testBrain

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)

	done := make(chan struct{}, 1)

	go func() {
		tree.Search(ctx, &state, &me, &brain)
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("Search did not exit in time")
	}
}

func TestExitDepth(t *testing.T) {
	me := game.Battlesnake{
		// Length 3, facing right
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 11,
			Width:  11,
			Snakes: []game.Battlesnake{me},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	var brain testBrain

	ctx := context.Background()

	done := make(chan struct{}, 1)

	go func() {
		tree.Search(ctx, &state, &me, &brain, tree.ConfigMaxDepth(2))
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(20 * time.Millisecond):
		t.Fatalf("Search did not exit in time")
	}
}
