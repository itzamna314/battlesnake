package multiverse_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/brain/multiverse"
	"github.com/itzamna314/battlesnake/game"
	"github.com/itzamna314/battlesnake/tree"
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
		minWeight multiverse.FloatWeight
		maxWeight multiverse.FloatWeight
	}{
		{"low_health", 15, multiverse.Avoid, multiverse.FloatWeight(1.0)},
		{"high_health", 99, multiverse.Death, multiverse.Base},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			me.Health = tt.health
			input.You = me
			input.Board.Snakes = []game.Battlesnake{me}

			var state multiverse.State
			state.Init(&input)

			state.Move(me.ID, game.Up)

			w := state.Weight(&tree.Node{
				SnakeID:   me.ID,
				Coord:     &me.Head,
				Direction: game.Up,
				Brain:     &state,
				Depth:     1,
			})
			if w.Compare(tt.minWeight) < 0 {
				t.Errorf("Expected weight above %v, got %v", tt.minWeight, w)
			}

			if w.Compare(tt.maxWeight) > 0 {
				t.Errorf("Expected weight below %v, got %v", tt.maxWeight, w)
			}
		})
	}
}
