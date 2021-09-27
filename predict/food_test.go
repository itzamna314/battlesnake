package predict_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/predict"
)

func TestSingleFood(t *testing.T) {
	me := game.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   game.Coord{X: 2, Y: 0},
		Body:   []game.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 15,
	}
	input := game.GameState{
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

	testCases := []struct {
		testName  string
		health    int32
		minWeight int32
		maxWeight int32
	}{
		{"low", 15, predict.Neutral, predict.CertainWin},
		{"high", 99, predict.CertainDeath, predict.Neutral},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			me.Health = tt.health
			input.You = me
			input.Board.Snakes = []game.Battlesnake{me}

			var state predict.State
			state.Init(&input)

			w := state.Weight(&game.Coord{2, 1}, &me)
			if w < tt.minWeight {
				t.Errorf("Expected weight above %v, got %v", tt.minWeight, w)
			}

			if w > tt.maxWeight {
				t.Errorf("Expected weight above %v, got %v", tt.maxWeight, w)
			}
		})
	}
}
