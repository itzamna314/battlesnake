package api_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/api"
	"github.com/itzamna314/battlesnake/model"
	"github.com/itzamna314/battlesnake/testdata"
)

func TestEatOne(t *testing.T) {
	// Arrange
	me := model.Battlesnake{
		// Length 3, facing right
		Head:   model.Coord{X: 2, Y: 0},
		Body:   []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := model.GameState{
		Board: model.Board{
			Height: 10,
			Width:  10,
			Snakes: []model.Battlesnake{me},
			Food: []model.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	nextMove := api.NextMove(state)
	if nextMove.Move != "up" {
		t.Errorf("snake did not eat food at (2,1), went %s", nextMove.Move)
	}
}

func TestEatFuture(t *testing.T) {
	// Arrange
	me := model.Battlesnake{
		// Length 3, facing right
		Head:   model.Coord{X: 2, Y: 0},
		Body:   []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 1,
	}
	state := model.GameState{
		Board: model.Board{
			Height: 10,
			Width:  10,
			Snakes: []model.Battlesnake{me},
			Food: []model.Coord{
				{1, 1},
				{3, 0},
				{4, 0},
				{7, 0},
			},
		},
		You: me,
	}

	nextMove := api.NextMove(state)
	if nextMove.Move != "right" {
		t.Errorf("snake did not eat 2 food, went %s", nextMove.Move)
	}
}

func TestWithEnemies(t *testing.T) {
	me := model.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   model.Coord{X: 2, Y: 0},
		Body:   []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 80,
	}
	enemy := model.Battlesnake{
		// Length 3, facing down
		ID:     "enemy",
		Head:   model.Coord{X: 3, Y: 1},
		Body:   []model.Coord{{X: 3, Y: 1}, {X: 3, Y: 2}, {X: 2, Y: 2}},
		Health: 80,
	}
	state := model.GameState{
		Board: model.Board{
			Height: 10,
			Width:  10,
			Snakes: []model.Battlesnake{me, enemy},
			Food: []model.Coord{
				{2, 1},
			},
		},
		You: me,
	}

	nextMove := api.NextMove(state)
	if nextMove.Move != "up" {
		t.Errorf("snake did not eat food at (2,1), went %s", nextMove.Move)
	}
}

func TestFoodOrDeath(t *testing.T) {
	me := model.Battlesnake{
		ID:   "me",
		Head: model.Coord{X: 4, Y: 6},
		Body: []model.Coord{
			{X: 4, Y: 6}, {X: 5, Y: 6}, {X: 6, Y: 6}, {X: 6, Y: 7}, {X: 7, Y: 7}, {X: 8, Y: 7},
		},
		Health: 20,
	}
	enemy := model.Battlesnake{
		ID:   "enemy",
		Head: model.Coord{X: 3, Y: 4},
		Body: []model.Coord{
			{X: 3, Y: 5}, {X: 4, Y: 5}, {X: 5, Y: 5}, {X: 6, Y: 5},
		},
		Health: 80,
	}
	state := model.GameState{
		Board: model.Board{
			Height: 10,
			Width:  10,
			Snakes: []model.Battlesnake{me, enemy},
			Food: []model.Coord{
				{4, 4},
				{5, 4},
				{4, 3},
			},
		},
		You: me,
	}

	nextMove := api.NextMove(state)
	if nextMove.Move == "down" {
		t.Errorf("snake went down into death")
	}
}

func TestFrames(t *testing.T) {
	testCases := []struct {
		frame        string
		allowedMoves []model.Direction
	}{
		{"afraid_to_eat", []model.Direction{model.Up}},
	}

	for _, tt := range testCases {
		t.Run(tt.frame, func(t *testing.T) {
			input, ok := testdata.Frame(tt.frame)
			if !ok {
				t.Fatalf("Failed to find frame %s", tt.frame)
			}

			state := input.Clone()
			nextMove := api.NextMove(state)

			for _, mv := range tt.allowedMoves {
				if mv.String() == nextMove.Move {
					t.Logf("Made allowed move %s", mv)
					return
				}
			}

			t.Errorf("Made disallowed move %s. Allowed: %v", nextMove.Move, tt.allowedMoves)
		})
	}
}
