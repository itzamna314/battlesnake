package move_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/move"
	"github.com/itzamna314/battlesnake/testdata"
)

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

	mv := move.Next(state)
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

	mv := move.Next(state)
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

	mv := move.Next(state)
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

	mv := move.Next(state)
	if mv == game.Down {
		t.Errorf("snake went down into death")
	}
}

func TestFrames(t *testing.T) {
	testCases := []struct {
		frame        string
		allowedMoves []game.Direction
	}{
		{"afraid_to_eat", []game.Direction{game.Up}},
		{"no_mercy", []game.Direction{game.Right}},
		{"enemy_ate", []game.Direction{game.Down}},
		{"leave_hazard", []game.Direction{game.Left}},
		{"over_chase", []game.Direction{game.Right}},
		{"corner_crash", []game.Direction{game.Left}},
		{"bad_joust", []game.Direction{game.Down}},
		{"pessimistic", []game.Direction{game.Left, game.Up}},
	}

	for _, tt := range testCases {
		t.Run(tt.frame, func(t *testing.T) {
			input, ok := testdata.Frame(tt.frame)
			if !ok {
				t.Fatalf("Failed to find frame %s", tt.frame)
			}

			state := input.Clone()
			next := move.Next(state)

			for _, allowed := range tt.allowedMoves {
				if allowed == next {
					t.Logf("Made allowed move %s", next)
					return
				}
			}

			t.Errorf("Made disallowed move %s. Allowed: %v", next, tt.allowedMoves)
		})
	}
}
