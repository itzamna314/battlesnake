package move_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/model"
	"github.com/itzamna314/battlesnake/move"
)

func TestSingleFood(t *testing.T) {
	me := model.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   model.Coord{X: 2, Y: 0},
		Body:   []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 15,
	}
	input := model.GameState{
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

	state := input.Clone()

	w := move.WeightFood(&state, &model.Coord{2, 1})
	if w != 0.15 {
		t.Errorf("Expected weight for guaranteed food, low health at (2,1) 0.15, got %v", w)
	}

	state.You.Health = 99
	w = move.WeightFood(&state, &model.Coord{2, 1})
	if w != -0.15 {
		t.Errorf("Expected weight for guaranteed food, high health at (2,1) -0.15, got %v", w)
	}
}

func TestMaxFood(t *testing.T) {
	me := model.Battlesnake{
		// Length 3, facing right
		ID:     "me",
		Head:   model.Coord{X: 1, Y: 1},
		Body:   []model.Coord{{X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 99,
	}
	input := model.GameState{
		Board: model.Board{
			Height: 5,
			Width:  5,
			Snakes: []model.Battlesnake{me},
			Food: []model.Coord{
				{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
				{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
				{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
				{0, 1}, {2, 1}, {3, 1}, {4, 1},
				{2, 0}, {3, 0}, {4, 0},
			},
		},
		You: me,
	}

	state := input.Clone()

	legalMoves := []model.Coord{{0, 1}, {1, 2}, {2, 1}}

	for _, l := range legalMoves {
		w := move.Weight(&state, &l)
		if w < 0 || w > 1 {
			t.Errorf("Weight exceeding bounds for %s: %v", l, w)
		}

		t.Logf("weight: %v", w)
	}
}
