package move_test

import (
	"testing"

	"github.com/itzamna314/battlesnake/handlers/move"
	"github.com/itzamna314/battlesnake/model"
)

func TestNeckAvoidance(t *testing.T) {
	// Arrange
	me := model.Battlesnake{
		// Length 3, facing right
		Head: model.Coord{X: 2, Y: 0},
		Body: []model.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := model.GameState{
		Board: model.Board{
			Snakes: []model.Battlesnake{me},
		},
		You: me,
	}

	// Act 1,000x (this isn't a great way to test, but it's okay for starting out)
	for i := 0; i < 1000; i++ {
		nextMove := move.Next(state)
		// Assert never move left
		if nextMove.Move == "left" {
			t.Errorf("snake moved onto its own neck, %s", nextMove.Move)
		}
	}
}

// TODO: More GameState test cases!
