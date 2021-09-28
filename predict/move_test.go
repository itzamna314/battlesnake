package predict_test

import (
	"context"
	"testing"
	"time"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
	"github.com/itzamna314/battlesnake/tree"
)

func testTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	return ctx
}

func TestEatOne(t *testing.T) {
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

	mv := tree.Search(testTimeout(), &state, &me, &predict.State{})
	if mv != game.Up {
		t.Errorf("snake did not eat food at (2,1), went %s", mv)
	}
}

func TestEatFuture(t *testing.T) {
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
				{1, 1},
				{3, 0},
				{4, 0},
				{7, 0},
			},
		},
		You: me,
	}

	mv := tree.Search(testTimeout(), &state, &me, &predict.State{})
	if mv != game.Right {
		t.Errorf("snake did not eat 2 food, went %s", mv)
	}
}

func TestWithEnemies(t *testing.T) {
	me := game.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 80,
	}
	enemy := game.Battlesnake{
		// Length 3, facing down
		ID:     "enemy",
		Head:   game.Coord{X: 3, Y: 1},
		Body:   []game.Coord{{X: 3, Y: 1}, {X: 3, Y: 2}, {X: 2, Y: 2}},
		Health: 80,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 10,
			Width:  10,
			Snakes: []game.Battlesnake{me, enemy},
			Food: []game.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	mv := tree.Search(testTimeout(), &state, &me, &predict.State{})
	if mv != game.Up {
		t.Errorf("snake did not eat food at (2,1), went %s", mv)
	}
}

func TestFoodOrDeath(t *testing.T) {
	me := game.Battlesnake{
		ID:   "me",
		Head: game.Coord{X: 4, Y: 6},
		Body: []game.Coord{
			{X: 4, Y: 6}, {X: 5, Y: 6}, {X: 6, Y: 6}, {X: 6, Y: 7}, {X: 7, Y: 7}, {X: 8, Y: 7},
		},
		Health: 20,
	}
	enemy := game.Battlesnake{
		ID:   "enemy",
		Head: game.Coord{X: 3, Y: 4},
		Body: []game.Coord{
			{X: 3, Y: 5}, {X: 4, Y: 5}, {X: 5, Y: 5}, {X: 6, Y: 5},
		},
		Health: 80,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 10,
			Width:  10,
			Snakes: []game.Battlesnake{me, enemy},
			Food: []game.Coord{
				{4, 4},
				{5, 4},
				{4, 3},
			},
		},
		You: me,
	}

	mv := tree.Search(testTimeout(), &state, &me, &predict.State{})
	if mv == game.Down {
		t.Errorf("snake went down into death")
	}
}

func TestAvoidDeath(t *testing.T) {
	me := game.Battlesnake{
		ID:   "me",
		Head: game.Coord{X: 4, Y: 6},
		Body: []game.Coord{
			{X: 8, Y: 6}, {X: 7, Y: 6}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 5, Y: 7}, {X: 4, Y: 7}, {X: 3, Y: 7}, {X: 2, Y: 7}, {X: 2, Y: 7},
		},
		Health: 100,
	}
	enemy := game.Battlesnake{
		ID:   "enemy",
		Head: game.Coord{X: 3, Y: 4},
		Body: []game.Coord{
			{X: 7, Y: 5}, {X: 6, Y: 5}, {X: 5, Y: 5}, {X: 4, Y: 5}, {X: 4, Y: 4}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 5, Y: 2}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1},
		},
		Health: 94,
	}
	state := game.GameState{
		Board: game.Board{
			Height: 11,
			Width:  11,
			Snakes: []game.Battlesnake{me, enemy},
			Food: []game.Coord{
				{8, 0},
			},
		},
		You: me,
	}

	mv := tree.Search(testTimeout(), &state, &state.You, &predict.State{}, tree.ConfigMaxDepth(4))
	if mv == game.Down {
		t.Errorf("snake went down into possible head collision")
	}
}
